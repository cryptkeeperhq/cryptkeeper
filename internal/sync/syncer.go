package sync

import (
	"errors"

	"github.com/cryptkeeperhq/cryptkeeper/config"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")

type SecretSyncer interface {
	CreateSecret(path string, key string, value string) error
	UpdateSecret(path string, key string, value string) error
	DeleteSecret(path string, key string) error
}

func NewSyncer(cfg *config.Config) (SecretSyncer, error) {
	switch cfg.Sync.Platform {
	case "vault":
		return NewVaultSyncer(cfg)
	// case "aws":
	// 	return NewAWSSyncer(cfg)
	// case "azure":
	// 	return NewAzureSyncer(cfg)
	default:
		return nil, ErrUnsupportedPlatform
	}
}
