package server

import (
	"context"
	"crypto/sha256"
	"net/http"
	"time"

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

type EnrollRequest struct {
	Passkey string `json:"passkey" binding:"required"`
	CSR     string `json:"csr" binding:"required"`
}

type EnrollResponse struct {
	Certificate string `json:"certificate"`
}

func (s *Server) generateDeviceCertificate(ctx *gin.Context) {
	var req EnrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error().Err(err).Msg("invalid enroll request")
		ctx.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	deviceIdRaw := ctx.Param("deviceId")

	deviceId, err := parseDeviceID(deviceIdRaw)
	if err != nil {
		s.logger.Error().Err(err).Msg("invalid device id")
		ctx.JSON(400, gin.H{
			"error": "invalid device id",
		})
		return
	}

	device, err := s.clients.DBClient.GetDeviceByID(ctx, deviceId)
	if err != nil {
		s.logger.Error().Err(err).Int("deviceId", deviceId).Msg("failed to get device")
		ctx.JSON(500, gin.H{
			"error": "failed to get device",
		})
		return
	}

	if device == nil || device.Blacklisted {
		s.logger.Warn().Int("deviceId", deviceId).Msg("device not found")
		ctx.JSON(404, gin.H{
			"error": "device not found",
		})
		return
	}

	passkeyHash := sha256.Sum256([]byte(req.Passkey))
	if passkeyHash != [32]byte(device.SecretHash) {
		s.logger.Warn().Int("deviceId", deviceId).Msg("authentication failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "passkey is invalid",
		})
	}

	csr, err := parseCSR(req.CSR)
	if err != nil {
		s.logger.Error().Err(err).Msg("invalid CSR")
		ctx.JSON(400, gin.H{
			"error": "invalid CSR",
		})
		return
	}

	certPEM, certErr := s.clients.CAClient.GenerateCertificate(csr, deviceId)
	if certErr != nil {
		s.logger.Error().Err(certErr.Error).Msg("failed to generate certificate")
		ctx.JSON(certErr.Code, gin.H{
			"error": certErr.Error.Error(),
		})
		return
	}

	ctxWTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = s.clients.DBClient.ActivateDevice(ctxWTimeout, deviceId)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to generate certificate")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate certificate",
		})
		return
	}

	ctx.JSON(200, EnrollResponse{
		Certificate: string(certPEM),
	})

}
