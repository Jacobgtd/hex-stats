package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
)

type EnrollRequest struct {
	Passkey string `json:"passkey"`
	CSR     string `json:"csr"`
}

type EnrollResponse struct {
	Certificate string `json:"certificate"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <deviceId> <passkey> [backend_url]")
		fmt.Println("Example: go run main.go device-12 mypasskey http://localhost:8443")
		os.Exit(1)
	}

	deviceIDStr := os.Args[1]
	passkey := os.Args[2]
	backendURL := "https://localhost:8080"
	if len(os.Args) > 3 {
		backendURL = os.Args[3]
	}

	// Generate RSA private key
	fmt.Println("=== GENERATING KEYS ===")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Error generating private key: %v\n", err)
		os.Exit(1)
	}

	// Print private key
	privateKeyPEM := privateKeyToPEM(privateKey)
	fmt.Println("\n=== PRIVATE KEY ===")
	fmt.Println(privateKeyPEM)

	// Print public key
	publicKeyPEM := publicKeyToPEM(&privateKey.PublicKey)
	fmt.Println("=== PUBLIC KEY ===")
	fmt.Println(publicKeyPEM)

	// Generate CSR
	fmt.Println("=== GENERATING CSR ===")
	csr, err := generateCSR(privateKey, deviceIDStr)
	if err != nil {
		fmt.Printf("Error generating CSR: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== CSR ===")
	fmt.Println(csr)

	// Send request to backend
	fmt.Printf("\n=== SENDING REQUEST TO BACKEND ===\n")
	fmt.Printf("Endpoint: %s/api/v1/devices/%s/certificate\n", backendURL, deviceIDStr)

	enrollReq := EnrollRequest{
		Passkey: passkey,
		CSR:     csr,
	}

	reqBody, err := json.Marshal(enrollReq)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		os.Exit(1)
	}

	endpoint := fmt.Sprintf("%s/api/v1/devices/%s/certificate", backendURL, deviceIDStr)

	// Create HTTP client that ignores certificate verification (for testing with self-signed certs)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to backend: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("Backend returned error status %d\n", resp.StatusCode)
		fmt.Printf("Response: %s\n", string(body))
		os.Exit(1)
	}

	var enrollResp EnrollResponse
	if err := json.Unmarshal(body, &enrollResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== CERTIFICATE ===")
	fmt.Println(enrollResp.Certificate)
}

func privateKeyToPEM(key *rsa.PrivateKey) string {
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	return string(privPEM)
}

func publicKeyToPEM(key *rsa.PublicKey) string {
	pubBytes, _ := x509.MarshalPKIXPublicKey(key)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	return string(pubPEM)
}

func generateCSR(privateKey *rsa.PrivateKey, deviceId string) (string, error) {
	subj := pkix.Name{
		CommonName: deviceId,
	}

	template := &x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, privateKey)
	if err != nil {
		return "", err
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	return string(csrPEM), nil
}
