package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/authn"
	"github.com/gin-gonic/gin"
)

type authenticationResponse struct {
	AuthToken string `json:"auth_token"`
}

func (s *Server) AuthGithub(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authorization format",
		})
		return
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "empty bearer token",
		})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, authErr := s.clients.GithubClient.IsAdmin(timeoutCtx, parts[1])
	if authErr != nil {
		c.JSON(authErr.Code, gin.H{
			"error": authErr.Error.Error(),
		})
		return
	}

	token, err := s.clients.AuthnClient.GenerateToken(authn.User{
		Type:        authn.UserHuman,
		Id:          user,
		Permissions: authn.PermissionsAdmin,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			},
		)
	}

	c.JSON(
		http.StatusOK, authenticationResponse{
			AuthToken: token,
		},
	)
}

func AuthDevice(c *gin.Context)
