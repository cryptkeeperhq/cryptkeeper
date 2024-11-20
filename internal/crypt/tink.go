package crypt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/mac"
	"github.com/tink-crypto/tink-go/v2/tink"

	"github.com/tink-crypto/tink-go/v2/aead"

	"github.com/cryptkeeperhq/cryptkeeper/config"
)

type Tink struct {
	MasterAead      tink.AEAD
	MasterKeyHandle *keyset.Handle
	Logger          *slog.Logger
}

func NewTinkOps(cfg *config.Config) (*Tink, error) {
	t := Tink{}

	t.MasterKeyHandle, _ = getTinkMasterKey()
	// t.Logger = config.Logger

	var err error
	t.MasterAead, err = aead.New(t.MasterKeyHandle)
	if err != nil {
		panic(err)
	}

	return &t, nil
}

// Generate a new Path Key (DEK)
// This function generates a new AES-256-GCM key handle for encrypting secrets at the path level.
func (t *Tink) GeneratePathKey() (*keyset.Handle, error) {
	// Create a new key handle using the AES256-GCM template, which provides a 256-bit key for AES encryption in GCM mode
	pathKeyHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	return pathKeyHandle, err
}

// Generate a new Path Key (DEK)
// This function generates a new AES-256-GCM key handle for encrypting secrets at the path level.
func (t *Tink) GeneratePathHmacKey() (*keyset.Handle, error) {
	// Create a new key handle using the AES256-GCM template, which provides a 256-bit key for AES encryption in GCM mode
	pathKeyHandle, err := keyset.NewHandle(mac.HMACSHA256Tag256KeyTemplate())
	return pathKeyHandle, err
}

// Serialize and Encrypt the Path Key
// This function serializes the path key handle and encrypts it using the master key (t.MasterAead).
func (t *Tink) EncryptPathKey(pathKey *keyset.Handle) ([]byte, error) {
	var buffer bytes.Buffer
	writer := keyset.NewBinaryWriter(&buffer)

	// Write the path key handle to the buffer, encrypting it with the master AEAD key
	err := pathKey.Write(writer, t.MasterAead)
	if err != nil {
		return nil, err
	}

	// Encrypt the serialized path key using the master AEAD key and return the encrypted key
	return t.MasterAead.Encrypt(buffer.Bytes(), nil)
}

// Decrypt the Serialized Path Key
// This function takes an encrypted path key, decrypts it using the master key, and reconstructs the path key handle.
func (t *Tink) DecryptPathKey(encryptedPathKey []byte) (*keyset.Handle, error) {
	// Decrypt the encrypted path key using the master AEAD key
	serializedPathKey, err := t.MasterAead.Decrypt(encryptedPathKey, nil)
	if err != nil {
		return nil, err
	}

	// Read the decrypted serialized path key from a buffer and reconstruct the key handle
	var buffer bytes.Buffer
	reader := keyset.NewBinaryReader(&buffer)
	buffer.Write(serializedPathKey) // Write decrypted data to buffer

	// Recreate the path key handle from the serialized key
	pathKeyHandle, err := keyset.Read(reader, t.MasterAead)
	return pathKeyHandle, err
}

// Encrypt Secret Value using a Secret Key and Encrypt the Key with Path Key
// This function generates a new secret key for each secret, encrypts the secret value with it, and then encrypts the secret key using the path key.
func (t *Tink) EncryptSecretValue(input string, pathKeyHandle *keyset.Handle) ([]byte, []byte, error) {
	// Create an AEAD primitive from the path key handle
	pathAead, err := aead.New(pathKeyHandle)
	if err != nil {
		return nil, nil, err
	}

	// Generate a new secret key handle to encrypt the secret value
	secretKeyHandle, _ := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	var buffer bytes.Buffer
	writer := keyset.NewBinaryWriter(&buffer)

	// Serialize the secret key handle and write it to buffer, encrypting it with the path key AEAD
	if err := secretKeyHandle.Write(writer, pathAead); err != nil {
		return nil, nil, err
	}

	// Create an AEAD primitive from the secret key handle for encrypting the actual secret value
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the secret value using the secret AEAD
	ciphertext, err := secretAead.Encrypt([]byte(input), nil)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the serialized secret key with the path AEAD for secure storage
	serializedSecretKey := buffer.Bytes()
	encryptedSecretKey, err := pathAead.Encrypt(serializedSecretKey, nil)
	if err != nil {
		return nil, nil, err
	}

	// Return the encrypted secret key and the encrypted secret value
	return encryptedSecretKey, ciphertext, err
}

// Decrypt Secret Value using the Path Key and Encrypted Secret Key
// This function decrypts the encrypted secret key using the path key, reconstructs the secret key handle, and decrypts the secret value.
func (t *Tink) DecryptSecretValue(encryptedDEK, encryptedValue []byte, pathKeyHandle *keyset.Handle) (string, error) {
	// Create an AEAD primitive from the path key handle to decrypt the secret key (DEK)
	pathAead, err := aead.New(pathKeyHandle)
	if err != nil {
		return "", err
	}

	// Decrypt the DEK (Data Encryption Key) using the path AEAD
	decryptedDek, err := pathAead.Decrypt(encryptedDEK, nil)
	if err != nil {
		return "", err
	}

	// Reconstruct the secret key handle from the decrypted DEK
	var buffer bytes.Buffer
	reader := keyset.NewBinaryReader(&buffer)
	buffer.Write(decryptedDek) // Write decrypted DEK into buffer

	// Read the key handle from the buffer to recreate the secret key
	secretKeyHandle, err := keyset.Read(reader, pathAead)
	if err != nil {
		return "", err
	}

	// Create an AEAD primitive from the secret key handle to decrypt the actual secret value
	secretAead, err := aead.New(secretKeyHandle)
	if err != nil {
		return "", err
	}

	// Decrypt the encrypted secret value
	plainText, err := secretAead.Decrypt(encryptedValue, nil)

	// Return the decrypted plaintext secret value
	return string(plainText), err
}

// Read the master key data from the file
func getTinkMasterKey() (*keyset.Handle, error) {
	masterKeyPath := "master.key"
	masterKeyData, err := LoadMasterKeyFromFile(masterKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read master key from file: %v", err)
	}

	// Decode the base64 encoded keyset data
	// decodedKeyset, err := base64.StdEncoding.DecodeString(string(masterKeyData))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decode master key data: %v", err)
	// }

	// Deserialize the decrypted DEK into a Tink keyset
	ksReader := keyset.NewBinaryReader(bytes.NewReader(masterKeyData))

	return insecurecleartextkeyset.Read(ksReader)
}

// Encode the master key as base64 and save to a file
func SaveMasterKeyToFile(masterKey []byte, filePath string) error {
	encodedKey := base64.StdEncoding.EncodeToString(masterKey)
	return os.WriteFile(filePath, []byte(encodedKey), 0600)
}

// Decode the base64-encoded master key from a file
func LoadMasterKeyFromFile(filePath string) ([]byte, error) {
	encodedKey, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read master key from file: %v", err)
	}
	return base64.StdEncoding.DecodeString(string(encodedKey))
}
