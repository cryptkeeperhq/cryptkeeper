package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)

var (
	// Paths to the .p12 file and CA certificate
	p12Path     = "./examples/client/client.p12"
	p12Password = "password"
	caPath      = "./examples/client/ca.pem"
)

// CryptoInterface defines the interface for cryptographic operations, focusing on HTTP API interaction.
type CryptoInterface interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	CallAPI(method string, url string, payload interface{}) ([]byte, error)
}

// CryptoStruct implements the CryptoInterface and provides methods for cryptographic operations and calling HTTP APIs.
type CryptoStruct struct {
	baseURL    string
	httpClient *http.Client
}

// Utility function to load the PKCS12 file downloaded from CryptKeer
func loadPKCS12(p12Path, password string) (tls.Certificate, error) {
	// Read the .p12 file
	p12Data, err := os.ReadFile(p12Path)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to read p12 file: %w", err)
	}

	// Decode the .p12 file to extract the certificate and private key
	privateKey, certificate, caCerts, err := pkcs12.DecodeChain(p12Data, password)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to decode p12: %w", err)
	}

	// Convert the private key to PEM format
	privateKeyPEM, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyBlock := pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyPEM}

	// Convert the certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw})

	// Bundle the certificate and private key into a tls.Certificate
	tlsCert, err := tls.X509KeyPair(certPEM, pem.EncodeToMemory(&privateKeyBlock))
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to create TLS certificate: %w", err)
	}

	// Optional: Print additional certificates in the chain
	for i, ca := range caCerts {
		fmt.Printf("CA Certificate %d:\n%s\n", i+1, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Raw}))
	}

	return tlsCert, nil
}

// Utility function to load CA
func loadCA(caPath string) (*x509.CertPool, error) {
	// Read the CA certificate
	caData, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file: %w", err)
	}

	// Add the CA certificate to a CertPool
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caData) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return caCertPool, nil
}

// Utility function to verify the certificate against the CA
func verifyCertificate(tlsCert tls.Certificate, certPool *x509.CertPool) error {
	// Parse the first certificate (the leaf certificate)
	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse x509 certificate: %w", err)
	}

	// Verify the certificate against the CA
	opts := x509.VerifyOptions{
		Roots: certPool,
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("certificate verification failed: %w", err)
	}

	return nil
}

func main() {
	// Load the client certificate and key from the .p12 file
	clientCert, err := loadPKCS12(p12Path, p12Password)
	if err != nil {
		fmt.Printf("Error loading PKCS#12 file: %v\n", err)
		return
	}

	// // Load the CA certificate
	// caCertPool, err := loadCA(caPath)
	// if err != nil {
	// 	fmt.Printf("Error loading CA certificate: %v\n", err)
	// 	return
	// }

	// // Verify it the certificate is valid
	// if err := verifyCertificate(clientCert, caCertPool); err != nil {
	// 	fmt.Printf("Certificate verification failed: %v\n", err)
	// }

	// Configure the TLS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		// RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Should be false for production, but we need it true due to self signed certificate
	}

	// Create an HTTP client with the TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	c := CryptoStruct{
		baseURL:    "https://localhost:8000/api/transit",
		httpClient: client,
	}

	// c.GetKeys()

	keyID := "a5cb8616-332f-4759-879d-6ccfa9d39776"
	encryptedData, err := c.Encrypt(keyID, []byte("Vishal Parikh"))

	fmt.Printf("Encrypted Value is: %s\n", encryptedData)
	if err != nil {
		panic(err)
	}

	decryptedData, err := c.Decrypt(keyID, []byte(encryptedData))
	fmt.Printf("Decrypted Value is: %s\n", decryptedData)
	if err != nil {
		panic(err)
	}

}

func (c *CryptoStruct) GetKeys() ([]byte, error) {
	// Make an mTLS API call
	body, err := c.CallAPI("GET", fmt.Sprintf("%s/keys", c.baseURL), nil)
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", body)

	return nil, nil

}
func (c *CryptoStruct) Encrypt(keyID string, plaintext []byte) ([]byte, error) {
	payload := make(map[string]interface{})
	payload["key_id"] = keyID
	payload["plaintext"] = base64.StdEncoding.EncodeToString(plaintext)

	body, err := c.CallAPI("POST", fmt.Sprintf("%s/encrypt", c.baseURL), payload)

	if err != nil {
		return nil, err
	}

	var response struct {
		Ciphertext string `json:"ciphertext"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return []byte(response.Ciphertext), nil
}

func (c *CryptoStruct) Decrypt(keyID string, ciphertext []byte) ([]byte, error) {
	payload := make(map[string]interface{})
	payload["key_id"] = keyID
	payload["ciphertext"] = string(ciphertext)

	body, err := c.CallAPI("POST", fmt.Sprintf("%s/decrypt", c.baseURL), payload)

	if err != nil {
		return nil, err
	}

	var response struct {
		Plaintext string `json:"plaintext"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(response.Plaintext)
}

func (c *CryptoStruct) CallAPI(method string, url string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
