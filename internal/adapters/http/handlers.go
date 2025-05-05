package http

import (
	"net/http"

	"github.com/gabrielrauch/reconcile-auth-service/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	app    *app.AuthService
	logger zerolog.Logger
}

func RunServer(port string, authApp *app.AuthService, logger zerolog.Logger) {
	r := gin.Default()

	r.Use(RequestIDMiddleware(logger))

	h := AuthHandler{app: authApp, logger: logger}

	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/auth/validate", h.Validate)
	r.GET("/admin/secret", RoleMiddleware("SUPER_ADMIN"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, super admin!"})
	})

	r.Run(":" + port)
}

func (h *AuthHandler) Register(c *gin.Context) {
	requestID := GetRequestID(c.Request.Context())

	tokenStr := c.GetHeader("Authorization")
	claims := jwt.MapClaims{}
	var authenticatedRole string

	if tokenStr != "" {
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret"), nil
		})

		if err != nil || !token.Valid {
			h.logger.Error().Str("request_id", requestID).Msg("invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		authenticatedRole = claims["role"].(string)
	}

	var input struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error().Err(err).Str("request_id", requestID).Msg("failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authenticatedRole == "SUPER_ADMIN" {
		if input.Role != "admin" {
			h.logger.Warn().Str("request_id", requestID).Msg("only admin can register admin users")
			c.JSON(http.StatusForbidden, gin.H{"error": "only super admin can register admin users"})
			return
		}
	} else {
		input.Role = "USER"
	}

	if err := h.app.Register(input.FirstName, input.Email, input.Password, input.Role, requestID); err != nil {
		h.logger.Error().Str("request_id", requestID).Msg("registration failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().Str("request_id", requestID).Msg("user registered successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	requestID := GetRequestID(c.Request.Context())

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error().Err(err).Str("request_id", requestID).Msg("failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.app.Login(input.Email, input.Password, requestID)
	if err != nil {
		h.logger.Error().Str("request_id", requestID).Msg("login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().Str("request_id", requestID).Msg("user logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Validate(c *gin.Context) {
	requestID := GetRequestID(c.Request.Context())

	token := c.GetHeader("Authorization")
	if token == "" || !h.app.ValidateToken(token, requestID) {
		h.logger.Warn().Str("request_id", requestID).Msg("invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	h.logger.Info().Str("request_id", requestID).Msg("token validated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "valid token"})
}
