package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type AuthClient struct {
	config *AuthClientConfig
	logger zerolog.Logger
}

func NewAuthClient(logger zerolog.Logger, config *AuthClientConfig) *AuthClient {
	return &AuthClient{
		logger: logger,
		config: config,
	}
}

func (c *AuthClient) GenerateToken(user User) (string, error) {

	now := time.Now()
	claims := jwtClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.config.ExpiryTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    c.config.ApiName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(c.config.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil

}

func (c *AuthClient) DecipherToken(tokenStr string) (*User, error) {

	claims := jwtClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return c.config.PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims.User, nil
}
