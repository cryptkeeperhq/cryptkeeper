package mhandlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/go-pg/pg/v10"
	// "your_project/messaging" // replace with correct import path
)

// PolicyUpdateHandler is an implementation of the MessageHandler interface.
type PolicyUpdateHandler struct {
	DB *pg.DB
}

// HandleMessage processes the incoming message from the specified topic.
func (h *PolicyUpdateHandler) HandleMessage(topic string, message []byte) error {
	log.Printf("Processing message from topic %s: %s", topic, string(message))
	fmt.Printf("Message from topic '%s': %s\n", topic, message)

	var policy models.Policy
	err := json.Unmarshal(message, &policy)
	if err != nil {
		return err
	}

	// log.Printf("Received policy update for topic %s: %+v", topic, policy)

	// Process the policy update using Zanzibar and the database
	h.applyHCLConfig(policy)

	// Returning nil to indicate successful handling.
	return nil
}

func (h *PolicyUpdateHandler) applyHCLConfig(policy models.Policy) {

	h.DB.Model((*models.UserPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()
	h.DB.Model((*models.GroupPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()
	h.DB.Model((*models.AppPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()

	for _, path := range policy.Paths {

		fmt.Println("-----------> ", path)

		// TODO: This does not work for users which are added after the policy is created
		for _, user := range path.Users {
			// h.Z.AddPathPermissions(path.ID, user, path.Permissions)
			var u models.User
			err := h.DB.Model(&u).Where("username = ?", user).First()
			if err != nil {
				continue
			}
			userPolicy := models.UserPolicy{
				UserID:       u.ID,
				PolicyID:     policy.ID,
				Capabilities: path.Permissions,
			}
			h.DB.Model(&userPolicy).OnConflict("(user_id, policy_id) DO UPDATE").Insert()
		}

		for _, app := range path.Apps {
			// h.Z.AddPathPermissions(path.ID, app, path.Permissions)
			var u models.AppRole
			err := h.DB.Model(&u).Where("ID = ?", app).First()
			if err != nil {
				continue
			}
			appPolicy := models.AppPolicy{
				AppID:        u.ID,
				PolicyID:     policy.ID,
				Capabilities: path.Permissions,
			}
			fmt.Println("APP POLICY", appPolicy)
			h.DB.Model(&appPolicy).OnConflict("(app_id, policy_id) DO UPDATE").Insert()
		}

		for _, group := range path.Groups {
			// h.Z.AddPathGroupPermissions(path.ID, group, path.Permissions)

			var g models.Group
			err := h.DB.Model(&g).Where("name = ?", group).First()
			if err != nil {
				continue
			}

			groupPolicy := models.GroupPolicy{
				GroupID:      g.ID,
				PolicyID:     policy.ID,
				Capabilities: path.Permissions,
			}
			h.DB.Model(&groupPolicy).OnConflict("(group_id, policy_id) DO UPDATE").Insert()
		}
	}

	// for _, secret := range policy.Secrets {
	// 	// h.Z.AddSecretToPath(secret.PathID, secret.ID)
	// 	for _, user := range secret.DenyUsers {
	// 		for _, perm := range secret.DenyPermissions {
	// 			h.Z.DenySecretPermission(secret.ID, user, "deny_"+perm)
	// 		}
	// 	}
	// }

}
