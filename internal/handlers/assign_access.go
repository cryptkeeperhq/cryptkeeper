package handlers

// func (h *Handler) AassignSecretAccess(w http.ResponseWriter, r *http.Request) {
// 	var tempRequest struct {
// 		SecretID    int64  `json:"secret_id"`
// 		UserID      string `json:"user_id,omitempty"`
// 		GroupID     string `json:"group_id,omitempty"`
// 		AccessLevel string `json:"access_level"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&tempRequest); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Convert string IDs to int64
// 	var userID, groupID *int64
// 	if tempRequest.UserID != "" {
// 		id, err := strconv.ParseInt(tempRequest.UserID, 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid user_id", http.StatusBadRequest)
// 			return
// 		}
// 		userID = &id
// 	}
// 	if tempRequest.GroupID != "" {
// 		id, err := strconv.ParseInt(tempRequest.GroupID, 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid group_id", http.StatusBadRequest)
// 			return
// 		}
// 		groupID = &id
// 	}

// 	if userID == nil && groupID == nil {
// 		http.Error(w, "Either user_id or group_id is required", http.StatusBadRequest)
// 		return
// 	}

// 	secretAccess := models.SecretAccess{
// 		SecretID:    tempRequest.SecretID,
// 		UserID:      userID,
// 		GroupID:     groupID,
// 		AccessLevel: tempRequest.AccessLevel,
// 	}

// 	_, err := h.DB.Model(&secretAccess).Insert()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }
