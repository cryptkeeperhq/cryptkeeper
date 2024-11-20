package handlers

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	enginepki "github.com/cryptkeeperhq/cryptkeeper/internal/engine/pki"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"software.sslmate.com/src/go-pkcs12"
)

func (h *Handler) DownloadCA(w http.ResponseWriter, r *http.Request) {

	pathID := r.URL.Query().Get("path_id")

	var subCA models.SubCA

	err := h.DB.Model(&subCA).Where("path_id = ?", pathID).First()
	if err != nil {
		http.Error(w, "subCA not found", http.StatusBadRequest)
	}

	var rootCA models.RootCA

	err = h.DB.Model(&rootCA).Where("id = ?", subCA.RootCA).First()
	if err != nil {
		http.Error(w, "rootCA not found", http.StatusBadRequest)
		return
	}

	rootCert, err := x509.ParseCertificate(rootCA.RootCert)
	if err != nil {
		http.Error(w, "subCA parsing failed", http.StatusBadRequest)
		return
	}

	// rootCACertPEM, err := os.ReadFile("root_ca.pem")
	// if err != nil {
	// 	fmt.Printf("Error reading Root CA certificate: %v\n", err)
	// 	return
	// }

	subCert, err := x509.ParseCertificate(subCA.SubCACert)
	if err != nil {
		http.Error(w, "subCA parsing failed", http.StatusBadRequest)
		return
	}

	// Generate CA Bundle
	caBundle, err := enginepki.GenerateCABundle(subCert, rootCert)
	if err != nil {
		http.Error(w, "error generating bundle", http.StatusBadRequest)
		return
	}

	// Save the CA Bundle to a file
	err = os.WriteFile("ca_bundle.pem", caBundle, 0644)
	if err != nil {
		fmt.Printf("Error writing CA bundle to file: %v\n", err)
		http.Error(w, "Error writing CA bundle to file", http.StatusBadRequest)
		return
	}

	// // Convert the subCert to PEM format
	// var pemFile []byte
	// pemBlock := &pem.Block{
	// 	Type:  "CERTIFICATE",
	// 	Bytes: subCert.Raw,
	// }

	// pemFile = pem.EncodeToMemory(pemBlock)

	// Set the response headers and body
	w.Header().Set("Content-Disposition", "attachment; filename=ca_bundle.pem")
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Write(caBundle)

}

func (h *Handler) DownloadCertificate(w http.ResponseWriter, r *http.Request) {
	secretID := r.URL.Query().Get("secret_id")
	if secretID == "" {
		http.Error(w, "Secret ID is required", http.StatusBadRequest)
		return
	}

	var secret models.Secret
	err := h.DB.Model(&secret).Where("id = ?", secretID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the path
	var pathDetails models.Path
	err = h.DB.Model(&pathDetails).Where("id = ?", secret.PathID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// GET the Sub CA from path
	var subCA models.SubCA
	err = h.DB.Model(&subCA).Where("path_id = ?", secret.PathID).First()
	if err != nil {
		http.Error(w, "subCA not found", http.StatusBadRequest)
	}

	subCert, err := x509.ParseCertificate(subCA.SubCACert)
	if err != nil {
		http.Error(w, "subCA parsing failed", http.StatusBadRequest)
	}

	var rootCA models.RootCA
	err = h.DB.Model(&rootCA).Where("id = ?", subCA.RootCA).First()
	if err != nil {
		http.Error(w, "rootCA not found", http.StatusBadRequest)
		return
	}

	rootCert, err := x509.ParseCertificate(rootCA.RootCert)
	if err != nil {
		http.Error(w, "rootCA parsing failed", http.StatusBadRequest)
		return
	}

	// Decrypt the path key
	decryptedPathKeyHandle, err := h.CryptoOps.DecryptPathKey(pathDetails.KeyData)
	if err != nil {
		h.Config.Logger.Error("failed to decrypt path key", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// GET the Private Key Value
	plaintextValue, err := h.CryptoOps.DecryptSecretValue(secret.EncryptedDEK, secret.EncryptedValue, decryptedPathKeyHandle)
	if err != nil {
		h.Config.Logger.Error("failed to get plaintextValue", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decodedPlainTextValue, _ := base64.StdEncoding.DecodeString(plaintextValue)
	privateKey, err := x509.ParsePKCS1PrivateKey(decodedPlainTextValue)
	if err != nil {
		h.Config.Logger.Error("failed to get plaintextValue", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// GET the Public Key Value
	certPEM, _ := base64.StdEncoding.DecodeString(secret.Metadata["public_key"].(string))
	entityCert, err := x509.ParseCertificate(certPEM)
	if err != nil {
		h.Config.Logger.Error("", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode to p12 File
	p12Data, err := pkcs12.Modern.Encode(privateKey, entityCert, []*x509.Certificate{subCert, rootCert}, "password")
	if err != nil {
		log.Println("Failed to create PKCS#12 file", err.Error())
		http.Error(w, "Failed to create PKCS#12 file", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("certificate.p12", p12Data, 0644)
	if err != nil {
		panic(err)
	}

	// Set the response headers and body to return the .p12 file
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.p12")
	w.Header().Set("Content-Type", "application/x-pkcs12")
	w.Write(p12Data)

	// var pemFile []byte
	// pemFile = append(pemFile, certPEM...)
	// pemFile = append(pemFile, keyPEM...)

	// // Set the response headers and body
	// w.Header().Set("Content-Disposition", "attachment; filename=certificate.pem")
	// w.Header().Set("Content-Type", "application/x-pem-file")
	// w.Write(pemFile)

	// publicKeyBase64, ok := secret.Metadata["public_key"].(string)
	// if !ok {
	// 	http.Error(w, "Public key not found in metadata", http.StatusInternalServerError)
	// 	return
	// }

	// publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	// if err != nil {
	// 	http.Error(w, "Failed to decode public key", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Disposition", "attachment; filename=certificate.pem")
	// w.Header().Set("Content-Type", "application/x-pem-file")
	// w.Write(pemFile)
}
