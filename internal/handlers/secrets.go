package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// TODO:
// - Implement write ahead log
// - Add checksum value to secrets DB

type CreateSecretRequest struct {
	PathID           string                 `json:"path_id"`
	Key              string                 `json:"key"`
	Value            string                 `json:"value"`
	MultiValue       map[string]interface{} `json:"multi_value"`
	ExpiresAt        *time.Time             `json:"expires_at,omitempty"`
	Metadata         map[string]interface{} `json:"metadata"`
	IsOneTime        bool                   `json:"is_one_time"`
	RotationInterval string                 `json:"rotation_interval"`
	IsMultiValue     bool                   `json:"is_multi_value"`
	Path             string                 `json:"path"`
	Tags             []string               `json:"tags"`
}

// Data Transfer Object (DTO) pattern for secret

type SecretResponse struct {
	ID               string                 `json:"id"`
	PathID           string                 `json:"path_id"`
	Key              string                 `json:"key"`
	Version          int                    `json:"version"`
	Checksum         string                 `json:"checksum"`
	Metadata         map[string]interface{} `json:"metadata"`
	IsMultiValue     bool                   `json:"is_multi_value"`
	Tags             []string               `json:"tags"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	IsOneTime        bool                   `json:"is_one_time"`
	ExpiresAt        *time.Time             `json:"expires_at,omitempty"`
	RotatedAt        *time.Time             `json:"rotated_at,omitempty"`
	RotationInterval string                 `json:"rotation_interval,omitempty"`
	LastRotatedAt    *time.Time             `json:"last_rotated_at,omitempty"`
	CreatedBy        string                 `json:"created_by"`
	Value            string                 `json:"value"`
	Path             string                 `json:"path"`
	KeyType          string                 `json:"key_type"`
	CreatedByUser    string                 `json:"created_by_user,omitempty"`
}

func dtoSecret(secret models.Secret) SecretResponse {
	response := SecretResponse{
		ID:               secret.ID,
		PathID:           secret.PathID,
		Key:              secret.Key,
		Version:          secret.Version,
		Checksum:         secret.Checksum,
		Metadata:         secret.Metadata,
		IsMultiValue:     secret.IsMultiValue,
		Tags:             secret.Tags,
		CreatedAt:        secret.CreatedAt,
		UpdatedAt:        secret.UpdatedAt,
		IsOneTime:        secret.IsOneTime,
		ExpiresAt:        secret.ExpiresAt,
		RotatedAt:        secret.RotatedAt,
		RotationInterval: secret.RotationInterval,
		LastRotatedAt:    secret.LastRotatedAt,
		CreatedBy:        secret.CreatedBy,
		Value:            secret.Value,
		Path:             secret.Path,
		KeyType:          secret.KeyType,
		CreatedByUser:    secret.CreatedByUser,
	}
	return response
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var req CreateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	userID := identity.GetID()

	// Check if it's a multi-value secret
	var value string
	if req.IsMultiValue {
		jsonValue, err := json.Marshal(req.MultiValue)
		if err != nil {
			http.Error(w, "Invalid multi-value data", http.StatusBadRequest)
			return
		}
		value = string(jsonValue)
	} else {
		value = req.Value
	}

	if req.Path != "" {
		path, err := db.GetPathByName(req.Path)
		if err != nil {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		req.PathID = path.ID
	}

	// Proceed with creating the secret in the database
	secret := models.Secret{
		ID:           uuid.New().String(),
		PathID:       req.PathID,
		Key:          req.Key,
		Value:        value,
		Checksum:     utils.MD5Checksum([]byte(value)),
		IsMultiValue: req.IsMultiValue,
		ExpiresAt:    req.ExpiresAt,
		Metadata:     req.Metadata,
		IsOneTime:    req.IsOneTime,
		CreatedBy:    identity.GetUsername(),
		Tags:         req.Tags,
	}

	// Parse the rotation interval if provided
	if req.RotationInterval != "" {
		duration, err := time.ParseDuration(req.RotationInterval)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s - %s", "Invalid rotation interval", err.Error()), http.StatusBadRequest)
			return
		}
		secret.RotationInterval = duration.String()
	}

	secret, err := db.WriteSecret(userID, secret, h.CryptoOps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log Action
	h.logAction(identity.GetUsername(), "create", &secret.ID, map[string]interface{}{
		"path":       secret.PathID,
		"key":        secret.Key,
		"version":    secret.Version,
		"created_by": identity.GetUsername(),
	})

	// // Publish secret update message for Zanzibar
	// topic := "secret_updates"
	// err = h.Producer.ProduceMessage(topic, secret)
	// if err != nil {
	// 	log.Printf("Failed to produce secret update: %s", err)
	// }

	// // Publish Event
	// user := events.User{
	// 	UserID:   identity.GetID(),
	// 	Username: identity.GetUsername(),
	// }
	// secretJson, _ := json.Marshal(secret)
	// details := map[string]interface{}{}
	// json.Unmarshal(secretJson, &details)
	// event := events.NewEvent(events.SecretCreated, user, details)
	// h.Producer.ProduceMessage("events", event)

	response := dtoSecret(secret)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetSecrets(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	secrets, err := db.GetSecretsByPathName(path)
	if err != nil && err != pg.ErrNoRows {
		h.Config.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []SecretResponse
	for _, secret := range secrets {
		response := dtoSecret(secret)
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")
	versionStr := r.URL.Query().Get("version")

	if path == "" || versionStr == "" {
		http.Error(w, "Path and version are required", http.StatusBadRequest)
		return
	}

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	secret, err := db.GetSecret(path, key, utils.ToInt64(versionStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.DeleteSecret(secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish Event
	// user := events.User{
	// 	UserID:   identity.GetID(),
	// 	Username: identity.GetUsername(),
	// }
	// secretJson, _ := json.Marshal(secret)
	// details := map[string]interface{}{}
	// json.Unmarshal(secretJson, &details)
	// event := events.NewEvent(events.SecretDeleted, user, details)
	// h.Producer.ProduceMessage("events", event)

	h.logAction(identity.GetUsername(), "delete", &secret.ID, map[string]interface{}{
		"path":    secret.PathID,
		"key":     secret.Key,
		"version": secret.Version,
	})

	response := dtoSecret(secret)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetDeletedSecrets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathID := vars["pathID"]

	var deletions []models.SecretDeletion
	err := h.DB.Model(&deletions).Where("path_id = ?", pathID).Order("deleted_at DESC").Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if deletions == nil {
		deletions = []models.SecretDeletion{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletions)
}

func (h *Handler) RestoreDeletedSecret(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pathID := vars["pathID"]
	secretID := vars["secretID"]

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var deletion models.SecretDeletion
	err = h.DB.Model(&deletion).Where("id = ? and path_id = ? and secret_id = ?", id, pathID, secretID).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			http.Error(w, "Deleted secret not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Restore the deleted secret
	secret := models.Secret{
		ID:             deletion.SecretID,
		PathID:         deletion.PathID,
		Version:        deletion.Version,
		Key:            deletion.Key,
		EncryptedDEK:   deletion.EncryptedDEK,
		EncryptedValue: deletion.EncryptedValue,
		Metadata:       deletion.Metadata,
		CreatedAt:      deletion.DeletedAt, // Use deleted_at as created_at
		UpdatedAt:      time.Now(),
	}

	_, err = h.DB.Model(&secret).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove the entry from secret_deletions
	_, err = h.DB.Model(&deletion).Where("id = ?", id).Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSecretHistory(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")
	if path == "" || key == "" {
		http.Error(w, "Path and key are required", http.StatusBadRequest)
		return
	}

	var secrets []models.Secret
	query := h.DB.Model(&secrets).Distinct().
		Column("secret.id", "p.path", "secret.path_id", "secret.key", "secret.version", "secret.metadata", "secret.created_by", "secret.created_at", "secret.updated_at", "secret.expires_at", "secret.rotated_at").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("secret.key = ?", key).
		Where("p.path = ?", path).
		Order("secret.version DESC").
		Limit(20) // Add limit to get only the last 20 versions

	// TODO: add pagination

	err := query.Select()
	if err != nil && err != pg.ErrNoRows {
		log.Printf("Error querying secret history: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []SecretResponse
	for _, secret := range secrets {
		responses = append(responses, dtoSecret(secret))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	version := r.URL.Query().Get("version")

	secret, err := db.GetSecretByID(id, utils.ToInt64(version))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoSecret(secret))
}

func (h *Handler) getSecretValue(secret *models.Secret) (string, error) {
	// Fetch the path to determine the engine type
	var pathDetails models.Path
	err := h.DB.Model(&pathDetails).Where("id = ?", secret.PathID).Select()
	if err != nil {
		return "", err
	}

	// Decrypt the path key
	decryptedPathKeyHandle, err := h.CryptoOps.DecryptPathKey(pathDetails.KeyData)
	fmt.Println(decryptedPathKeyHandle)
	if err != nil {
		return "", err
	}

	// Get plaintext value
	plaintextValue, err := h.CryptoOps.DecryptSecretValue(secret.EncryptedDEK, secret.EncryptedValue, decryptedPathKeyHandle)
	if err != nil {
		return "", err
	}

	// if pathDetails.EngineType == "pki" {
	// 	fmt.Println("plaintextValue", string(plaintextValue))
	// 	decodedValue, _ := base64.StdEncoding.DecodeString(plaintextValue)
	// 	return string(decodedValue), nil
	// }
	return plaintextValue, nil
}

func (h *Handler) RotateSecret(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")

	// var request models.Secret
	var request CreateSecretRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	userID := identity.GetID()

	existingSecret, err := db.GetSecret(path, key, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	existingSecret.CreatedBy = userID
	if request.ExpiresAt != nil {
		existingSecret.ExpiresAt = request.ExpiresAt
	}
	existingSecret.IsOneTime = request.IsOneTime

	// Keep the same value and just rotate the encryption keys
	if request.Value == "" && request.MultiValue == nil {
		// Get plaintext value
		existingSecret.Value, err = h.getSecretValue(&existingSecret)
		if err != nil {
			h.Config.Logger.Error("failed to get plaintextValue", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// existingSecret.Value = request.Value
		// var value string
		if request.IsMultiValue {
			jsonValue, err := json.Marshal(request.MultiValue)
			if err != nil {
				http.Error(w, "Invalid multi-value data", http.StatusBadRequest)
				return
			}
			existingSecret.Value = string(jsonValue)
		} else {
			existingSecret.Value = request.Value
		}
	}

	// if req.IsMultiValue {
	// 	jsonValue, err := json.Marshal(req.MultiValue)
	// 	if err != nil {
	// 		http.Error(w, "Invalid multi-value data", http.StatusBadRequest)
	// 		return
	// 	}
	// 	value = string(jsonValue)
	// } else {
	// 	value = req.Value
	// }

	// if existingSecret.IsMultiValue {
	// 	// Multi-value secret handling
	// 	existingSecret.Value = request.Value
	// }  else {
	// 	// Single-value secret handling
	// 	existingSecret.Value = request.Value
	// }

	existingSecret.Checksum = utils.MD5Checksum([]byte(existingSecret.Value))
	existingSecret.CreatedBy = identity.GetUsername()

	newSecret, err := db.WriteSecret(userID, existingSecret, h.CryptoOps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// topic := "secret_updates"
	// err = h.Producer.ProduceMessage(topic, newSecret)
	// if err != nil {
	// 	log.Printf("Failed to produce secret update: %s", err)
	// }

	h.logAction(identity.GetUsername(), "rotate", &existingSecret.ID, map[string]interface{}{
		"path":        newSecret.PathID,
		"key":         newSecret.Key,
		"version":     newSecret.Version,
		"old_version": existingSecret.Version,
	})

	// // Publish Event
	// user := events.User{
	// 	UserID:   identity.GetID(),
	// 	Username: identity.GetUsername(),
	// }
	// secretJson, _ := json.Marshal(newSecret)
	// details := map[string]interface{}{}
	// json.Unmarshal(secretJson, &details)
	// event := events.NewEvent(events.SecretRotated, user, details)
	// h.Producer.ProduceMessage("events", event)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoSecret(newSecret))
}

func (h *Handler) GetSecretVersion(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")
	version := r.URL.Query().Get("version")
	if path == "" && key == "" {
		http.Error(w, "Key and Version are required", http.StatusBadRequest)
		return
	}

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	userID := identity.GetID()

	secret, err := db.GetSecret(path, key, utils.ToInt64(version))
	if err != nil {
		log.Printf("Error querying secret: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the secret has expired
	if secret.ExpiresAt != nil && time.Now().After(*secret.ExpiresAt) {
		http.Error(w, "Secret has expired", http.StatusGone)
		return
	}

	plaintextValue, err := h.getSecretValue(&secret)
	if err != nil {
		h.Config.Logger.Error("failed to get plaintextValue", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(plaintextValue)

	// if utils.MD5Checksum([]byte(plaintextValue)) != secret.Checksum {
	// 	http.Error(w, "invalid checksum for plain value", http.StatusInternalServerError)
	// 	return
	// }

	// Set the plaintext value in the response
	// if pathDetails.EngineType == "pki" {
	// 	decodedValue, _ := base64.StdEncoding.DecodeString(plaintextValue)
	// 	secret.Value = string(decodedValue)
	// } else {
	secret.Value = string(plaintextValue)
	// }

	// If the secret is a one-time secret, delete it after retrieving the value
	if secret.IsOneTime {
		err = db.DeleteSecret(secret)
		if err != nil {
			log.Printf("Error deleting one-time secret: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = h.LogSecretAccess(secret.ID, userID, identity.GetAuthType())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logAction(identity.GetUsername(), "view", &secret.ID, map[string]interface{}{
		"source":  identity.GetAuthType(),
		"path":    secret.PathID,
		"key":     secret.Key,
		"version": secret.Version,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoSecret(secret))
}

func (h *Handler) SearchSecrets(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")
	metadata := r.URL.Query().Get("metadata")
	createdAfter := r.URL.Query().Get("created_after")
	createdBefore := r.URL.Query().Get("created_before")

	var secrets []models.Secret
	query := h.DB.Model(&secrets).Distinct().
		Column("secret.*", "p.path").
		Join("JOIN paths p ON p.id = secret.path_id")
		// Join("JOIN secret_accesses ON secret_accesses.secret_id = secret.id")
		// Where("secret_accesses.user_id = ? OR secret_accesses.group_id IN (SELECT group_id FROM user_groups WHERE user_id = ?)", userID, userID)

	if path != "" {
		query.Where("p.path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}

	if key != "" {
		query.Where("secret.key LIKE ?", fmt.Sprintf("%%%s%%", key))
	}

	if metadata != "" {
		query.Where("secret.metadata @> ?", metadata)
	}

	if createdAfter != "" {
		query.Where("secret.created_at >= ?", createdAfter)
	}

	if createdBefore != "" {
		query.Where("secret.created_at <= ?", createdBefore)
	}

	err := query.Order("secret.version DESC").Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

func (h *Handler) UpdateSecretMetadata(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// path := vars["pathID"]

	path := r.URL.Query().Get("path")
	key := r.URL.Query().Get("key")
	versionStr := r.URL.Query().Get("version")

	if path == "" || key == "" || versionStr == "" {
		http.Error(w, "Path and version are required", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// userID := identity.GetID()

	var metadata map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var secret models.Secret
	err = h.DB.Model(&secret).
		// Join("JOIN secret_accesses ON secret_accesses.secret_id = secret.id").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("p.path = ?", path).
		Where("secret.key = ? AND secret.version = ?", key, version).
		// Where("secret_accesses.user_id = ? OR secret_accesses.group_id IN (SELECT group_id FROM user_groups WHERE user_id = ?)", userID, userID).
		Limit(1). // Ensure only one row is returned
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			http.Error(w, "Secret not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	oldMetadata := secret.Metadata
	secret.Metadata = metadata
	secret.UpdatedAt = time.Now()

	_, err = h.DB.Model(&secret).WherePK().Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logAction(identity.GetUsername(), "updated", &secret.ID, map[string]interface{}{
		"path":             secret.PathID,
		"key":              secret.Key,
		"version":          secret.Version,
		"metadata":         secret.Metadata,
		"previous_metdata": oldMetadata,
	})

	// Publish Event
	// user := events.User{
	// 	UserID:   identity.GetID(),
	// 	Username: identity.GetUsername(),
	// }
	// secretJson, _ := json.Marshal(secret)
	// details := map[string]interface{}{}
	// json.Unmarshal(secretJson, &details)

	// event := events.NewEvent(events.SecretUpdated, user, details)
	// h.Producer.ProduceMessage("events", event)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoSecret(secret))
}
