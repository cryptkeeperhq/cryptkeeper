package sync

// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/secretsmanager"
// )

// type AWSSyncer struct {
// 	client *secretsmanager.SecretsManager
// }

// func NewAWSSyncer(sess *session.Session) *AWSSyncer {
// 	return &AWSSyncer{
// 		client: secretsmanager.New(sess),
// 	}
// }

// func (a *AWSSyncer) CreateSecret(path string, key string, value string) error {
// 	_, err := a.client.CreateSecret(&secretsmanager.CreateSecretInput{
// 		Name:         aws.String(path + "/" + key),
// 		SecretString: aws.String(value),
// 	})
// 	return err
// }

// func (a *AWSSyncer) UpdateSecret(path string, key string, value string) error {
// 	_, err := a.client.UpdateSecret(&secretsmanager.UpdateSecretInput{
// 		SecretId:     aws.String(path + "/" + key),
// 		SecretString: aws.String(value),
// 	})
// 	return err
// }

// func (a *AWSSyncer) DeleteSecret(path string, key string) error {
// 	_, err := a.client.DeleteSecret(&secretsmanager.DeleteSecretInput{
// 		SecretId: aws.String(path + "/" + key),
// 	})
// 	return err
// }
