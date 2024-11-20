package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/events"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WorkflowRequest struct {
	Event    map[string]interface{} `json:"event"`
	Workflow models.Workflow        `json:"workflow"`
}

// type Event struct {
// 	Name       string   `json:"name"`
// 	Attributes []string `json:"attributes"`
// }

// var events []Event

// func init() {
// 	// Define a list of events with predefined attributes
// 	events = []Event{
// 		{Name: "journey", Attributes: []string{"userJoined", "first_name", "last_name"}},
// 		{Name: "karma", Attributes: []string{"userCreated", "first_name", "last_name"}},
// 		// Add more events with attributes as needed
// 	}
// }

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events := events.GetEventList()
	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *Handler) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	var request WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// requestMap := make(map[string]interface{})
	// requestMap["name"] = "vishal"

	// Execute the workflow with the provided event and workflowJson
	// go func(request WorkflowRequest) {
	t := workflow.Temporal{}
	result, err := t.Execute(request.Workflow, request.Event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(result)
	// }(requestWorkflow)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	// http://localhost:8080/namespaces/default/workflows/3ebc76f6-bb60-4c53-862a-c41ede651e3a/e329b757-02fe-490c-9992-3eefc6ffb9d0/history
	w.Write([]byte(result))
}

func (h *Handler) GetWorkflows(w http.ResponseWriter, r *http.Request) {

	var workflow []models.Workflow
	h.DB.Model(&workflow).Select()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

func (h *Handler) SaveOrCreateWorkflow(w http.ResponseWriter, r *http.Request) {
	identity, ok := r.Context().Value("identity").(utils.Identity)
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workflow.UpdatedAt = time.Now()

	var err error
	if workflow.ID == "" {
		workflow.ID = uuid.New().String()
		workflow.CreatedAt = time.Now()
		workflow.CreatedBy = identity.GetID()
		_, err = h.DB.Model(&workflow).Insert()
	} else {
		_, err = h.DB.Model(&workflow).Where("workflow.id = ?", workflow.ID).Update()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

func (h *Handler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var workflowObj models.Workflow
	query := h.DB.Model(&workflowObj).
		Where("workflow.id = ?", id)

	query.Limit(1).Select()
	h.DB.Model(&workflowObj).Select()

	t := workflow.Temporal{}
	t.GetWorkflowHistory(workflowObj)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflowObj)
}
