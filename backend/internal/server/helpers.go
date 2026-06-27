package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jacobgtd/hex-stats/backend/internal/authn"
	"github.com/gin-gonic/gin"
)

func newSecret() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	secret := base64.StdEncoding.EncodeToString(b)
	return secret, nil
}

func parseBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}
	return token, nil
}

func getUser(c *gin.Context) (*authn.User, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, fmt.Errorf("could not get user")
	}

	u, ok := user.(*authn.User)
	if !ok {
		return nil, fmt.Errorf("could not get user")
	}
	return u, nil
}

func getIdFromPath(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
