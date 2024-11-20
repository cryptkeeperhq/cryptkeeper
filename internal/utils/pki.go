package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// ConvertCertificateToPEM converts a DER encoded certificate to PEM format
func ConvertCertificateToPEM(certDER []byte) string {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	return string(certPEM)
}

// ConvertPrivateKeyToPEM converts a private key to PEM format
func ConvertPrivateKeyToPEM(privateKey interface{}) (string, error) {
	var pemBlock *pem.Block
	var err error

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		}
	case *ecdsa.PrivateKey:
		var b []byte
		b, err = x509.MarshalECPrivateKey(key)
		if err != nil {
			return "", err
		}
		pemBlock = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: b,
		}
	default:
		return "", errors.New("unsupported private key type")
	}

	return string(pem.EncodeToMemory(pemBlock)), nil
}
