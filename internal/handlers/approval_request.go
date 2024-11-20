package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/google/uuid"
)

func (h *Handler) CreateApprovalRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SecretID *string       `json:"secret_id,omitempty"`
		Action   string        `json:"action"`
		Details  models.Secret `json:"details"`
	}

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

	approvalRequest := models.ApprovalRequest{
		UserID:    userID,
		SecretID:  request.SecretID,
		Action:    request.Action,
		Details:   request.Details,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := h.DB.Model(&approvalRequest).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approvalRequest)
}

func (h *Handler) ListApprovalRequests(w http.ResponseWriter, r *http.Request) {
	// Get the status query parameter
	status := r.URL.Query().Get("status")

	var requests []models.ApprovalRequest
	query := h.DB.Model(&requests)

	// If status is specified, filter by status
	if status != "" {
		query.Where("status = ?", status)
	}

	err := query.Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func (h *Handler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RequestID int64 `json:"request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var approvalRequest models.ApprovalRequest
	err := h.DB.Model(&approvalRequest).Where("id = ?", request.RequestID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	approvalRequest.Status = "rejected"
	approvalRequest.UpdatedAt = time.Now()

	_, err = h.DB.Model(&approvalRequest).WherePK().Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approvalRequest)
}

func (h *Handler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RequestID int64 `json:"request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var approvalRequest models.ApprovalRequest
	err := h.DB.Model(&approvalRequest).Where("id = ?", request.RequestID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	approvalRequest.Status = "approved"
	approvalRequest.UpdatedAt = time.Now()

	_, err = h.DB.Model(&approvalRequest).WherePK().Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Perform the requested action based on approvalRequest.Action
	switch approvalRequest.Action {
	case "create":
		err = h.createSecretFromApprovalRequest(approvalRequest)
	case "update":
		err = h.updateSecretFromApprovalRequest(approvalRequest)
	case "delete":
		err = h.deleteSecretFromApprovalRequest(approvalRequest)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approvalRequest)
}

// Helper functions to handle secret creation, update, and deletion

func (h *Handler) createSecretFromApprovalRequest(approvalRequest models.ApprovalRequest) error {
	secret := approvalRequest.Details
	secret.ID = uuid.New().String()

	// Extract necessary fields from the details
	// var secret models.Secret
	// secret.PathID = int(details["path_id"].(float64))
	// // details["path_id"].(int)
	// secret.Value = details["value"].(string)
	// secret.Metadata = details["metadata"].(map[string]interface{})
	// expiresAtStr, _ := details["expires_at"].(string)
	// secret.IsOneTime, _ = details["is_one_time"].(bool)
	// secret.RotationInterval = details["rotation_interval"].(string)
	// secret.Key = details["key"].(string)

	// if expiresAtStr != "" {
	// 	t, err := time.Parse(time.RFC3339, expiresAtStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	secret.ExpiresAt = &t
	// }

	_, err := db.WriteSecret(approvalRequest.UserID, secret, h.CryptoOps)
	return err
}

func (h *Handler) updateSecretFromApprovalRequest(approvalRequest models.ApprovalRequest) error {
	secret := approvalRequest.Details
	// Find the secret by path
	// var secret models.Secret
	err := h.DB.Model(&secret).
		Where("path = ? and key = ?", secret.Path, secret.Key).
		Order("version DESC").
		Limit(1).
		Select()
	if err != nil {
		return err
	}

	var path models.Path
	err = h.DB.Model(&path).Where("id = ?", secret.PathID).Select()
	if err != nil {
		return err
	}

	// Decrypt the path key
	decryptedPathKeyHandle, err := h.CryptoOps.DecryptPathKey(path.KeyData)
	if err != nil {
		log.Printf("failed to decrypt path key: %v\n", err)
		return err
	}

	// encryptedDEKBase64, encryptedValue, err := utils.EncryptSecretValue(secret.Value, encryptionKey)
	encryptedDEKBase64, encryptedValue, err := h.CryptoOps.EncryptSecretValue(secret.Value, decryptedPathKeyHandle)
	if err != nil {
		return err
	}

	// Update the secret details
	secret.EncryptedDEK = encryptedDEKBase64
	secret.EncryptedValue = encryptedValue
	// secret.Metadata = metadata
	// secret.IsOneTime = isOneTime
	// secret.ExpiresAt = expiresAt
	secret.UpdatedAt = time.Now()
	// secret.RotationInterval = rotationInterval

	_, err = h.DB.Model(&secret).WherePK().Update()
	if err != nil {
		return err
	}

	// h.logAction(utils.Identity.GetUsername(), "update", &secret.ID, map[string]interface{}{
	// 	"path":    secret.PathID,
	// 	"key":     secret.Key,
	// 	"version": secret.Version,
	// })

	return nil
}

func (h *Handler) deleteSecretFromApprovalRequest(approvalRequest models.ApprovalRequest) error {
	secret := approvalRequest.Details

	// // Extract the path from the details
	// path := secret.Path
	// key := secret.Key

	// // Find the secret by path
	// var secret models.Secret
	err := h.DB.Model(&secret).
		Where("path = ? and key = ?", secret.PathID, secret.Key).
		Order("version DESC").
		Limit(1).
		Select()
	if err != nil {
		return err
	}

	// Track the deleted secret
	deletion := models.SecretDeletion{
		SecretID:       secret.ID,
		PathID:         secret.PathID,
		Key:            secret.Key,
		Version:        secret.Version,
		EncryptedDEK:   secret.EncryptedDEK,
		EncryptedValue: secret.EncryptedValue,
		Metadata:       secret.Metadata,
		DeletedAt:      time.Now(),
	}
	_, err = h.DB.Model(&deletion).Insert()
	if err != nil {
		return err
	}

	// Delete the secret
	_, err = h.DB.Model(&secret).WherePK().Delete()

	// h.logAction(approvalRequest.UserID, "delete", &secret.ID, map[string]interface{}{
	// 	"path":    secret.PathID,
	// 	"key":     secret.Key,
	// 	"version": secret.Version,
	// })

	return err
}
