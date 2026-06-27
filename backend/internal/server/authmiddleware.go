package server

import (
	"net/http"

	"github.com/Jacobgtd/hex-stats/backend/internal/authn"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (s *Server) authMiddleware(log zerolog.Logger, permissions authn.Permissions) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := parseBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.clients.AuthnClient.DecipherToken(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		if !user.IsAuthorized(permissions) {
			c.AbortWithStatusJSON(
				http.StatusForbidden, gin.H{
					"error": "forbidden",
				},
			)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
