package handlers

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	enginepki "github.com/cryptkeeperhq/cryptkeeper/internal/engine/pki"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"software.sslmate.com/src/go-pkcs12"
)

func (h *Handler) ListClientCerts(w http.ResponseWriter, r *http.Request) {
	var certificates []models.Certificate
	err := h.DB.Model(&certificates).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, certificates)
}

func (h *Handler) DownloadClientCA(w http.ResponseWriter, r *http.Request) {
	// Load CA private key and certificate
	caCertPEM, err := os.ReadFile("./scripts/certs/ca.pem")
	if err != nil {
		http.Error(w, "Failed to read CA certificate", http.StatusInternalServerError)
		return
	}

	block, _ := pem.Decode([]byte(caCertPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		http.Error(w, "Failed to decode PEM block containing the certificate", http.StatusBadRequest)
		return
	}

	subCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		http.Error(w, "Failed to ParseCertificate CA certificate", http.StatusInternalServerError)
		return
	}

	var bundle bytes.Buffer

	// Encode Sub CA Certificate in PEM format
	err = pem.Encode(&bundle, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: subCert.Raw,
	})
	if err != nil {
		http.Error(w, "Failed to Encode Bundle", http.StatusInternalServerError)
		return
	}

	// Set the response headers and body
	w.Header().Set("Content-Disposition", "attachment; filename=ca.pem")
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Write(bundle.Bytes())

}

