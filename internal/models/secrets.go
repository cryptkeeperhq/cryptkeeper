package models

import "time"

type Secret struct {
	ID             string                 `json:"id" pg:"id,pk"`
	PathID         string                 `json:"path_id" pg:"fk:path_id,unique:path_id_key_version"`
	Key            string                 `json:"key" pg:"key,pk,unique:path_id_key_version"`
	Version        int                    `json:"version" pg:"version,pk,unique:path_id_key_version"`
	EncryptedDEK   []byte                 `json:"encrypted_dek"`
	EncryptedValue []byte                 `json:"encrypted_value"`
	Checksum       string                 `json:"checksum"`
	Metadata       map[string]interface{} `json:"metadata"`
	IsMultiValue   bool                   `json:"is_multi_value"`

	Tags []string `json:"tags"`

	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	IsOneTime        bool       `json:"is_one_time"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	RotatedAt        *time.Time `json:"rotated_at,omitempty"`
	RotationInterval string     `json:"rotation_interval,omitempty"`
	LastRotatedAt    *time.Time `json:"last_rotated_at,omitempty"`
	CreatedBy        string     `json:"created_by"`

	Value string `pg:"-" json:"value"`

	Path          string `pg:"-" json:"path"`
	KeyType       string `pg:"-" json:"key_type"`
	CreatedByUser string `pg:"-" json:"created_by_user,omitempty"`
}

type SecretDeletion struct {
	ID             int64                  `json:"id"`
	SecretID       string                 `json:"secret_id" pg:"fk:secret_id"`
	PathID         string                 `json:"path_id" pg:"fk:path_id"`
	Key            string                 `json:"key"`
	Version        int                    `json:"version"`
	EncryptedDEK   []byte                 `json:"encrypted_dek"`   // Encrypted Data Encryption Key
	EncryptedValue []byte                 `json:"encrypted_value"` // Encrypted Secret Value
	Metadata       map[string]interface{} `json:"metadata"`
	DeletedAt      time.Time              `json:"deleted_at"`
	// Accesses       []SecretAccess         `json:"accesses"`
}

// type SecretAccess struct {
// 	tableName   struct{} `pg:"secret_accesses"` // Add the correct table name
// 	ID          int64    `json:"id"`
// 	SecretID    int64    `json:"secret_id" pg:"fk:secret_id"`
// 	UserID      *int64   `json:"user_id,omitempty"`
// 	GroupID     *int64   `json:"group_id,omitempty"`
// 	AccessLevel string   `json:"access_level"`
// }

// Struct for shared link
type SharedLink struct {
	ID        int64     `json:"id"`
	LinkID    string    `json:"link_id"`
	SecretID  string    `json:"secret_id" pg:"fk:secret_id"`
	Version   int       `json:"version"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ApprovalRequest struct {
	ID       int64   `json:"id"`
	UserID   string  `json:"user_id"`
	SecretID *string `json:"secret_id,omitempty"`
	Action   string  `json:"action"`
	Status   string  `json:"status"`
	// Details   map[string]interface{} `json:"details"`
	Details   Secret    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Notification struct {
	ID        int64     `json:"id"`
	SecretID  string    `json:"secret_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID        int64                  `json:"id"`
	Username  string                 `json:"username"`
	Action    string                 `json:"action"`
	SecretID  *string                `json:"secret_id,omitempty"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time              `json:"timestamp"`
	// Username  string                 `json:"username" pg:"-"`
}
