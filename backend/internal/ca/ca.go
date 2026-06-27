package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/rs/zerolog"
)

type CAClient struct {
	config *CAConfig
	logger zerolog.Logger
}

func NewCAClient(logger zerolog.Logger, config *CAConfig) *CAClient {
	return &CAClient{
		config: config,
		logger: logger,
	}
}

func (c *CAClient) GenerateCertificate(commonName string) ([]byte, []byte, error) {
	// gen pk
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	//crt template
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// generate with ca
	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		c.config.crt,
		&clientKey.PublicKey,
		c.config.key,
	)
	if err != nil {
		return nil, nil, err
	}

	// 4. Encode CERT to PEM bytes
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	// 5. Encode KEY to PEM bytes
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
	})

	return certPEM, keyPEM, nil
}

func (c *CAClient) VerifyCertificate(cert *x509.Certificate) error {

	roots := x509.NewCertPool()
	roots.AddCert(c.config.crt)

	_, err := cert.Verify(x509.VerifyOptions{
		Roots: roots,
	})
	if err != nil {
		return err
	}
	return nil
}
