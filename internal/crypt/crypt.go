package crypt

import (
	"errors"

	"github.com/tink-crypto/tink-go/v2/keyset"

	"github.com/cryptkeeperhq/cryptkeeper/config"
)

var ErrUnsupportedOption = errors.New("unsupported option")

type CryptographicOperations interface {
	GeneratePathKey() (*keyset.Handle, error)
	EncryptPathKey(pathKey *keyset.Handle) ([]byte, error)
	DecryptPathKey(encryptedPathKey []byte) (*keyset.Handle, error)
	EncryptSecretValue(input string, pathKeyHandle *keyset.Handle) ([]byte, []byte, error)
	DecryptSecretValue(encryptedDEK, encryptedValue []byte, pathKeyHandle *keyset.Handle) (string, error)
}

func New(config *config.Config, option string) (CryptographicOperations, error) {
	switch option {
	case "hardware":
		return NewHSM(config)
	default:
		return NewTinkOps(config)
	}
}
