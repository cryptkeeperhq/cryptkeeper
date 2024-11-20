package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create a shared link
func (h *Handler) CreateSharedLink(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var request struct {
		SecretID  string    `json:"secret_id"`
		ExpiresAt time.Time `json:"expires_at"`
		Version   int       `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	linkID := uuid.New().String()
	sharedLink := models.SharedLink{
		LinkID:    linkID,
		SecretID:  request.SecretID,
		ExpiresAt: request.ExpiresAt,
		Version:   request.Version,
	}

	_, err := h.DB.Model(&sharedLink).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logAction(identity.GetUsername(), "shared", &request.SecretID, map[string]interface{}{
		"source": identity.GetAuthType(),
		"link":   sharedLink,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sharedLink)
}

// Access a shared link
func (h *Handler) AccessSharedLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID := vars["linkID"]

	var sharedLink models.SharedLink
	err := h.DB.Model(&sharedLink).Where("link_id = ?", linkID).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			http.Error(w, "Link not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if time.Now().After(sharedLink.ExpiresAt) {
		http.Error(w, "Link has expired", http.StatusGone)
		return
	}

	// var secret models.Secret
	// err = h.DB.Model(&secret).Where("id = ?", sharedLink.SecretID).Select()
	secret, err := db.GetSecretByID(sharedLink.SecretID, int64(sharedLink.Version))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var path models.Path
	err = h.DB.Model(&path).Where("id = ?", secret.PathID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decrypt the path key
	decryptedPathKeyHandle, err := h.CryptoOps.DecryptPathKey(path.KeyData)
	if err != nil {
		log.Printf("failed to decrypt path key: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decrypt the secret value
	plaintextValue, err := h.CryptoOps.DecryptSecretValue(secret.EncryptedDEK, secret.EncryptedValue, decryptedPathKeyHandle)
	if err != nil {
		h.Config.Logger.Error("Error decrypting secret", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the plaintext value in the response
	secret.Value = string(plaintextValue)

	// TODO: Should I delete the one time use here? What is there is an access issue....
	// // If the secret is a one-time secret, delete it after retrieving the value
	// if secret.IsOneTime {
	// 	err = h.DB.RunInTransaction(r.Context(), func(tx *pg.Tx) error {
	// 		_, err := tx.Model(&secret).Where("id = ?", secret.ID).Delete()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		_, err = tx.Model((*models.SecretAccess)(nil)).Where("secret_id = ?", secret.ID).Delete()
	// 		return err
	// 	})
	// 	if err != nil {
	// 		log.Printf("Error deleting one-time secret: %v", err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	h.logAction("", "shared_link_access", &sharedLink.SecretID, map[string]interface{}{
		"link": sharedLink,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}
