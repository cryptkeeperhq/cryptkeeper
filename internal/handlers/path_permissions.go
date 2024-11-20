package handlers

// func (h *Handler) ListPaths(w http.ResponseWriter, r *http.Request) {
// 	userID, err := utils.GetUserIDFromContext(r.Context())
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusForbidden)
// 		return
// 	}

// 	// Get the groups the user is part of
// 	var groupIDs []int64
// 	err = h.DB.Model((*models.UserGroupRole)(nil)).
// 		Column("group_id").
// 		Where("user_id = ?", userID).
// 		Select(&groupIDs)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if len(groupIDs) == 0 {
// 		http.Error(w, "No groups found for the user", http.StatusNotFound)
// 		return
// 	}

// 	// Get the paths based on the group memberships
// 	var paths []string
// 	err = h.DB.Model((*models.PathPermission)(nil)).
// 		ColumnExpr("DISTINCT path").
// 		Where("group_id IN (?)", pg.In(groupIDs)).
// 		Select(&paths)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(paths)
// }

// // Create or update a path permission
// func (h *Handler) CreateOrUpdatePathPermission(w http.ResponseWriter, r *http.Request) {
// 	http.Error(w, "", http.StatusMethodNotAllowed)
// 	// var pathPermission models.PathPermission
// 	// if err := json.NewDecoder(r.Body).Decode(&pathPermission); err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// if pathPermission.GroupID == 0 || pathPermission.Path == "" || pathPermission.PermissionLevel == "" {
// 	// 	http.Error(w, "GroupID, Path, and PermissionLevel are required", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// // Check if the permission already exists
// 	// var existingPermission models.PathPermission
// 	// err := h.DB.Model(&existingPermission).
// 	// 	Where("group_id = ? AND path = ?", pathPermission.GroupID, pathPermission.Path).
// 	// 	Select()
// 	// if err != nil && err != pg.ErrNoRows {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// // _, err = h.DB.Model(&existingPermission).
// 	// // 	OnConflict("(group_id, path) DO UPDATE").
// 	// // 	Set("permission_level = EXCLUDED.permission_level").
// 	// // 	Insert()
// 	// // if err != nil {
// 	// // 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// // 	return
// 	// // }

// 	// // If the permission exists, update it
// 	// if err == nil {
// 	// 	pathPermission.ID = existingPermission.ID
// 	// 	pathPermission.CreatedAt = existingPermission.CreatedAt
// 	// 	pathPermission.UpdatedAt = time.Now()
// 	// 	_, err = h.DB.Model(&pathPermission).WherePK().Update()
// 	// } else {
// 	// 	pathPermission.CreatedAt = time.Now()
// 	// 	pathPermission.UpdatedAt = time.Now()
// 	// 	_, err = h.DB.Model(&pathPermission).Insert()
// 	// }

// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// w.Header().Set("Content-Type", "application/json")
// 	// json.NewEncoder(w).Encode(pathPermission)
// }

// // Get path permissions for a user
// func (h *Handler) GetPathPermissions(w http.ResponseWriter, r *http.Request) {
// 	groupID := r.URL.Query().Get("group_id")
// 	if groupID == "" {
// 		http.Error(w, "Group ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	var permissions []models.PathPermission
// 	err := h.DB.Model(&permissions).
// 		Where("group_id = ?", groupID).
// 		Select()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(permissions)
// }
