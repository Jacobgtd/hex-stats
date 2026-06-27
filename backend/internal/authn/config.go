package authn

import (
	"crypto/rsa"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/configpack"
	"github.com/golang-jwt/jwt/v5"
)

type AuthnClientConfig struct {
	ExpiryTime time.Duration
	ApiName    string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadAuthnConfig() (*AuthnClientConfig, error) {
	err := configpack.Load("jwt.config")
	if err != nil {
		return nil, err
	}

	expiryTimeInt := configpack.IntOrDefault("EXPIRY_TIME_SECONDS", 600)

	apiName, err := configpack.String("API_NAME")
	if err != nil {
		return nil, err
	}

	privateKeyStr, err := configpack.LoadFile("jwt_private.pem")
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyStr))
	if err != nil {
		return nil, err
	}

	publicKeyStr, err := configpack.LoadFile("jwt_public.pem")
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyStr))
	if err != nil {
		return nil, err
	}

	return &AuthnClientConfig{
		ExpiryTime: time.Second * time.Duration(expiryTimeInt),
		ApiName:    apiName,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
