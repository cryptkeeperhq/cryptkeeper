// handlers/lineage.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

type LineageNode struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	ParentID string `json:"parent_id,omitempty"`
}

type LineageEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type LineageResponse struct {
	Nodes []LineageNode `json:"nodes"`
	Edges []LineageEdge `json:"edges"`
}

func (h *Handler) GetSecretLineage(w http.ResponseWriter, r *http.Request) {
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

	var accessLogs []models.AccessLog
	err = h.DB.Model(&accessLogs).
		ColumnExpr("access_log.*").
		ColumnExpr("u.username as username").
		Where("secret_id = ?", secret.ID).
		Join("LEFT JOIN users u ON u.id = access_log.user_id").
		Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	edges := []LineageEdge{}

	nodes := []LineageNode{
		{ID: fmt.Sprintf("secret-%d", secret.ID), Label: fmt.Sprintf("%d/%s", secret.PathID, secret.Key), Type: "secret"},
	}

	if secret.CreatedBy != "" {
		var creator models.User
		err = h.DB.Model(&creator).Where("id = ?", secret.CreatedBy).Select()
		if err == nil {
			nodes = append(nodes, LineageNode{ID: fmt.Sprintf("user-%d", secret.ID), Label: creator.Username, Type: "create", ParentID: fmt.Sprintf("secret-%d", secret.ID)})
			edges = append(edges, LineageEdge{ID: fmt.Sprintf("e%d-%d", creator.ID, secret.ID), Source: fmt.Sprintf("user-%d", secret.ID), Target: fmt.Sprintf("secret-%d", secret.ID)})
		}
	}

	for _, log := range accessLogs {
		node := LineageNode{
			ID:       fmt.Sprintf("log-%d", log.ID),
			Label:    fmt.Sprintf("%s Accessed on %s", log.Username, log.AccessedAt.String()),
			Type:     "access",
			ParentID: fmt.Sprintf("secret-%d", secret.ID),
		}
		nodes = append(nodes, node)
		edges = append(edges, LineageEdge{ID: fmt.Sprintf("e%d-%d", log.ID, secret.ID), Source: fmt.Sprintf("log-%d", log.ID), Target: fmt.Sprintf("secret-%d", secret.ID)})
	}

	response := LineageResponse{Nodes: nodes, Edges: edges}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
