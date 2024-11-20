package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/google/uuid"
)

func (h *Handler) GetAppRoles(w http.ResponseWriter, r *http.Request) {
	var approles []models.AppRole
	err := h.DB.Model(&approles).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approles)

}

func (h *Handler) CreateAppRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approle := &models.AppRole{
		ID:          fmt.Sprintf("%s-%s", "app", uuid.New().String()),
		Name:        req.Name,
		Description: req.Description,
		RoleID:      uuid.New().String(),
		SecretID:    uuid.New().String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if _, err := h.DB.Model(approle).Insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approle)

}
