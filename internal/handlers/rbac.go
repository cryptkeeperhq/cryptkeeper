package handlers

// // func (h *Handler) CheckPermission(permissionName string, next http.Handler) http.Handler {
// // 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		userID, err := utils.GetUserIDFromContext(r.Context())
// // 		if err != nil {
// // 			http.Error(w, err.Error(), http.StatusForbidden)
// // 			return
// // 		}

// // 		fmt.Println("USER ID", userID)

// // 		var permission models.Permission
// // 		err = h.DB.Model(&permission).
// // 			Join("JOIN role_permissions rp ON rp.permission_id = permission.id").
// // 			Join("JOIN user_roles ur ON ur.role_id = rp.role_id").
// // 			Where("ur.user_id = ? AND permission.name = ?", userID, permissionName).
// // 			Select()
// // 		if err != nil {
// // 			http.Error(w, "Access denied", http.StatusForbidden)
// // 			return
// // 		}

// // 		next.ServeHTTP(w, r)
// // 	})
// // }

// func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
// 	var roles []models.Role
// 	err := h.DB.Model(&roles).Select()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(roles)
// }

// func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
// 	var role models.Role
// 	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if _, err := h.DB.Model(&role).Insert(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(role)
// }

// func (h *Handler) GetPermissions(w http.ResponseWriter, r *http.Request) {
// 	var permissions []models.Permission
// 	err := h.DB.Model(&permissions).Select()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(permissions)
// }

// func (h *Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
// 	var permission models.Permission
// 	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if _, err := h.DB.Model(&permission).Insert(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(permission)
// }

// func (h *Handler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
// 	var request struct {
// 		RoleID string `json:"role_id"`
// 		UserID string `json:"user_id"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var userRole models.UserRole
// 	if request.UserID != "" {
// 		id, err := strconv.Atoi(request.UserID)
// 		if err != nil {
// 			http.Error(w, "Invalid user_id", http.StatusBadRequest)
// 			return
// 		}
// 		userRole.UserID = id
// 	}

// 	if request.RoleID != "" {
// 		id, err := strconv.Atoi(request.RoleID)
// 		if err != nil {
// 			http.Error(w, "Invalid role_id", http.StatusBadRequest)
// 			return
// 		}
// 		userRole.RoleID = id
// 	}

// 	// var userRole models.UserRole
// 	// if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 	// 	return
// 	// }

// 	if _, err := h.DB.Model(&userRole).Insert(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(userRole)
// }

// func (h *Handler) AssignPermissionToRole(w http.ResponseWriter, r *http.Request) {
// 	var request struct {
// 		RoleID       string `json:"role_id"`
// 		PermissionID string `json:"permission_id"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var rolePermission models.RolePermission
// 	if request.RoleID != "" {
// 		id, err := strconv.Atoi(request.RoleID)
// 		if err != nil {
// 			http.Error(w, "Invalid user_id", http.StatusBadRequest)
// 			return
// 		}
// 		rolePermission.RoleID = id
// 	}

// 	if request.PermissionID != "" {
// 		id, err := strconv.Atoi(request.PermissionID)
// 		if err != nil {
// 			http.Error(w, "Invalid role_id", http.StatusBadRequest)
// 			return
// 		}
// 		rolePermission.PermissionID = id
// 	}

// 	if _, err := h.DB.Model(&rolePermission).Insert(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(rolePermission)
// }
