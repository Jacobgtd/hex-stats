package ca

import "github.com/gin-gonic/gin"

type GenerateCertResponse struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

func (c *CAClient) GenerateCert(ctx *gin.Context) {
	cert, key, err := c.GenerateCertificate("test-client")
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "failed to generate certificate",
		})
		return
	}
	ctx.JSON(200, GenerateCertResponse{
		Certificate: string(cert),
		Key:         string(key),
	})
}
