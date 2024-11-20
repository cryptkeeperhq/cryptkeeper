package sync

// import (
// 	"context"

// 	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
// 	"github.com/Azure/go-autorest/autorest/azure/auth"
// )

// type AzureSyncer struct {
// 	client       *keyvault.BaseClient
// 	vaultBaseURL string
// }

// func NewAzureSyncer(vaultBaseURL string) (*AzureSyncer, error) {
// 	client := keyvault.New()
// 	authorizer, err := auth.NewAuthorizerFromEnvironment()
// 	if err != nil {
// 		return nil, err
// 	}
// 	client.Authorizer = authorizer
// 	return &AzureSyncer{
// 		client:       &client,
// 		vaultBaseURL: vaultBaseURL,
// 	}, nil
// }

// func (a *AzureSyncer) CreateSecret(path string, key string, value string) error {
// 	_, err := a.client.SetSecret(context.Background(), a.vaultBaseURL, path+"-"+key, keyvault.SecretSetParameters{
// 		Value: &value,
// 	})
// 	return err
// }

// func (a *AzureSyncer) UpdateSecret(path string, key string, value string) error {
// 	return a.CreateSecret(path, key, value)
// }

// func (a *AzureSyncer) DeleteSecret(path string, key string) error {
// 	_, err := a.client.DeleteSecret(context.Background(), a.vaultBaseURL, path+"-"+key)
// 	return err
// }
