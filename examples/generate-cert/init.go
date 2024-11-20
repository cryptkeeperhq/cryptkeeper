package main

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"log/slog"
// 	"os"

// 	"github.com/google/tink/go/aead"
// 	"github.com/google/tink/go/insecurecleartextkeyset"
// 	"github.com/google/tink/go/keyset"
// 	"github.com/google/tink/go/tink"
// 	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
// )

// type Tink struct {
// 	MasterAead      tink.AEAD
// 	MasterKeyHandle *keyset.Handle
// 	Logger          *slog.Logger
// }

// // type TinkPrimitivate interface {
// // 	CreateSecret(path string, key string, value string) error
// // 	UpdateSecret(path string, key string, value string) error
// // 	DeleteSecret(path string, key string) error
// // }

// func getTinkMasterKey() (*keyset.Handle, error) {
// 	// Read the master key data from the file
// 	masterKeyPath := "master.key"
// 	masterKeyData, err := os.ReadFile(masterKeyPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read master key from file: %v", err)
// 	}

// 	// Decode the base64 encoded keyset data
// 	decodedKeyset, err := base64.StdEncoding.DecodeString(string(masterKeyData))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode master key data: %v", err)
// 	}

// 	// Deserialize the decrypted DEK into a Tink keyset
// 	ksReader := keyset.NewBinaryReader(bytes.NewReader(decodedKeyset))

// 	return insecurecleartextkeyset.Read(ksReader)
// }

// func Init(config *config.Config) *Tink {
// 	// fmt.Println("Got master key from file", t.MasterKeyHandle)

// 	// // Generate and encrypt a path key
// 	// encryptedPathKey, _, err := t.GenerateAndEncryptPathKey()
// 	// if err != nil {
// 	// 	log.Fatalf("failed to generate and encrypt path key: %v", err)
// 	// }
// 	// // writeKeyHandleToFile(pathKeyHandle, "path.key")

// 	// // Decrypt the path key
// 	// decryptedPathKeyHandle, err := t.DecryptPathKey(encryptedPathKey)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to decrypt path key: %v", err)
// 	// }

// 	// // Encrypt a secret value
// 	// input := "1234567890"
// 	// fmt.Printf("Input Value: %s\n", input)

// 	// encryptedDEK, encryptedValue, err := t.EncryptSecretValue(input, decryptedPathKeyHandle)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to encrypt secret value: %v", err)
// 	// }

// 	// fmt.Printf("Encrypted DEK: %s\n", base64.StdEncoding.EncodeToString(encryptedDEK))
// 	// fmt.Printf("Encrypted value: %s\n", base64.StdEncoding.EncodeToString(encryptedValue))

// 	// // Decrypt the secret value
// 	// decryptedValue, err := t.DecryptSecretValue(encryptedDEK, encryptedValue, decryptedPathKeyHandle)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to decrypt secret value: %v", err)
// 	// }

// 	// fmt.Printf("Decrypted value: %s\n", decryptedValue)
// 	// GenerateMasterKey()

// 	// share1, _ := hex.DecodeString("942d440f98d7f007df28ac434c70f62ed612fae09573c51d600d042dab7f8a96f3")
// 	// share2, _ := hex.DecodeString("21c0b8c946aace37338722c5620dbbd0ed1a57a8d5dfe777bfb7a79562c2715345")
// 	// share3, _ := hex.DecodeString("e595524751921c0f4ea4dffdd098882e3d981e4817b3c648e91720d02a61961146")

// 	// // During system restart
// 	// shares := [][]byte{share1, share2, share3}
// 	// masterKey, err := CombineShares(shares)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to combine shares: %v", err)
// 	// }

// 	// masterKeyString := base64.StdEncoding.EncodeToString(masterKey)
// 	// fmt.Printf("Reconstructed Key: %x\n", masterKey)
// 	// err = crypto.InitKeyManagement(masterKeyString)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	t := Tink{}

// 	t.MasterKeyHandle, _ = getTinkMasterKey()
// 	t.Logger = config.Logger

// 	var err error
// 	t.MasterAead, err = aead.New(t.MasterKeyHandle)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &t

// }

// // This is for illustration purposes only. Insecurecleartextkeyset should not be used in production!
// // Use a Key Management Service (KMS) or Tink's key management features for secure key storage.
// //  const useInsecureCleartextKeyset = false // Set to true to use insecurecleartextkeyset (not recommended)
// // func GenerateMasterKey() {

// // 	// Initialize Tink AEAD
// // 	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
// // 	if err != nil {
// // 		fmt.Println("Error generating Tink key handle:", err)
// // 		return
// // 	}

// // 	// Prepare to write the keyset into a memory buffer
// // 	memKeyset := &keyset.MemReaderWriter{}
// // 	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
// // 	if err := insecurecleartextkeyset.Write(kh, memKeyset); err != nil {
// // 		fmt.Println("Failed to write keyset:", err)
// // 		return
// // 	}

// // 	// Serialize the keyset (stored in memKeyset) to a byte slice
// // 	dekBuf, err := proto.Marshal(memKeyset.Keyset)
// // 	if err != nil {
// // 		fmt.Println("Failed to marshal keyset:", err)
// // 		return
// // 	}

// // 	// Store the encoded keyset in a file
// // 	err = os.WriteFile("master.key", []byte(base64.StdEncoding.EncodeToString(dekBuf)), 0600)
// // 	if err != nil {
// // 		fmt.Println("Error writing master key to file:", err)
// // 		return
// // 	}

// // 	fmt.Println("Master key generated and stored successfully!")

// // }

// // func writeKeyHandleToFile(kh *keyset.Handle, fileNm string) {
// // 	// Prepare to write the keyset into a memory buffer
// // 	memKeyset := &keyset.MemReaderWriter{}
// // 	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
// // 	if err := insecurecleartextkeyset.Write(kh, memKeyset); err != nil {
// // 		fmt.Println("Failed to write keyset:", err)
// // 		return
// // 	}

// // 	// Serialize the keyset (stored in memKeyset) to a byte slice
// // 	dekBuf, err := proto.Marshal(memKeyset.Keyset)
// // 	if err != nil {
// // 		fmt.Println("Failed to marshal keyset:", err)
// // 		return
// // 	}

// // 	// Store the encoded keyset in a file
// // 	err = os.WriteFile(fileNm, []byte(base64.StdEncoding.EncodeToString(dekBuf)), 0600)
// // 	if err != nil {
// // 		fmt.Println("Error writing master key to file:", err)
// // 		return
// // 	}

// // 	fmt.Println("File Written ", fileNm)
// // }

// // func readKeyHandleFromFile(fileNm string) (*keyset.Handle, []byte, error) {
// // 	masterKeyData, err := os.ReadFile(fileNm)
// // 	if err != nil {
// // 		return nil, nil, fmt.Errorf("failed to read master key from file: %v", err)
// // 	}

// // 	// Decode the base64 encoded keyset data
// // 	decodedKeyset, err := base64.StdEncoding.DecodeString(string(masterKeyData))
// // 	if err != nil {
// // 		return nil, nil, fmt.Errorf("failed to decode master key data: %v", err)
// // 	}

// // 	// Deserialize the decrypted DEK into a Tink keyset
// // 	ksReader := keyset.NewBinaryReader(bytes.NewReader(decodedKeyset))

// // 	kh, err := insecurecleartextkeyset.Read(ksReader)

// // 	return kh, decodedKeyset, nil
// // }
