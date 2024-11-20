package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/shamir"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"google.golang.org/protobuf/proto"
)

// // GenerateMasterKey generates a new master key
// func GenerateMasterKey() ([]byte, error) {
// 	masterKey := make([]byte, 32)
// 	_, err := rand.Read(masterKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return masterKey, nil
// }

func GenerateMasterKey() ([]byte, error) {

	// Initialize Tink AEAD
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		fmt.Println("Error generating Tink key handle:", err)
		return nil, err
	}

	// Prepare to write the keyset into a memory buffer
	memKeyset := &keyset.MemReaderWriter{}
	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
	if err := insecurecleartextkeyset.Write(kh, memKeyset); err != nil {
		fmt.Println("Failed to write keyset:", err)
		return nil, err
	}

	// Serialize the keyset (stored in memKeyset) to a byte slice
	dekBuf, err := proto.Marshal(memKeyset.Keyset)
	if err != nil {
		fmt.Println("Failed to marshal keyset:", err)
		return nil, err
	}

	// Store the encoded keyset in a file
	err = os.WriteFile("master.key", []byte(base64.StdEncoding.EncodeToString(dekBuf)), 0600)
	if err != nil {
		fmt.Println("Error writing master key to file:", err)
		return nil, err
	}

	fmt.Println("Master key generated and stored successfully!")
	return dekBuf, err
}

// SplitMasterKey splits the master key into shares
func SplitMasterKey(masterKey []byte, totalShares, threshold int) ([][]byte, error) {
	shares, err := shamir.Split(masterKey, totalShares, threshold)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

// CombineShares combines the shares to reconstruct the master key
func CombineShares(shares [][]byte) ([]byte, error) {
	masterKey, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

func main() {
	masterKey, err := GenerateMasterKey()
	if err != nil {
		log.Fatalf("Failed to generate master key: %v", err)
	}

	fmt.Printf("Master Key: %x\n", masterKey)

	totalShares := 5
	threshold := 3

	shares, err := SplitMasterKey(masterKey, totalShares, threshold)
	if err != nil {
		log.Fatalf("Failed to split master key: %v", err)
	}

	fmt.Println("Shares:")
	for i, share := range shares {
		fmt.Printf("Share %d: %x\n", i+1, share)
	}

	// Simulate reconstructing the master key from shares
	reconstructedKey, err := CombineShares(shares[:threshold])
	if err != nil {
		log.Fatalf("Failed to combine shares: %v", err)
	}

	fmt.Printf("Reconstructed Key: %x\n", reconstructedKey)
}
