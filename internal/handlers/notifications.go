package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	var notifications []models.Notification
	err := h.DB.Model(&notifications).Limit(100).Order("created_at DESC").Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
