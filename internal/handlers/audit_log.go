package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func (h *Handler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	username := r.URL.Query().Get("username")
	action := r.URL.Query().Get("action")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var logs []models.AuditLog
	query := h.DB.Model(&logs).
		ColumnExpr("audit_log.*").
		Order("timestamp DESC").
		Limit(limit).
		Offset((page - 1) * limit)

	if username != "" {
		query = query.Where("username = ?", username)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("timestamp BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.DB.Model(&models.AuditLog{}).Count()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Logs  []models.AuditLog `json:"logs"`
		Total int               `json:"total"`
	}{
		Logs:  logs,
		Total: totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
