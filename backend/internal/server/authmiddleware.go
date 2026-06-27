package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/ca"
	"github.com/Jacobgtd/hex-stats/backend/internal/github"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func adminAuthMiddleware(logger zerolog.Logger, githubClient *github.GithubClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "empty bearer token",
			})
			return
		}

		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, authErr := githubClient.IsAdmin(timeoutCtx, parts[1])
		if authErr != nil {
			c.AbortWithStatusJSON(authErr.Code, gin.H{
				"error": authErr.Error.Error(),
			})
			return
		}

		c.Set("auth-id", user)
		c.Set("auth-type", "github-admin")
		c.Next()
	}
}
func caAuthMiddleware(logger zerolog.Logger, caClient *ca.CAClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		tlsState := c.Request.TLS

		if tlsState == nil || len(tlsState.PeerCertificates) == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "client certificate required",
			})
			return
		}

		cert := tlsState.PeerCertificates[0]

		// verify it is signed by your CA
		err := caClient.VerifyCertificate(cert)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "invalid client certificate",
			})
			return
		}

		c.Set("auth-id", cert.Subject.CommonName)
		c.Set("auth-type", "certificate")
		c.Next()
	}
}
