package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	CACertFilePath = "../client/ca.pem"
	CertFilePath   = "../client/client-certificate.pem"
	KeyFilePath    = "../client/client-key.pem"
)

func main() {
	c, _ := NewCryptoStruct(3)
	encryptedData, err := c.Encrypt([]byte("Vishal Parikh"))

	fmt.Printf("Encrypted Value is: %s\n", encryptedData)
	if err != nil {
		panic(err)
	}

	decryptedData, err := c.Decrypt([]byte(encryptedData))
	fmt.Printf("Decrypted Value is: %s\n", decryptedData)
	if err != nil {
		panic(err)
	}
}

// CryptoInterface defines the interface for cryptographic operations, focusing on HTTP API interaction.
type CryptoInterface interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	CallAPI(method string, url string, payload interface{}) ([]byte, error)
}

// CryptoStruct implements the CryptoInterface and provides methods for cryptographic operations and calling HTTP APIs.
type CryptoStruct struct {
	keyID      int
	baseURL    string
	httpClient *http.Client
}

// NewCryptoStruct creates a new CryptoStruct instance with a randomly generated key and nonce, and an HTTP client configured using the provided key ID.
func NewCryptoStruct(id int) (*CryptoStruct, error) {
	// Load the combined certificate and key
	pemFile, err := os.ReadFile("./certificate.pem")
	if err != nil {
		log.Fatalf("Failed to read PEM file: %v", err)
	}

	// Prepare to decode the PEM file
	var certPEMBlock, keyPEMBlock []byte
	var certFound, keyFound bool

	// Decode the PEM file and separate the cert and key
	for {
		block, rest := pem.Decode(pemFile)
		if block == nil {
			break // No more PEM data
		}
		pemFile = rest // Update pemFile to remaining data

		switch block.Type {
		case "CERTIFICATE":
			certPEMBlock = append(certPEMBlock, pem.EncodeToMemory(block)...)
			certFound = true
		case "RSA PRIVATE KEY":
			keyPEMBlock = append(keyPEMBlock, pem.EncodeToMemory(block)...)
			keyFound = true
		}
	}

	if !certFound || !keyFound {
		log.Fatal("Failed to find both certificate and private key in PEM file")
	}

	// Configure the client to trust TLS server certs issued by a CA.
	// caCertPool := x509.NewCertPool()
	caCertPool, err := x509.SystemCertPool()

	// // Load the CA certificate
	// caCert, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCert, err := os.ReadFile(CACertFilePath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("Failed to append CA certificate")
	}

	// Load the client certificate and private key
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// If we have 2 seperate files than load them like this.
	// cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &CryptoStruct{
		keyID:      id,
		baseURL:    "https://localhost:8000/api/transit",
		httpClient: client,
	}, nil
}

func (c *CryptoStruct) Encrypt(plaintext []byte) ([]byte, error) {
	payload := make(map[string]interface{})
	payload["key_id"] = c.keyID
	payload["plaintext"] = string(plaintext)

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

func (c *CryptoStruct) Decrypt(ciphertext []byte) ([]byte, error) {
	payload := make(map[string]interface{})
	payload["key_id"] = c.keyID
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

	return []byte(response.Plaintext), nil
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
