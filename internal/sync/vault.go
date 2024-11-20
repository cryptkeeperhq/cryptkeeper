package sync

import (
	"fmt"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/hashicorp/vault/api"
)

type VaultSyncer struct {
	client *api.Client
}

func NewVaultSyncer(cfg *config.Config) (*VaultSyncer, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.Sync.Vault.Address

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Sync.Vault.Token)

	return &VaultSyncer{client: client}, nil
}

func (s *VaultSyncer) CreateSecret(path, key, value string) error {

	// data := map[string]interface{}{
	// 	"data": map[string]interface{}{
	// 		key: value,
	// 	},
	// }

	data := map[string]interface{}{
		"value": value,
	}

	fullPath := fmt.Sprintf("cryptkeeper/%s/%s", path, key)
	_, err := s.client.Logical().Write(fullPath, data)
	return err
}

func (s *VaultSyncer) UpdateSecret(path string, key string, value string) error {
	return nil
}

func (s *VaultSyncer) DeleteSecret(path string, key string) error {
	return nil
}
