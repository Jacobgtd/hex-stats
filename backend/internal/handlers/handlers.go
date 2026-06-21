package handlers

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

type GenerateCertResponse struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}
