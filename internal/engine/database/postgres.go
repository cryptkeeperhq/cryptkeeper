package database

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

func CreateRoleInDatabase(connectionString, roleTemplate, dbUsername, dbPassword string) error {
	db, err := connectToDatabase(connectionString)
	if err != nil {
		return err
	}

	err = executeRoleTemplate(db, roleTemplate, dbUsername, dbPassword)
	if err != nil {
		return err
	}

	return nil

}

func connectToDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func executeRoleTemplate(db *sql.DB, roleTemplate, username, password string) error {
	query := strings.ReplaceAll(roleTemplate, "{{name}}", username)
	query = strings.ReplaceAll(query, "{{password}}", password)
	_, err := db.Exec(query)
	return err
}
