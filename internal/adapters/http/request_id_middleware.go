package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type key int

const requestIDKey key = 0

func RequestIDMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()

		ctx := context.WithValue(c.Request.Context(), requestIDKey, reqID)

		c.Request = c.Request.WithContext(ctx)

		logger = logger.With().Str("request_id", reqID).Logger()

		logger.Info().Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Msg("request received")

		c.Header("X-Request-ID", reqID)

		c.Next()
	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
