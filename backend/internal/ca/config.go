package ca

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/Jacobgtd/hex-stats/backend/internal/configpack"
)

type CAConfig struct {
	crt     *x509.Certificate
	crtPath string
	key     *rsa.PrivateKey
	keyPath string
}

func LoadCAConfig() (*CAConfig, error) {
	crt, err := configpack.LoadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	crtPath, err := configpack.GetPath("ca.crt")
	if err != nil {
		return nil, err
	}

	key, err := configpack.LoadFile("ca.key")
	if err != nil {
		return nil, err
	}

	keyPath, err := configpack.GetPath("ca.key")
	if err != nil {
		return nil, err
	}

	caCertBlock, _ := pem.Decode([]byte(crt))
	if caCertBlock == nil {
		return nil, fmt.Errorf("failed to decode ca.crt")
	}

	caKeyBlock, _ := pem.Decode([]byte(key))
	if caKeyBlock == nil {
		return nil, fmt.Errorf("failed to decode ca.key")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ca.crt: %w", err)
	}

	caKeyInterface, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ca.key: %w", err)
	}

	caKey, ok := caKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("ca.key is not an RSA private key")
	}

	return &CAConfig{
		crt:     caCert,
		crtPath: crtPath,
		key:     caKey,
		keyPath: keyPath,
	}, nil
}
