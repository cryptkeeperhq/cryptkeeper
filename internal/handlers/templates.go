package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	templates := []map[string]interface{}{
		{
			"name": "Database Credentials",
			"fields": map[string]string{
				"username": "db_user",
				"password": "db_pass",
				"host":     "localhost",
				"port":     "5432",
				"database": "example_db",
			},
		},
		{
			"name": "API Key",
			"fields": map[string]string{
				"api_key": "your_api_key_here",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}
