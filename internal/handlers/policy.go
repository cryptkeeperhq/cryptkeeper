package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/policy"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	var policies []models.Policy
	err := h.DB.Model(&policies).Select()
	if err != nil {
		log.Printf("Error fetching policies: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, policies)
}

func (h *Handler) GetPathPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathID := vars["pathID"]

	var policy models.Policy
	h.DB.Model(&policy).Where("path_id = ?", pathID).First()
	respondWithJSON(w, policy)
}

func (h *Handler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		respondDenied(w)
		return
	}

	vars := mux.Vars(r)
	policyID := vars["id"]

	err := db.DeletePolicy(policyID, identity.GetUsername())

	if err != nil {
		log.Printf("Error deleting policy: %v", err)
		respondWithError(w, err)
		return
	}

	respondWithOK(w)
}

func (h *Handler) SavePolicy(w http.ResponseWriter, r *http.Request) {

	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req models.Policy
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	policy, err := policy.ParseHCLPolicy(req.HCL)
	if err != nil {
		http.Error(w, "Failed to decode HCL", http.StatusBadRequest)
		return
	}
	policy.HCL = req.HCL
	policy.ID = req.ID

	policy.Name = req.Name
	policy.Description = req.Description
	policy.PathID = req.PathID

	for idx, secret := range policy.Secrets {
		if secret.DenyApps == nil {
			policy.Secrets[idx].DenyApps = &[]string{} // Default empty slice if needed
		}
		if secret.DenyCertificates == nil {
			policy.Secrets[idx].DenyCertificates = &[]string{} // Default empty slice if needed
		}
		if secret.DenyGroups == nil {
			policy.Secrets[idx].DenyGroups = &[]string{} // Default empty slice if needed
		}
		if secret.DenyUsers == nil {
			policy.Secrets[idx].DenyUsers = &[]string{} // Default empty slice if needed
		}
	}

	err = db.SavePolicy(policy, identity.GetUsername())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defaultPolicyJson, _ := json.Marshal(policy)
	if err := h.Producer.SendMessage("policy_updates", defaultPolicyJson); err != nil {
		log.Printf("Error sending message: %v", err)
	}

	respondWithJSON(w, policy)
}

func (h *Handler) GetPolicyAuditLogs(w http.ResponseWriter, r *http.Request) {
	var logs []models.PolicyAuditLog
	err := h.DB.Model(&logs).
		Order("timestamp DESC").
		Select()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, logs)
}
