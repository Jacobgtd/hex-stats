package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func loggerMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("auth-type", c.GetString("auth-type")).
			Str("auth-id", c.GetString("auth-id")).
			Msg("request")
	}
}
