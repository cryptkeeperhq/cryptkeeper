package main

import (
	"context"
	"log"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/proto/transit"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := transit.NewTransitServiceClient(conn)

	// Call the Encrypt method
	req := &transit.EncryptRequest{
		KeyId:     "31a22fe8-a748-49b9-ad74-f17307d4e123",
		Plaintext: []byte("Hello, World!"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Encrypt(ctx, req)
	if err != nil {
		log.Fatalf("could not encrypt: %v", err)
	}

	log.Printf("Ciphertext: %s", string(res.Ciphertext))
}
