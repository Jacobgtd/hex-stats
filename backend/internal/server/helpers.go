package server

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
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

func parseCSR(pemCSR string) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode([]byte(pemCSR))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}
	return x509.ParseCertificateRequest(block.Bytes)
}

// parseDeviceID extracts the numeric ID from a device identifier like "device-12"
func parseDeviceID(deviceID string) (int, error) {
	parts := strings.Split(deviceID, "-")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid device ID format: expected 'device-<id>', got '%s'", deviceID)
	}

	if parts[0] != "device" {
		return 0, fmt.Errorf("invalid device ID format: expected prefix 'device', got '%s'", parts[0])
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid device ID: '%s' is not a valid number", parts[1])
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid device ID: must be positive, got %d", id)
	}

	return id, nil
}
