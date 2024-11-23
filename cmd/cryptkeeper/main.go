package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/cryptkeeperhq/cryptkeeper/internal/handlers"
	pb "github.com/cryptkeeperhq/cryptkeeper/proto/transit"
	ghandlers "github.com/gorilla/handlers"
	"google.golang.org/grpc"
)

type TransitServer struct {
	pb.UnimplementedTransitServiceServer
	handler *handlers.Handler
}

// Implement the Encrypt method
func (s *TransitServer) Encrypt(ctx context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	// ciphertext, _ := s.handler.GetEncryptedValue(req.KeyId, string(req.Plaintext))
	ciphertext := []byte("return")
	return &pb.EncryptResponse{Ciphertext: []byte(base64.StdEncoding.EncodeToString(ciphertext))}, nil
}

func main() {

	// Load Config
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	h := handlers.Init(config)
	router := h.NewHandler()
	headersOk := ghandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := ghandlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := ghandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Start HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: ghandlers.CORS(originsOk, headersOk, methodsOk)(router),
	}

	// http.Handle("/metrics", promhttp.Handler())

	go func() {
		if config.TLS.Enabled {
			// Load CA certificate
			caCert, err := os.ReadFile(config.TLS.CaFile)
			if err != nil {
				log.Fatalf("Failed to read CA certificate: %v", err)
			}

			// Create a new certificate pool and append the CA certificate
			caCertPool := x509.NewCertPool()
			if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
				log.Fatal("Failed to append CA certificate")
			}

			cert, err := tls.LoadX509KeyPair(config.TLS.CertFile, config.TLS.KeyFile)
			if err != nil {
				log.Fatal(err)
			}
			// Configure the TLS server
			server.TLSConfig = &tls.Config{
				MinVersion:               tls.VersionTLS12,
				ClientAuth:               tls.VerifyClientCertIfGiven,
				ClientCAs:                caCertPool,
				PreferServerCipherSuites: true,
				InsecureSkipVerify:       false,
				Certificates:             []tls.Certificate{cert},
			}

			config.Logger.Info(fmt.Sprintf("Server running on %s:%d with HTTPS", config.Server.Host, config.Server.Port))
			if err := server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s:%d: %v", config.Server.Host, config.Server.Port, err)
			}
		} else {
			config.Logger.Info(fmt.Sprintf("Server running on %s:%d", config.Server.Host, config.Server.Port))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s:%d: %v", config.Server.Host, config.Server.Port, err)
			}
		}
	}()

	// gRPC Server setup
	grpcServer := grpc.NewServer()
	pb.RegisterTransitServiceServer(grpcServer, &TransitServer{handler: h})

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
		if err != nil {
			log.Fatalf("Failed to listen on port %d: %v", 50051, err)
		}

		log.Printf("gRPC Server running on port %d", 50051)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Wait for exit signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	// Graceful shutdown
	config.Logger.Debug("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// h.Consumer.Close()
	// h.Producer.Close()

}
