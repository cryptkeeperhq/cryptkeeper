package transit

import (
	"errors"
	"fmt"

	"github.com/tink-crypto/tink-go/v2/mac"
	"github.com/tink-crypto/tink-go/v2/tink"

	"github.com/cryptkeeperhq/cryptkeeper/internal/fpe"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/daead"
)

// func (h *Handler) getDecryptedKeysetHandle(encryptedKeyData []byte) (*keyset.Handle, error) {
// 	pathAead, err := aead.New(h.PathKeyHandle)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create AEAD: %v", err)
// 	}

// 	decryptedDek, err := pathAead.Decrypt(encryptedKeyData, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decrypt data: %v", err)
// 	}

// 	return keyset.Read(keyset.NewBinaryReader(bytes.NewReader(decryptedDek)), pathAead)
// }

func (h *Handler) getPrimitive() (interface{}, error) {
	// handle, err := h.getDecryptedKeysetHandle(encryptedKeyData)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read keyset: %v", err)
	// }
	handle := h.Handle
	switch h.KeyType {
	case "fpe":
		return fpe.New(handle)
	case "aes256S-iv":
		return daead.New(handle)
	case "aes128-gcm96", "aes256-gcm96", "chacha20-poly1305":
		return aead.New(handle)
	case "hmac":
		return mac.New(handle)

	default:
		return nil, errors.New("invalid key type")
	}
}

func (h *Handler) Encrypt(plainText []byte, associatedData []byte) ([]byte, error) {
	if h.KeyType != "aes256S-iv" && h.KeyType != "aes128-gcm96" && h.KeyType != "aes256-gcm96" && h.KeyType != "chacha20-poly1305" && h.KeyType != "fpe" {
		return nil, errors.New("unsupported key type for encryption")
	}

	a, err := h.getPrimitive()
	if err != nil {
		return nil, fmt.Errorf("error getting primitive: %v", err)
	}

	switch h.KeyType {
	case "fpe":
		cipherText, err := a.(fpe.FPEPrimitive).EncryptFPE(string(plainText))
		if err != nil {
			return nil, fmt.Errorf("error encrypting plaintext: %v", err)
		}

		return []byte(cipherText), nil
	case "aes256S-iv":
		return a.(tink.DeterministicAEAD).EncryptDeterministically(plainText, associatedData)
	default:
		cipherText, err := a.(tink.AEAD).Encrypt(plainText, nil)
		if err != nil {
			return nil, fmt.Errorf("error encrypting plaintext: %v", err)
		}
		return cipherText, nil
	}

}

func (h *Handler) Decrypt(cipherText []byte, associatedData []byte) ([]byte, error) {
	if h.KeyType != "aes256S-iv" && h.KeyType != "aes128-gcm96" && h.KeyType != "aes256-gcm96" && h.KeyType != "chacha20-poly1305" && h.KeyType != "fpe" {
		return nil, errors.New("unsupported key type for decryption")
	}

	a, err := h.getPrimitive()
	if err != nil {
		return nil, fmt.Errorf("error getting primitive: %v", err)
	}

	switch h.KeyType {
	case "fpe":
		plainText, err := a.(fpe.FPEPrimitive).DecryptFPE(string(cipherText))
		if err != nil {
			return nil, fmt.Errorf("error encrypting plaintext: %v", err)
		}

		return []byte(plainText), nil
	case "aes256S-iv":
		return a.(tink.DeterministicAEAD).DecryptDeterministically(cipherText, associatedData)
	default:
		plainText, err := a.(tink.AEAD).Decrypt(cipherText, nil)
		if err != nil {
			return nil, fmt.Errorf("error decrypting ciphertext: %v", err)
		}

		return plainText, nil
	}

}
