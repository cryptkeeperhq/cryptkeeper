package models

import "time"

type AppRole struct {
	ID          string    `json:"id" pg:"fk:id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RoleID      string    `json:"role_id" pg:"fk:role_id"`
	SecretID    string    `json:"secret_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
