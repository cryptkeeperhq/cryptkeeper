package handlers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	enginepki "github.com/cryptkeeperhq/cryptkeeper/internal/engine/pki"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/policy"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var permissions = []string{"list", "read", "create", "update", "delete", "rotate"}

func (h *Handler) GetPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathID := vars["pathID"]

	var path models.Path

	err := h.DB.Model(&path).Where("id = ?", pathID).First()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoPath(path))
}

func (h *Handler) GetPathPermissions(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	pathID := vars["pathID"]

	var path models.Path

	err := h.DB.Model(&path).Where("id = ?", pathID).First()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var policy models.Policy
	err = h.DB.Model(&policy).Where("path_id = ?", path.ID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pathPermissions := make(map[string]bool)

	userGroups, _ := db.GetUserGroups(identity.GetID())

	denyPermissions := []string{}

	for _, permission := range permissions {

		pathPermissions[permission] = false

		for _, policyPath := range policy.Paths {
			// If deny on any policy than don't iterate further
			if utils.Contains(denyPermissions, permission) {
				continue
			}

			// Check if user or any group has permission for this path
			if utils.Contains(policyPath.Users, identity.GetUsername()) || h.hasGroupAccess(policyPath.Groups, userGroups) {
				if utils.Contains(policyPath.DenyPermissions, permission) {
					denyPermissions = append(denyPermissions, permission)
					pathPermissions[permission] = false
					break
				}

				if utils.Contains(policyPath.Permissions, permission) {
					pathPermissions[permission] = true
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"permissions": pathPermissions,
	})
}

func (h *Handler) CreatePath(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var path models.Path
	if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path.ID = uuid.New().String()
	keyHandle, err := h.CryptoOps.GeneratePathKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate and handle engine type specific logic
	switch path.EngineType {
	case "kv":
		// Handle KV engine specific initialization
	case "pki":
		// Handle PKI engine specific initialization
	case "database":
		// Handle Database secrets engine specific initialization
		// path.Metadata = map[string]interface{}{
		// 	"connection_string": "your-database-connection-string",
		// 	"role_template":     "CREATE ROLE '{{name}}' WITH LOGIN PASSWORD '{{password}}'",
		// }
	case "transit":
		// Handle Transit engine specific initialization
	default:
		http.Error(w, "Unsupported engine type", http.StatusBadRequest)
		return
	}

	encryptedPathKey, err := h.CryptoOps.EncryptPathKey(keyHandle)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path.KeyData = encryptedPathKey
	err = h.DB.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		path.CreatedAt = time.Now()
		path.UpdatedAt = time.Now()

		_, err = h.DB.Model(&path).Insert()
		if err != nil {
			return err
		}

		if path.EngineType == "pki" {

			// if path.Metadata["root_ca"].(string) == "cryptkeeper_ca" {
			// 	err := h.DB.Model(&rootCA).First()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	rootCert, err := x509.ParseCertificate(rootCA.RootCert)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	rootKey, err := x509.ParsePKCS1PrivateKey(rootCA.RootKey)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	subCert, subKey, err := enginepki.GenerateSubCA(rootCert, rootKey, path.Path)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	subCA := models.SubCA{
			// 		PathID:    path.ID,
			// 		SubCACert: subCert.Raw,
			// 		SubCAKey:  x509.MarshalPKCS1PrivateKey(subKey),
			// 	}

			// 	_, err = h.DB.Model(&subCA).Insert()
			// 	if err != nil {
			// 		return err
			// 	}
			// } else {
			fmt.Println("Genertaing CA for ", path.Metadata["root_ca"].(string))

			// var ca models.CertificateAuthority
			// err := h.DB.Model(&ca).Where("id = ?", utils.ToInt(path.Metadata["root_ca"].(string))).First()
			// if err != nil {
			// 	return err
			// }

			var rootCA models.RootCA
			err := h.DB.Model(&rootCA).Where("ca = ?", utils.ToInt(path.Metadata["root_ca"].(string))).First()
			if err != nil {
				return err
			}

			// // Decode PEM CA certificate
			// block, _ := pem.Decode([]byte(rootCA.RootCert))
			// if block == nil || block.Type != "CERTIFICATE" {
			// 	log.Println("Failed to decode PEM block containing the certificate")
			// 	return err
			// }

			block, _ := pem.Decode([]byte(rootCA.RootCert))
			if block == nil || block.Type != "CERTIFICATE" {
				log.Println("Failed to decode PEM block containing the certificate")
				return err
			}

			// var rootCA models.RootCA
			// rootCA.RootCert = block.Bytes
			rootCert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Println("Invalid CA Certificate")
				return err
			}

			// rootCert, err := x509.ParseCertificate(rootCA.RootCert)
			// if err != nil {
			// 	log.Println("Invalid CA Certificate")
			// 	return err
			// }

			// var rootCA models.RootCA
			// rootCA.RootCert = block.Bytes

			// Decode PEM CA private key
			// block, _ := pem.Decode([]byte(rootCA.RootKey))
			// if block == nil {
			// 	log.Println("Failed to decode PEM block containing the private key")
			// 	// http.Error(w, "Failed to decode PEM block containing the private key", http.StatusInternalServerError)
			// 	return err
			// }
			// rootCA.RootKey = block.Bytes

			var rootKey interface{}
			rootKey, err = x509.ParsePKCS1PrivateKey(rootCA.RootKey)
			if err != nil {
				rootKey, err = x509.ParsePKCS8PrivateKey(rootCA.RootKey)
				if err != nil {
					log.Println("Invalid CA Private Key", err.Error())
					return err
				}
			}

			subCert, subKey, err := enginepki.GenerateSubCA(rootCert, rootKey.(*rsa.PrivateKey), path.Path)
			if err != nil {
				return err
			}

			subCA := models.SubCA{
				RootCA:    rootCA.ID,
				PathID:    path.ID,
				SubCACert: subCert.Raw,
				SubCAKey:  x509.MarshalPKCS1PrivateKey(subKey),
			}

			_, err = h.DB.Model(&subCA).Insert()
			if err != nil {
				return err
			}
		}
		// }

		defaultPolicy, err := policy.GetDefaultHCLPolicy(path.Path)
		if err != nil {
			return err
		}

		// defaultPolicy.ID = uuid.New().String()
		defaultPolicy.PathID = path.ID
		defaultPolicy.Name = "Default Policy"
		defaultPolicy.Description = "Created during path creation"
		defaultPolicy.CreatedAt = time.Now()
		defaultPolicy.UpdatedAt = time.Now()

		err = db.SavePolicy(defaultPolicy, identity.GetUsername())

		if err != nil {
			h.Config.Logger.Error(fmt.Sprintf("Failed to create default policy: %s", err))
		}

		defaultPolicyJson, _ := json.Marshal(defaultPolicy)
		if err := h.Producer.SendMessage("policy_updates", defaultPolicyJson); err != nil {
			log.Printf("Error sending message: %v", err)
		}

		// topic := "policy_updates"
		// err = h.Producer.SendMessage(topic, defaultPolicy)
		// if err != nil {
		// 	h.Config.Logger.Error(fmt.Sprintf("Failed to produce policy update: %s", err))
		// }
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoPath(path))
}

func (h *Handler) UpdatePath(w http.ResponseWriter, r *http.Request) {
	var path models.Path
	if err := json.NewDecoder(r.Body).Decode(&path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingPath, err := db.GetPathByID(path.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingPath.Path = path.Path
	existingPath.Metadata = path.Metadata
	existingPath.UpdatedAt = time.Now()

	_, err = h.DB.Model(&existingPath).Where("id = ?", existingPath.ID).Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoPath(path))
}

// ListAllPaths returns all paths regardless of permissions
func (h *Handler) ListAllPaths(w http.ResponseWriter, r *http.Request) {
	var paths []models.Path
	var err error

	pathQuery := r.URL.Query().Get("path")
	if pathQuery != "" {
		var path models.Path
		path, err = db.GetPathByName(pathQuery)
		paths = append(paths, path)
	} else {
		err = h.DB.Model(&paths).Select()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}

type PathResponse struct {
	ID         string                 `json:"id"`
	Path       string                 `json:"path"`
	EngineType string                 `json:"engine_type"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedBy  int64                  `json:"created_by"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

func dtoPath(path models.Path) PathResponse {
	return PathResponse{
		ID:         path.ID,
		Path:       path.Path,
		EngineType: path.EngineType,
		Metadata:   path.Metadata,
		CreatedBy:  path.CreatedBy,
		CreatedAt:  path.CreatedAt,
		UpdatedAt:  path.UpdatedAt,
	}
}

func (h *Handler) ListUserPaths(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	engine := r.URL.Query().Get("engine")

	userID := identity.GetID()
	var userPolicyIds []string
	var groupPolicyIds []string

	if identity.GetAuthType() == "user" {
		// Get user policies
		h.DB.Model((*models.UserPolicy)(nil)).ColumnExpr("array_agg(distinct policy_id)").Where("user_id = ?", userID).Select(pg.Array(&userPolicyIds))

		// Get groups the user belongs to
		var userGroups []string
		err := h.DB.Model((*models.UserGroup)(nil)).
			Column("group_id").
			Where("user_id = ?", userID).
			Select(&userGroups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("User Policy IDs", groupPolicyIds)

		// Get group policies
		if len(userGroups) > 0 {
			err = h.DB.Model((*models.GroupPolicy)(nil)).ColumnExpr("array_agg(distinct policy_id)").Where("group_id IN (?)", pg.In(userGroups)).Select(pg.Array(&groupPolicyIds))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if identity.GetAuthType() == "approle" {
		// Get app role policies
		h.DB.Model((*models.AppPolicy)(nil)).ColumnExpr("array_agg(distinct policy_id)").Where("app_id = ?", userID).Select(pg.Array(&userPolicyIds))
	}

	fmt.Println("USER Policy IDs", userPolicyIds)

	// if identity.GetAuthType() == "approle" {
	// 	// Get app role policies
	// 	h.DB.Model((*models.AppRolesPolicy)(nil)).ColumnExpr("array_agg(distinct policy_id)").Where("app_role_id = ?", userID).Select(pg.Array(&userPolicyIds))
	// 	fmt.Println("IDs", userPolicyIds)
	// }

	policyIdSet := append(userPolicyIds, groupPolicyIds...)
	fmt.Println("Policy IDS", policyIdSet)
	paths := h.pathsByPolicies(policyIdSet, engine)

	// paths, _ := db.GetPaths()

	var responses []PathResponse

	for _, path := range paths {
		responses = append(responses, dtoPath(path))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *Handler) pathsByPolicies(policyIds []string, engine string) []models.Path {

	// Fetch paths
	var paths []models.Path

	if len(policyIds) > 0 {
		query := h.DB.Model(&paths).
			Column("path.*").
			Join("join policies on policies.path_id = path.id").
			Where("policies.id IN (?)", pg.In(policyIds))

		if engine != "" {
			query.Where("engine_type = ?", engine)
		}
		query.Select()
	}

	return paths
}

// func (h *Handler) ListPathSecrets(w http.ResponseWriter, r *http.Request) {
// 	pathID, err := strconv.ParseInt(mux.Vars(r)["pathID"], 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid path ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch secrets for the given path
// 	secrets, err := db.GetSecretsByPathID(pathID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(secrets)

// }
