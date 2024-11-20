package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
)

func (h *Handler) LogSecretAccess(secretID, userID string, source string) error {
	accessLog := models.AccessLog{
		SecretID:   secretID,
		UserID:     userID,
		AccessedAt: time.Now(),
		Source:     source,
	}

	_, err := h.DB.Model(&accessLog).Insert()
	return err
}

func (h *Handler) AccessSecret(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	userID := identity.GetID()

	secretID := r.URL.Query().Get("secret_id")
	if secretID == "" {
		http.Error(w, "secret_id is required", http.StatusBadRequest)
		return
	}

	var secret models.Secret
	err := h.DB.Model(&secret).Where("id = ?", secretID).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.LogSecretAccess(secret.ID, userID, "UI")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}
