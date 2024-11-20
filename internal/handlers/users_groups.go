package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) UserInGroup(userID int64, groupName string) bool {
	var group []models.Group
	err := h.DB.Model(&group).
		Column("g.*").
		TableExpr("groups AS g").
		Join("JOIN user_groups ugr ON ugr.group_id = g.id").
		Where("ugr.user_id = ? AND g.name = ?", userID, groupName).
		Select()

	fmt.Println("Group", group)
	return err == nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.KAuth.CreateUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// user.Password = string(hashedPassword)
	// user.CreatedAt = time.Now()

	// _, err = h.DB.Model(&user).Insert()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Config.Auth.SSOEnabled {
		err := h.KAuth.CreateGroup(group.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		group.ID = uuid.New().String()
		group.CreatedAt = time.Now()
		_, err := h.DB.Model(&group).Insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// type KUserGroup struct {
// 	UserID  string `json:"user_id"`
// 	GroupID string `json:"group_id"`
// }

func (h *Handler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	var userGroup models.UserGroup
	if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Config.Auth.SSOEnabled {
		err := h.KAuth.AddUserToGroup(userGroup.UserID, userGroup.GroupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := h.DB.Model(&userGroup).Insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// var g models.Group
		// h.DB.Model(&g).Where("id = ?", userGroup.GroupID).First()

		// var u models.User
		// h.DB.Model(&u).Where("id = ?", userGroup.UserID).First()

		// h.Z.AddUserToGroup(u.Username, g.Name)

	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	// var userGroup models.UserGroup
	// if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// _, err := h.DB.Model(&userGroup).Where("user_id = ? AND group_id = ?", userGroup.UserID, userGroup.GroupID).Delete()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	var userGroup models.UserGroup
	if err := json.NewDecoder(r.Body).Decode(&userGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// vars := mux.Vars(r)
	// userID := vars["userID"]

	err := h.KAuth.RemoveUserFromGroup(userGroup.UserID, userGroup.GroupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {

	if h.Config.Auth.SSOEnabled {
		users, err := h.KAuth.GetUsers("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		return
	}

	var users []models.User
	err := h.DB.Model(&users).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, users)
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {

	if h.Config.Auth.SSOEnabled {
		groups, err := h.KAuth.GetGroups("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}

	var groups []models.Group
	err := h.DB.Model(&groups).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *Handler) ListGroupUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	var users []models.User

	if h.Config.Auth.SSOEnabled {
		kUsers, err := h.KAuth.GetUsers(groupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, user := range kUsers {
			users = append(users, models.User{
				ID:       *user.ID,
				Username: *user.Username,
			})
		}
	} else {
		err := h.DB.Model(&users).Distinct().
			Column("u.id", "u.username", "u.created_at").
			TableExpr("users AS u").
			Join("JOIN user_groups AS ug ON ug.user_id = u.id").
			Where("ug.group_id = ?", groupID).
			Select()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var groups []models.Group

	if h.Config.Auth.SSOEnabled {
		kGroups, err := h.KAuth.GetGroups(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, group := range kGroups {
			groups = append(groups, models.Group{
				ID:   *group.ID,
				Name: *group.Name,
			})
		}
	} else {
		err := h.DB.Model(&groups).Distinct().
			Column("g.id", "g.name", "g.created_at").
			TableExpr("groups AS g").
			Join("JOIN user_groups AS ug ON ug.group_id = g.id").
			Where("ug.user_id = ?", userID).
			Select()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// func (h *Handler) AssignGroupAccess(w http.ResponseWriter, r *http.Request) {
// 	var access models.SecretAccess
// 	if err := json.NewDecoder(r.Body).Decode(&access); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if access.GroupID == nil {
// 		http.Error(w, "Group ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	_, err := h.DB.Model(&access).Insert()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *Handler) getUserGroups(userID int64) ([]int64, error) {
// 	var groupIDs []int64
// 	err := h.DB.Model((*models.UserGroup)(nil)).
// 		Column("group_id").
// 		Where("user_id = ?", userID).
// 		Select(&groupIDs)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return groupIDs, nil
// }
