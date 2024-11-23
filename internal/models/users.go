package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" pg:"id,pk"`
	Username  string    `json:"username" pg:"username,pk"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID        string    `json:"id" pg:"id,pk"`
	Name      string    `json:"name" pg:"name,pk"`
	CreatedAt time.Time `json:"created_at"`
}

type UserGroup struct {
	UserID  string `json:"user_id" pg:"user_id,pk"`
	GroupID string `json:"group_id" pg:"group_id,pk"`
}

type Certificate struct {
	Name           string    `json:"name" pg:"name,pk"`
	Description    string    `json:"description"`
	DNSName        string    `json:"dns_names"`
	IPAddresses    string    `json:"ip_addresses"`
	EmailAddresses string    `json:"email_addresses"`
	Key            string    `json:"key"`
	Cert           string    `json:"cert"`
	CreatedAt      time.Time `json:"created_at"`
	Revoked        bool      `json:"revoked"`
}
