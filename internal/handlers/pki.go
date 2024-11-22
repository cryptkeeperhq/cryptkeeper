package handlers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/go-pg/pg/v10"
)

type CertificateAuthorityInput struct {
	Name        string `json:"name"`
	CACert      string `json:"ca_cert"`
	PrivateKey  string `json:"private_key"`
	Description string `json:"description"`
}

type GeneratedCertificate struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

func (h *Handler) getCAs(w http.ResponseWriter, r *http.Request) {
	var cas []models.CertificateAuthority
	err := h.DB.Model(&cas).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cas)
}

// Add a new CA
func (h *Handler) addCA(w http.ResponseWriter, r *http.Request) {
	var input CertificateAuthorityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	caCert, err := base64.StdEncoding.DecodeString(input.CACert)
	if err != nil {
		http.Error(w, "Invalid CA Certificate", http.StatusBadRequest)
		return
	}

	privateKey, err := base64.StdEncoding.DecodeString(input.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid Private Key", http.StatusBadRequest)
		return
	}

	ca := models.CertificateAuthority{
		Name:        input.Name,
		CACert:      caCert,
		PrivateKey:  privateKey,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = h.DB.Model(&ca).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate Root CA?
	err = h.DB.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		// Decode PEM CA certificate
		block, _ := pem.Decode([]byte(ca.CACert))
		if block == nil || block.Type != "CERTIFICATE" {
			log.Println("Failed to decode PEM block containing the certificate")
			return err
		}

		// var rootCA models.RootCA
		// rootCA.RootCert = block.Bytes
		// caCert, err := x509.ParseCertificate(block.Bytes)
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Println("Invalid CA Certificate")
			return err
		}

		// Decode PEM CA private key
		block, _ = pem.Decode([]byte(ca.PrivateKey))
		if block == nil {
			log.Println("Failed to decode PEM block containing the private key")
			// http.Error(w, "Failed to decode PEM block containing the private key", http.StatusInternalServerError)
			return err
		}
		// rootCA.RootKey = block.Bytes

		var caKey interface{}
		caKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			caKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				log.Println("Invalid CA Private Key", err.Error())
				return err
			}
		}

		// rootCert, rootKey, err := enginepki.GenerateRootCA(caCert, caKey.(*rsa.PrivateKey), input.Name)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(rootCert, rootKey)
		// rootCA := models.RootCA{
		// 	CA:       ca.ID,
		// 	RootCert: rootCert.Raw,
		// 	RootKey:  x509.MarshalPKCS1PrivateKey(rootKey),
		// }

		rootCA := models.RootCA{
			CA:       ca.ID,
			RootCert: ca.CACert,
			RootKey:  x509.MarshalPKCS1PrivateKey(caKey.(*rsa.PrivateKey)),
		}

		_, err = h.DB.Model(&rootCA).Insert()
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ca)
}

func (h *Handler) getTemplates(w http.ResponseWriter, r *http.Request) {

	var templates []models.CertificateTemplate
	err := h.DB.Model(&templates).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(templates)
}

// Add a new certificate template
func (h *Handler) addTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.CertificateTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	_, err := h.DB.Model(&template).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// func (h *Handler) issueCertificate(w http.ResponseWriter, r *http.Request) {
// 	var certReq models.CertificateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&certReq); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var ca models.CertificateAuthority
// 	if err := h.DB.Model(&ca).Where("id = ?", certReq.CACertID).Select(); err != nil {
// 		http.Error(w, "Invalid CA ID", http.StatusBadRequest)
// 		return
// 	}

// 	var template models.CertificateTemplate
// 	if err := h.DB.Model(&template).Where("id = ?", certReq.CATemplateID).Select(); err != nil {
// 		http.Error(w, "Invalid Template ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Decode PEM CA certificate
// 	block, _ := pem.Decode([]byte(ca.CACert))
// 	if block == nil || block.Type != "CERTIFICATE" {
// 		log.Println("Failed to decode PEM block containing the certificate")
// 		http.Error(w, "Failed to decode PEM block containing the certificate", http.StatusInternalServerError)
// 		return
// 	}

