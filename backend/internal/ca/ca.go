package ca

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/common"
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

func (c *CAClient) GenerateCertificate(csr *x509.CertificateRequest, deviceId int) ([]byte, *common.StatusError) {

	if err := csr.CheckSignature(); err != nil {
		c.logger.Error().Err(err).Str("subject", csr.Subject.CommonName).Msg("invalid CSR signature")
		return nil, &common.StatusError{
			Code:  http.StatusBadRequest,
			Error: errors.New("Invalid CSR Signature"),
		}
	}

	log.Println(csr.Subject.CommonName)

	if csr.Subject.CommonName != fmt.Sprintf("device-%d", deviceId) {
		c.logger.Error().Int("deviceId", deviceId).Str("common_name", csr.Subject.CommonName).Msg("CSR common name does not match expected format")
		return nil, &common.StatusError{
			Code:  http.StatusBadRequest,
			Error: errors.New("CSR Common Name must be in format 'device-{id}'"),
		}
	}

	//crt template
	serial := big.NewInt(time.Now().UnixNano())

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: csr.Subject.CommonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// generate with ca
	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		c.config.crt,
		csr.PublicKey,
		c.config.key,
	)
	if err != nil {
		return nil, &common.StatusError{
			Code:  http.StatusInternalServerError,
			Error: errors.New("Failed to create certificate"),
		}
	}

	// 4. Encode CERT to PEM bytes
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	return certPEM, nil
}

func (c *CAClient) VerifyCertificate(cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	roots.AddCert(c.config.crt)

	_, err := cert.Verify(x509.VerifyOptions{
		Roots:       roots,
		CurrentTime: time.Now(),
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	})
	if err != nil {
		c.logger.Error().Err(err).Str("subject", cert.Subject.CommonName).Str("eku", fmt.Sprintf("%v", cert.ExtKeyUsage)).Str("keyUsage", fmt.Sprintf("%v", cert.KeyUsage)).Msg("failed to verify certificate")
		return err
	}
	return nil
}
