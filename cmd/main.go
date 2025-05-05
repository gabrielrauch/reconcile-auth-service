package main

import (
	"os"

	"github.com/gabrielrauch/reconcile-auth-service/internal/adapters/http"
	"github.com/gabrielrauch/reconcile-auth-service/internal/adapters/jwt"
	"github.com/gabrielrauch/reconcile-auth-service/internal/adapters/postgres"
	"github.com/gabrielrauch/reconcile-auth-service/internal/app"
	"github.com/gabrielrauch/reconcile-auth-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	cfg := config.LoadConfig()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	db := postgres.NewPostgres(cfg.DatabaseURL)
	repo := postgres.NewUserRepository(db)
	tokenProvider := jwt.NewTokenProvider(cfg.JwtSecret)

	authApp := app.NewAuthService(repo, tokenProvider, logger)

	r := gin.Default()
	r.Use(http.RequestIDMiddleware(logger))

	http.RunServer(cfg.Port, authApp, logger)
}
