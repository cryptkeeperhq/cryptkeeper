package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/config"
	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	enginepki "github.com/cryptkeeperhq/cryptkeeper/internal/engine/pki"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
)

var (
	CertFilePath = "../server/server-certificate.pem"
	KeyFilePath  = "../server/server-key.pem"
)

func main() {
	// Load Config
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	// Initialize Postgres DB
	db := db.Init(config)

	pathID := 4

	var subCA models.SubCA

	err = db.Model(&subCA).Where("path_id = ?", pathID).First()
	if err != nil {
		panic(err)
	}

	subCert, err := x509.ParseCertificate(subCA.SubCACert)
	if err != nil {
		panic(err)
	}

	subKey, err := x509.ParsePKCS1PrivateKey(subCA.SubCAKey)
	if err != nil {
		panic(err)
	}

	dnsNames := []string{}
	ipAddresses := []string{}
	emailAddresses := []string{}
	expiresAt := time.Now().AddDate(0, 0, 365)

	entityCert, entityKey, err := enginepki.GenerateEndEntityCertificateNew(subCert, subKey, "localhost", "localhost", dnsNames, ipAddresses, emailAddresses, expiresAt)
	if err != nil {
		panic(err)
	}

	// Private Key: This is the private part of the certificate that should be kept secure. It is used to sign data and decrypt information encrypted with the public certificate.
	// Private Key: Users use the private key to sign data or decrypt information. For example, in a web server scenario, the server uses the private key to decrypt incoming HTTPS traffic.
	certPEM := utils.ConvertCertificateToPEM(entityCert.Raw)
	privateKeyPEM, err := utils.ConvertPrivateKeyToPEM(entityKey)
	if err != nil {
		panic(err)
	}

	// // Convert the subCert to PEM format
	// // var pemFile []byte
	// pemBlock := &pem.Block{
	// 	Type:  "CERTIFICATE",
	// 	Bytes: subCert.Raw,
	// }

	// Save the certificate to a PEM file
	// certFilePath := "./certificate.pem" // Set your desired path
	// keyFilePath := "./key.pem"          // Set your desired path

	// Save the certificate to a file
	certFile, err := os.Create(CertFilePath)
	if err != nil {
		panic(err)
	}
	defer certFile.Close()

	certFile.Write([]byte(certPEM))
	// if err := pem.Encode(certFile, certPEM); err != nil {
	// 	panic(err)
	// }

	// // Assuming you have a way to obtain the private key (subCA.SubCAKey)
	// // Convert the private key to PEM format (this is a placeholder; implement as needed)
	// keyPEMBlock := &pem.Block{
	// 	Type:  "PRIVATE KEY",
	// 	Bytes: []byte(subCA.SubCAKey), // Replace this with actual key bytes
	// }

	// Save the private key to a file
	keyFile, err := os.Create(KeyFilePath)
	if err != nil {
		panic(err)
	}
	defer keyFile.Close()

	keyFile.Write([]byte(privateKeyPEM))
	// if err := pem.Encode(keyFile, keyPEMBlock); err != nil {
	// 	panic(err)
	// }

	fmt.Printf("Certificate and key have been saved to %s and %s\n", CertFilePath, KeyFilePath)

	// pemFile = pem.EncodeToMemory(pemBlock)

	// entityCert, entityKey, err := enginepki.GenerateEndEntityCertificate(subCert, subKey, request.Key, request.ExpiresAt)
	// if err != nil {
	// 	return models.Secret{}, errors.New("failed to generate end-entity certificate")
	// }

	// // Private Key: This is the private part of the certificate that should be kept secure. It is used to sign data and decrypt information encrypted with the public certificate.
	// // Private Key: Users use the private key to sign data or decrypt information. For example, in a web server scenario, the server uses the private key to decrypt incoming HTTPS traffic.
	// certPEM := utils.ConvertCertificateToPEM(entityCert.Raw)
	// privateKeyPEM, err := utils.ConvertPrivateKeyToPEM(entityKey)
	// if err != nil {
	// 	return models.Secret{}, err
	// }

}
