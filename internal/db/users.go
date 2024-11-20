package db

import (
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func GetUser(username string) (models.User, error) {
	var user models.User
	err := DB.Model(&user).Where("username = ?", username).First()
	return user, err
}

func GetUserGroups(userID string) ([]string, error) {
	var userGroups []string
	err := DB.Model((*models.UserGroup)(nil)).
		Join("JOIN groups p ON p.id = group_id").
		Column("p.name").Where("user_id = ?", userID).Select(&userGroups)
	return userGroups, err
}
