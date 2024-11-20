package handlers

import (
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func (h *Handler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	var summary struct {
		SecretsCount  int `json:"secrets_count"`
		PathsCount    int `json:"paths_count"`
		PoliciesCount int `json:"policies_count"`
	}

	var err error
	if summary.SecretsCount, err = h.DB.Model((*models.Secret)(nil)).Count(); err != nil {
		http.Error(w, "Failed to fetch secrets count", http.StatusInternalServerError)
		return
	}

	if summary.PathsCount, err = h.DB.Model((*models.Path)(nil)).Count(); err != nil {
		http.Error(w, "Failed to fetch paths count", http.StatusInternalServerError)
		return
	}

	if summary.PoliciesCount, err = h.DB.Model((*models.Policy)(nil)).Count(); err != nil {
		http.Error(w, "Failed to fetch policies count", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, summary)

}

func (h *Handler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	var logs []models.AuditLog
	if err := h.DB.Model(&logs).Order("timestamp DESC").Limit(10).Select(); err != nil {
		http.Error(w, "Failed to fetch recent activity", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, logs)
}
