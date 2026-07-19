package server

import (
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func loggerMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		user := auth.User{}
		value, exists := c.Get("user")
		if exists {
			innerUser, ok := value.(*auth.User)
			if ok {
				user = *innerUser
			}
		}

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("auth-type", string(user.Type)).
			Str("auth-id", string(user.Id)).
			Msg("request")
	}
}
