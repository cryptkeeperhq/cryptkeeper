package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
)

func (h *Handler) ScanForSecrets(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secrets := utils.ScanForSecrets(request.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}