// 	caCert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		log.Println("Invalid CA Certificate")
// 		http.Error(w, "Invalid CA Certificate", http.StatusInternalServerError)
// 		return
// 	}

// 	// Decode PEM CA private key
// 	block, _ = pem.Decode([]byte(ca.PrivateKey))
// 	if block == nil {
// 		log.Println("Failed to decode PEM block containing the private key")
// 		http.Error(w, "Failed to decode PEM block containing the private key", http.StatusInternalServerError)
// 		return
// 	}

// 	var caKey interface{}
// 	caKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		caKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
// 		if err != nil {
// 			log.Println("Invalid CA Private Key", err.Error())
// 			http.Error(w, "Invalid CA Private Key", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	expiresAt := time.Now().AddDate(0, 0, template.ValidityPeriod)
// 	// expiresAt := time.Now().AddDate(0, 0, certReq.ValidityPeriod)

// 	// Generate end-entity certificate using the template and CA
// 	entityCert, entityKey, err := enginepki.GenerateEndEntityCertificateNew(caCert, caKey.(*rsa.PrivateKey), certReq.CommonName, template.Organization, expiresAt)
// 	if err != nil {
// 		log.Println("Failed to generate end-entity certificate", err.Error())
// 		http.Error(w, "Failed to generate end-entity certificate", http.StatusInternalServerError)
// 		return
// 	}

// 	// // Combine certificate and private key into a single PEM file
// 	// var pemFile []byte

// 	// cACert := utils.ConvertCertificateToPEM(entityCert.Raw)
// 	// privateKey, _ := utils.ConvertPrivateKeyToPEM(entityKey)

// 	// pemFile = append(pemFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: entityCert.Raw})...)
// 	// pemFile = append(pemFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(entityKey)})...)

// 	// pemFile = append(pemFile, []byte(cACert)...)
// 	// pemFile = append(pemFile, []byte(privateKey)...)

// 	// Create a PKCS#12 (PFX) file with the certificate, private key, and CA chain

// 	p12Data, err := pkcs12.Modern.Encode(entityKey, entityCert, []*x509.Certificate{caCert}, "password")
// 	if err != nil {
// 		log.Println("Failed to create PKCS#12 file", err.Error())
// 		http.Error(w, "Failed to create PKCS#12 file", http.StatusInternalServerError)
// 		return
// 	}

// 	err = ioutil.WriteFile("certificate.p12", p12Data, 0644)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Set the response headers and body to return the .p12 file
// 	w.Header().Set("Content-Disposition", "attachment; filename=certificate.p12")
// 	w.Header().Set("Content-Type", "application/x-pkcs12")
// 	w.Write(p12Data)

// 	// Set the response headers and body
// 	// w.Header().Set("Content-Disposition", "attachment; filename=certificate.pem")
// 	// w.Header().Set("Content-Type", "application/x-pem-file")
// 	// w.Write(pemFile)

// 	// w.Header().Set("Content-Disposition", "attachment; filename=certificate.pem")
// 	// w.Header().Set("Content-Type", "application/x-pem-file")
// 	// w.Write(publicKey)

// 	// w.WriteHeader(http.StatusCreated)
// 	// json.NewEncoder(w).Encode(map[string]interface{}{
// 	// 	"certificate": entityCert,
// 	// 	"private_key": entityKey,
// 	// })
// }

// func (h *Handler) downloadCertificate(w http.ResponseWriter, r *http.Request) {
// 	cert := r.URL.Query().Get("cert")
// 	w.Header().Set("Content-Disposition", "attachment; filename=certificate.crt")
// 	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
// 	w.Write([]byte(cert))
// }

// func (h *Handler) downloadPrivateKey(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Query().Get("key")
// 	w.Header().Set("Content-Disposition", "attachment; filename=private.key")
// 	w.Header().Set("Content-Type", "application/x-pem-file")
// 	w.Write([]byte(key))
// }
