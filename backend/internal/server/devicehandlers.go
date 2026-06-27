package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewDeviceResponse struct {
	ID     int    `json:"id"`
	Secret string `json:"secret"`
}

func (s *Server) newDeviceHandler(ctx *gin.Context) {

	secret, err := newSecret()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to generate device secret")
		ctx.JSON(500, gin.H{
			"error": "failed to generate device secret",
		})
		return
	}

	newDevice, err := s.clients.DBClient.CreateDevice(ctx, secret)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create device")
		ctx.JSON(500, gin.H{
			"error": "failed to create device",
		})
		return
	}

	ctx.JSON(http.StatusCreated, NewDeviceResponse{
		ID:     newDevice.ID,
		Secret: secret,
	})
}
