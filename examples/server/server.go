package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	CACertFilePath = "../client/ca.pem"
	CertFilePath   = "../server/server-certificate.pem"
	KeyFilePath    = "../server/server-key.pem"
)

// TLSAuthMiddleware is a middleware to handle TLS authentication
func TLSAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientCerts := r.TLS.PeerCertificates
		if len(clientCerts) > 0 {
			clientCert := clientCerts[0] // Get the first client certificate
			// Log the subject and serial number of the client certificate
			log.Printf("Client Certificate Subject: %s, Serial Number: %s", clientCert.Subject.String(), clientCert.SerialNumber.String())
		} else {
			log.Println("No client certificate provided")
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load CA certificate
	caCert, err := os.ReadFile(CACertFilePath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	// Create a new certificate pool and append the CA certificate
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("Failed to append CA certificate")
	}

	cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                caCertPool,
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
	}

	r := mux.NewRouter()

	r.Handle("/", TLSAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.RemoteAddr)
	})))

	// Create a server with the TLS configuration
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   r,
	}

	// Start the server
	fmt.Println("Starting server on https://localhost:8443")
	if err := server.ListenAndServeTLS(CertFilePath, KeyFilePath); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
