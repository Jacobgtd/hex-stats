package server

import (
	"context"
	"crypto/sha256"
	"net/http"
	"strconv"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type authenticationResponse struct {
	AuthToken string `json:"auth_token"`
}

func (s *Server) authGithub(c *gin.Context) {

	ghToken, err := parseBearerToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, authErr := s.clients.GithubClient.IsAdmin(timeoutCtx, ghToken)
	if authErr != nil {
		c.JSON(authErr.Code, gin.H{
			"error": authErr.Error.Error(),
		})
		return
	}

	token, err := s.clients.AuthClient.GenerateToken(auth.User{
		Type:        auth.UserHuman,
		Id:          user,
		Permissions: auth.PermissionsAdmin,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK, authenticationResponse{
			AuthToken: token,
		},
	)
}

type authDeviceRequest struct {
	Id      uint   `json:"id"`
	Passkey string `json:"passkey"`
}

func (s *Server) authDevice(c *gin.Context) {

	request := authDeviceRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
	}

	device, err := s.clients.DBClient.GetDeviceByID(c, int(request.Id))
	if err != nil {
		s.logger.Error().Err(err).Int("deviceId", int(request.Id)).Msg("failed to get device")
		c.JSON(500, gin.H{
			"error": "failed to get device",
		})
		return
	}

	if device == nil || device.Blacklisted {
		s.logger.Warn().Int("deviceId", int(request.Id)).Msg("device not found")
		c.JSON(404, gin.H{
			"error": "device not found",
		})
		return
	}

	passkeyHash := sha256.Sum256([]byte(request.Passkey))
	if passkeyHash != [32]byte(device.SecretHash) {
		s.logger.Warn().Int("deviceId", device.ID).Msg("authentication failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "passkey is invalid",
		})
	}

	if err := s.clients.DBClient.ActivateDevice(c, int(request.Id)); err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			},
		)
		return
	}

	token, err := s.clients.AuthClient.GenerateToken(auth.User{
		Type:        auth.UserDevice,
		Id:          strconv.Itoa(device.ID),
		Permissions: auth.PermissionsDefault,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK, authenticationResponse{
			AuthToken: token,
		},
	)
}

func (s *Server) checkAuth(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
