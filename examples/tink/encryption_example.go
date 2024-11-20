package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
	"google.golang.org/protobuf/proto"
)

// This is for illustration purposes only. Insecurecleartextkeyset should not be used in production!
// Use a Key Management Service (KMS) or Tink's key management features for secure key storage.
const useInsecureCleartextKeyset = false // Set to true to use insecurecleartextkeyset (not recommended)

func EncryptWithTink(plaintext []byte, pathKeyHandle *keyset.Handle) ([]byte, []byte, error) {

	pathAead, _ := aead.New(pathKeyHandle)

	// Generate Key Handle
	secretKeyHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	var buffer bytes.Buffer
	writer := keyset.NewBinaryWriter(&buffer)
	if err := secretKeyHandle.Write(writer, pathAead); err != nil {
		return nil, nil, err
	}

	// New AEAD primitive
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the value
	ciphertext, err := secretAead.Encrypt(plaintext, nil)

	// Lastly encrypt the DEK using Path Key
	serializedSecretKey := buffer.Bytes()
	encryptedSecretKey, err := pathAead.Encrypt(serializedSecretKey, nil)
	if err != nil {
		return nil, nil, err
	}

	return encryptedSecretKey, ciphertext, err
}

func DecryptWithTink(encryptedDek, ciphertext []byte, pathKeyHandle *keyset.Handle) ([]byte, error) {

	pathAead, _ := aead.New(pathKeyHandle)

	// Firstly decrypt the DEK using Path Key
	decryptedDek, _ := pathAead.Decrypt(encryptedDek, nil)

	// Read the decrypted DEK into Key Handle
	var buffer bytes.Buffer
	reader := keyset.NewBinaryReader(&buffer)
	buffer.Write(decryptedDek)
	secretKeyHandle, err := keyset.Read(reader, pathAead)

	// New AEAD primitive
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return nil, err
	}

	return secretAead.Decrypt(ciphertext, nil)
}

func EncryptSecretValue(input string, pathKeyHandle *keyset.Handle) ([]byte, []byte, error) {
	// Encrypt the secret value using the DEK
	encryptedDEK, encryptedValue, err := EncryptWithTink([]byte(input), pathKeyHandle)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the DEK using the Path Key
	// pathAead, err := aead.New(pathKeyHandle)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// encryptedDEK, err := pathAead.Encrypt(dek, nil)
	// if err != nil {
	// 	return nil, nil, err
	// }

	return encryptedDEK, encryptedValue, nil
}

func DecryptSecretValue(encryptedDEK, encryptedValue []byte, pathKeyHandle *keyset.Handle) (string, error) {

	// Decrypt the DEK using the Path Key
	// pathAead, err := aead.New(pathKeyHandle)
	// if err != nil {
	// 	return "", err
	// }

	// decryptedDek, err := pathAead.Decrypt(encryptedDEK, nil)
	// if err != nil {
	// 	return "", err
	// }

	// Decrypt the secret value using the DEK
	decryptedValue, err := DecryptWithTink(encryptedDEK, encryptedValue, pathKeyHandle)
	if err != nil {
		return "", err
	}

	return string(decryptedValue), nil
}

func GenerateAndEncryptPathKey(masterKeyHandle *keyset.Handle) ([]byte, *keyset.Handle, error) {
	// Generate a new Path Key (DEK)
	pathKeyHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		return nil, nil, err
	}

	// Serialize the Path Key
	var buffer bytes.Buffer
	writer := keyset.NewBinaryWriter(&buffer)
	masterAead, err := aead.New(masterKeyHandle)
	err = pathKeyHandle.Write(writer, masterAead)
	if err != nil {
		return nil, nil, err
	}

	serializedPathKey := buffer.Bytes()

	encryptedPathKey, err := masterAead.Encrypt(serializedPathKey, nil)
	if err != nil {
		return nil, nil, err
	}

	return encryptedPathKey, pathKeyHandle, nil
}

func DecryptPathKey(encryptedPathKey []byte, masterKeyHandle *keyset.Handle) (*keyset.Handle, error) {
	masterAead, err := aead.New(masterKeyHandle)
	if err != nil {
		return nil, err
	}

	serializedPathKey, err := masterAead.Decrypt(encryptedPathKey, nil)
	if err != nil {
		return nil, err
	}

	// Serialize the Path Key
	var buffer bytes.Buffer
	reader := keyset.NewBinaryReader(&buffer)
	buffer.Write(serializedPathKey) // Write decrypted data to buffer
	pathKeyHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	pathKeyHandle, err = keyset.Read(reader, masterAead)
	return pathKeyHandle, err
}

