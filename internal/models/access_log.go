package models

import "time"

type AccessLog struct {
	ID         int64     `json:"id"`
	SecretID   string    `json:"secret_id"`
	UserID     string    `json:"user_id"`
	AccessedAt time.Time `json:"accessed_at"`
	Source     string    `json:"source"`
	Username   string    `pg:"-" json:"username"`
}

type SchedulerMetadata struct {
	ID            int64     `json:"id"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}
