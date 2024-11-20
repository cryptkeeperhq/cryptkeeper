package models

import (
	"time"
)

type CertificateAuthority struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CACert      []byte    `json:"ca_cert"`
	PrivateKey  []byte    `json:"private_key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RootCA struct {
	ID        int64     `json:"id"`
	CA        int64     `json:"ca_id"`
	RootCert  []byte    `json:"root_cert"`
	RootKey   []byte    `json:"root_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubCA struct {
	ID        int64     `json:"id"`
	PathID    string    `json:"path_id"`
	SubCACert []byte    `json:"sub_ca_cert"`
	SubCAKey  []byte    `json:"sub_ca_key"`
	RootCA    int64     `json:"root_ca_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CertificateTemplate struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	CommonName     string    `json:"common_name"`
	Organization   string    `json:"organization"`
	ValidityPeriod int       `json:"validity_period"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CertificateRequest struct {
	ID             int64     `json:"id"`
	CommonName     string    `json:"common_name"`
	Organization   string    `json:"organization"`
	ValidityPeriod int       `json:"validity_period"`
	CATemplateID   int64     `json:"ca_template_id"`
	CACertID       int64     `json:"ca_cert_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