func GetMasterKey() (*keyset.Handle, error) {
	// Read the master key data from the file
	masterKeyPath := "master.key"
	masterKeyData, err := os.ReadFile(masterKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read master key from file: %v", err)
	}

	// Decode the base64 encoded keyset data
	decodedKeyset, err := base64.StdEncoding.DecodeString(string(masterKeyData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode master key data: %v", err)
	}

	// Deserialize the decrypted DEK into a Tink keyset
	ksReader := keyset.NewBinaryReader(bytes.NewReader(decodedKeyset))

	kh, err := insecurecleartextkeyset.Read(ksReader)

	return kh, nil
}

func main() {
	// GenerateMasterKey()

	masterKeyHandle, err := GetMasterKey()
	if err != nil {
		panic(err)
	}

	fmt.Println("Got master key from file", masterKeyHandle)

	// Generate and encrypt a path key
	encryptedPathKey, _, err := GenerateAndEncryptPathKey(masterKeyHandle)
	if err != nil {
		log.Fatalf("failed to generate and encrypt path key: %v", err)
	}
	// writeKeyHandleToFile(pathKeyHandle, "path.key")

	// Decrypt the path key
	decryptedPathKeyHandle, err := DecryptPathKey(encryptedPathKey, masterKeyHandle)
	if err != nil {
		log.Fatalf("failed to decrypt path key: %v", err)
	}

	// Encrypt a secret value
	input := "1234567890"
	fmt.Printf("Input Value: %s\n", input)

	encryptedDEK, encryptedValue, err := EncryptSecretValue(input, decryptedPathKeyHandle)
	if err != nil {
		log.Fatalf("failed to encrypt secret value: %v", err)
	}

	fmt.Printf("Encrypted DEK: %s\n", base64.StdEncoding.EncodeToString(encryptedDEK))
	fmt.Printf("Encrypted value: %s\n", base64.StdEncoding.EncodeToString(encryptedValue))

	// Decrypt the secret value
	decryptedValue, err := DecryptSecretValue(encryptedDEK, encryptedValue, decryptedPathKeyHandle)
	if err != nil {
		log.Fatalf("failed to decrypt secret value: %v", err)
	}

	fmt.Printf("Decrypted value: %s\n", decryptedValue)
}

func GenerateMasterKey() {

	// Initialize Tink AEAD
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		fmt.Println("Error generating Tink key handle:", err)
		return
	}

	// Prepare to write the keyset into a memory buffer
	memKeyset := &keyset.MemReaderWriter{}
	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
	if err := insecurecleartextkeyset.Write(kh, memKeyset); err != nil {
		fmt.Println("Failed to write keyset:", err)
		return
	}

	// Serialize the keyset (stored in memKeyset) to a byte slice
	dekBuf, err := proto.Marshal(memKeyset.Keyset)
	if err != nil {
		fmt.Println("Failed to marshal keyset:", err)
		return
	}

	// Store the encoded keyset in a file
	err = os.WriteFile("master.key", []byte(base64.StdEncoding.EncodeToString(dekBuf)), 0600)
	if err != nil {
		fmt.Println("Error writing master key to file:", err)
		return
	}

	fmt.Println("Master key generated and stored successfully!")

}

func writeKeyHandleToFile(kh *keyset.Handle, fileNm string) {
	// Prepare to write the keyset into a memory buffer
	memKeyset := &keyset.MemReaderWriter{}
	// Use insecurecleartextkeyset to write the key handle (kh) into the memKeyset buffer
	if err := insecurecleartextkeyset.Write(kh, memKeyset); err != nil {
		fmt.Println("Failed to write keyset:", err)
		return
	}

	// Serialize the keyset (stored in memKeyset) to a byte slice
	dekBuf, err := proto.Marshal(memKeyset.Keyset)
	if err != nil {
		fmt.Println("Failed to marshal keyset:", err)
		return
	}

	// Store the encoded keyset in a file
	err = os.WriteFile(fileNm, []byte(base64.StdEncoding.EncodeToString(dekBuf)), 0600)
	if err != nil {
		fmt.Println("Error writing master key to file:", err)
		return
	}

	fmt.Println("File Written ", fileNm)
}

func readKeyHandleFromFile(fileNm string) (*keyset.Handle, []byte, error) {
	masterKeyData, err := os.ReadFile(fileNm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read master key from file: %v", err)
	}

	// Decode the base64 encoded keyset data
	decodedKeyset, err := base64.StdEncoding.DecodeString(string(masterKeyData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode master key data: %v", err)
	}

	// Deserialize the decrypted DEK into a Tink keyset
	ksReader := keyset.NewBinaryReader(bytes.NewReader(decodedKeyset))

	kh, err := insecurecleartextkeyset.Read(ksReader)

	return kh, decodedKeyset, nil
}
