package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// deriveKey generates a symmetric key from the identity using HMAC-SHA256
func deriveKey(identity string, masterKey []byte) ([]byte, error) {
	h := hmac.New(sha256.New, masterKey)
	_, err := h.Write([]byte(identity))
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}
	return h.Sum(nil)[:32], nil // Use only the first 32 bytes for AES-256
}

// encrypt encrypts plaintext using AES-GCM with the derived key
func encrypt(plaintext string, key []byte) (string, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nonce, nil
}

// decrypt decrypts ciphertext using AES-GCM with the derived key
func decrypt(ciphertextHex string, nonce, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}
	return string(plaintext), nil
}

func main() {
	// Identity and master key
	identity := "user@example.com"
	masterKey := []byte("a very secret master key of 32 bytes...") // 32 bytes for AES-256

	// Derive key based on identity
	key, err := deriveKey(identity, masterKey)
	if err != nil {
		fmt.Println("Error deriving key:", err)
		return
	}

	// Encrypt the message
	plaintext := "Hello, IBSE world!"
	ciphertext, nonce, err := encrypt(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	fmt.Println("Encrypted text:", ciphertext)
	fmt.Println("Nonce:", hex.EncodeToString(nonce))

	// Decrypt the message
	decryptedText, err := decrypt(ciphertext, nonce, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}

	fmt.Println("Decrypted text:", decryptedText)
}
