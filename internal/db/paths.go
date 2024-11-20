package db

import (
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func GetPaths() ([]models.Path, error) {
	var paths []models.Path
	err := DB.Model(&paths).Select()
	return paths, err
}

func GetPathByName(name string) (models.Path, error) {

	var path models.Path
	err := DB.Model(&path).Where("path = ?", name).Select()
	return path, err
}

func GetPathByID(id string) (models.Path, error) {
	var path models.Path
	err := DB.Model(&path).Where("id = ?", id).First()
	return path, err
}