func (h *Handler) CreateClientCert(w http.ResponseWriter, r *http.Request) {
	var req models.Certificate

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Load and parse CA certificate and key
	subCert, rootKey, err := loadCACertificateAndKey("./scripts/certs/ca.pem", "./scripts/certs/ca.key")
	if err != nil {
		http.Error(w, "Failed to load CA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the certificate already exists in the database
	err = h.DB.Model(&req).Where("name = ?", req.Name).Select()
	if err == nil {
		// If exists, retrieve and return the .p12 file
		if err := returnExistingPKCS12(w, req, subCert); err != nil {
			http.Error(w, "Failed to process existing certificate: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Generate new certificate
	entityCert, entityKey, err := generateNewCertificate(req, subCert, rootKey.(*rsa.PrivateKey))
	if err != nil {
		http.Error(w, "Failed to generate certificate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store certificate and key in database
	if err := storeCertificate(h.DB, &req, entityCert, entityKey); err != nil {
		http.Error(w, "Failed to store certificate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created .p12 file
	if err := returnPKCS12(w, entityCert, entityKey, subCert); err != nil {
		http.Error(w, "Failed to create PKCS#12 file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadCACertificateAndKey(caCertPath, caKeyPath string) (*x509.Certificate, interface{}, error) {
	// Load and parse CA certificate
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading CA certificate: %w", err)
	}
	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil || caCertBlock.Type != "CERTIFICATE" {
		return nil, nil, errors.New("invalid CA certificate PEM")
	}
	subCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing CA certificate: %w", err)
	}

	// Load and parse CA private key
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading CA private key: %w", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		return nil, nil, errors.New("invalid CA private key PEM")
	}

	rootKey, err := parsePrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing CA private key: %w", err)
	}

	return subCert, rootKey, nil
}

func parsePrivateKey(keyBytes []byte) (interface{}, error) {
	if key, err := x509.ParsePKCS1PrivateKey(keyBytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		return key, nil
	}
	return nil, errors.New("unsupported private key format")
}

func generateNewCertificate(req models.Certificate, subCert *x509.Certificate, rootKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	dnsNames := strings.Split(req.DNSName, ",")
	ipAddresses := strings.Split(req.IPAddresses, ",")
	emailAddresses := strings.Split(req.EmailAddresses, ",")

	expiresAt := time.Now().AddDate(0, 0, 365)
	return enginepki.GenerateEndEntityCertificateNew(subCert, rootKey, req.Name, "", dnsNames, ipAddresses, emailAddresses, expiresAt)
}

func storeCertificate(db *pg.DB, req *models.Certificate, cert *x509.Certificate, key *rsa.PrivateKey) error {
	privateKeyPEM := x509.MarshalPKCS1PrivateKey(key)
	req.Key = base64.StdEncoding.EncodeToString(privateKeyPEM)
	req.Cert = base64.StdEncoding.EncodeToString(cert.Raw)

	_, err := db.Model(req).Insert()
	return err
}

func returnExistingPKCS12(w http.ResponseWriter, req models.Certificate, subCert *x509.Certificate) error {
	keyBytes, _ := base64.StdEncoding.DecodeString(req.Key)
	entityKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("parsing existing private key: %w", err)
	}

	certBytes, _ := base64.StdEncoding.DecodeString(req.Cert)
	entityCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return fmt.Errorf("parsing existing certificate: %w", err)
	}

	return returnPKCS12(w, entityCert, entityKey, subCert)
}

func returnPKCS12(w http.ResponseWriter, cert *x509.Certificate, key *rsa.PrivateKey, subCert *x509.Certificate) error {
	p12Data, err := pkcs12.Modern.Encode(key, cert, []*x509.Certificate{subCert}, "password")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Disposition", "attachment; filename=certificate.p12")
	w.Header().Set("Content-Type", "application/x-pkcs12")
	_, err = w.Write(p12Data)
	return err
}

func (h *Handler) UserInGroup(userID int64, groupName string) bool {
	var group []models.Group
	err := h.DB.Model(&group).
		Column("g.*").
		TableExpr("groups AS g").
		Join("JOIN user_groups ugr ON ugr.group_id = g.id").
		Where("ugr.user_id = ? AND g.name = ?", userID, groupName).
		Select()

	fmt.Println("Group", group)
	return err == nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.KAuth.CreateUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// user.Password = string(hashedPassword)
	// user.CreatedAt = time.Now()

	// _, err = h.DB.Model(&user).Insert()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Config.Auth.SSOEnabled {
		err := h.KAuth.CreateGroup(group.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		group.ID = uuid.New().String()
		group.CreatedAt = time.Now()
		_, err := h.DB.Model(&group).Insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// type KUserGroup struct {
// 	UserID  string `json:"user_id"`
// 	GroupID string `json:"group_id"`
// }

func (h *Handler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	var userGroup models.UserGroup
	if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Config.Auth.SSOEnabled {
		err := h.KAuth.AddUserToGroup(userGroup.UserID, userGroup.GroupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := h.DB.Model(&userGroup).Insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// var g models.Group
		// h.DB.Model(&g).Where("id = ?", userGroup.GroupID).First()

		// var u models.User
		// h.DB.Model(&u).Where("id = ?", userGroup.UserID).First()

		// h.Z.AddUserToGroup(u.Username, g.Name)

	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	// var userGroup models.UserGroup
	// if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// _, err := h.DB.Model(&userGroup).Where("user_id = ? AND group_id = ?", userGroup.UserID, userGroup.GroupID).Delete()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	var userGroup models.UserGroup
	if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// vars := mux.Vars(r)
	// userID := vars["userID"]

	err := h.KAuth.RemoveUserFromGroup(userGroup.UserID, userGroup.GroupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {

	if h.Config.Auth.SSOEnabled {
		users, err := h.KAuth.GetUsers("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		return
	}

	var users []models.User
	err := h.DB.Model(&users).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, users)
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {

	if h.Config.Auth.SSOEnabled {
		groups, err := h.KAuth.GetGroups("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}

	var groups []models.Group
	err := h.DB.Model(&groups).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *Handler) ListGroupUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	var users []models.User

	if h.Config.Auth.SSOEnabled {
		kUsers, err := h.KAuth.GetUsers(groupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, user := range kUsers {
			users = append(users, models.User{
				ID:       *user.ID,
				Username: *user.Username,
			})
		}
	} else {
		err := h.DB.Model(&users).Distinct().
			Column("u.id", "u.username", "u.created_at").
			TableExpr("users AS u").
			Join("JOIN user_groups AS ug ON ug.user_id = u.id").
			Where("ug.group_id = ?", groupID).
			Select()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var groups []models.Group

	if h.Config.Auth.SSOEnabled {
		kGroups, err := h.KAuth.GetGroups(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, group := range kGroups {
			groups = append(groups, models.Group{
				ID:   *group.ID,
				Name: *group.Name,
			})
		}
	} else {
		err := h.DB.Model(&groups).Distinct().
			Column("g.id", "g.name", "g.created_at").
			TableExpr("groups AS g").
			Join("JOIN user_groups AS ug ON ug.group_id = g.id").
			Where("ug.user_id = ?", userID).
			Select()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// func (h *Handler) AssignGroupAccess(w http.ResponseWriter, r *http.Request) {
// 	var access models.SecretAccess
// 	if err := json.NewDecoder(r.Body).Decode(&access); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if access.GroupID == nil {
// 		http.Error(w, "Group ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	_, err := h.DB.Model(&access).Insert()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *Handler) getUserGroups(userID int64) ([]int64, error) {
// 	var groupIDs []int64
// 	err := h.DB.Model((*models.UserGroup)(nil)).
// 		Column("group_id").
// 		Where("user_id = ?", userID).
// 		Select(&groupIDs)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return groupIDs, nil
// }
