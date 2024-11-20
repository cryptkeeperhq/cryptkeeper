package db

import (
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/google/uuid"
)

func logPolicyChange(policyID string, action, username string, details map[string]interface{}) error {
	auditLog := models.PolicyAuditLog{
		PolicyID: policyID,
		Action:   action,
		Username: username,
		Details:  details,
	}
	_, err := DB.Model(&auditLog).Insert()
	return err
}

func DeletePolicy(id string, username string) error {

	var policy models.Policy
	DB.Model(&policy).Where("id = ?", id).First()

	// Delete the policy
	_, err := DB.Model((*models.Policy)(nil)).
		Where("id = ?", id).
		Delete()

	if err != nil {
		return err
	}

	details := map[string]interface{}{
		"policy": policy,
	}
	logPolicyChange(policy.ID, "delete", username, details)
	return nil
}

func SavePolicy(policy models.Policy, username string) error {

	details := map[string]interface{}{
		"policy": policy, // Add relevant policy details here
	}

	if policy.ID == "" {
		policy.ID = uuid.New().String()
		_, err := DB.Model(&policy).Insert()
		if err != nil {
			return err
		}
		logPolicyChange(policy.ID, "create", username, details)
	} else {
		_, err := DB.Model(&policy).WherePK().Update()
		if err != nil {
			return err
		}
		logPolicyChange(policy.ID, "update", username, details)
	}

	return nil

	// var err error

	// // Remove existing group and user policies for this policy
	// _, err = DB.Model((*models.GroupPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()
	// if err != nil {
	// 	return err
	// }

	// _, err = DB.Model((*models.UserPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()
	// if err != nil {
	// 	return err
	// }

	// _, err = DB.Model((*models.AppRolesPolicy)(nil)).Where("policy_id = ?", policy.ID).Delete()
	// if err != nil {
	// 	return err
	// }

	// // Add group policies
	// for _, groupPolicy := range policy.Groups {
	// 	var group models.Group
	// 	err := DB.Model(&group).Where("name = ?", groupPolicy.Name).Select()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	groupPolicy := models.GroupPolicy{
	// 		GroupID:      group.ID,
	// 		PolicyID:     policy.ID,
	// 		Capabilities: groupPolicy.Capabilities,
	// 	}

	// 	_, err = DB.Model(&groupPolicy).Insert()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // Add user policies
	// for _, userPolicy := range policy.Users {
	// 	var user models.User
	// 	err := DB.Model(&user).Where("username = ?", userPolicy.Name).Select()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	userPolicy := models.UserPolicy{
	// 		UserID:       user.ID,
	// 		PolicyID:     policy.ID,
	// 		Capabilities: userPolicy.Capabilities,
	// 	}
	// 	_, err = DB.Model(&userPolicy).Insert()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // Add App Role Policies
	// for _, appRolePolicy := range policy.AppRoles {
	// 	var appRole models.AppRole
	// 	err := DB.Model(&appRole).Where("role_id = ?", appRolePolicy.Name).Select()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	appRolePolicy := models.AppRolesPolicy{
	// 		AppRoleID:    appRole.ID,
	// 		PolicyID:     policy.ID,
	// 		Capabilities: appRolePolicy.Capabilities,
	// 	}
	// 	_, err = DB.Model(&appRolePolicy).Insert()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// return nil
}
