package models

import (
	"time"
)

type Path struct {
	ID         string                 `json:"id" pg:"id,pk"`
	Path       string                 `json:"path" pg:"path,pk,unique"`
	EngineType string                 `json:"engine_type"`
	Metadata   map[string]interface{} `json:"metadata"`

	// TODO: Change to string
	CreatedBy int64 `json:"created_by"`

	KeyData []byte `json:"key_data"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
