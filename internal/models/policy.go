package models

import (
	"time"
)

/*
	path "/zanzibar-kv" {
		permissions = ["list", "create", "update", "delete", "rotate"]
		users       = ["admin"]
		allowed_roles = ["admin"]

		secret "attribute/foo" {
			deny_permissions = ["update"]
			deny_users      = ["admin"]
		}

		secret "attribute/foo/bar" {
			deny_permissions = ["update"]
			deny_users      = ["admin"]
		}
	}
*/

type PolicyJson struct {
	PathId     string `json:"pathId"`
	Statements []struct {
		Principal  string   `json:"principal"`
		Actions    []string `json:"actions"`
		Effect     string   `json:"effect"`
		Conditions struct {
			IPAddress string `json:"ipAddress"`
		} `json:"conditions,omitempty"`
	} `json:"statements"`
}

type PolicyAuditLog struct {
	ID        int                    `json:"id"`
	PolicyID  string                 `json:"policy_id"`
	Action    string                 `json:"action"`
	Username  string                 `json:"username"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

type Policy struct {
	ID          string         `json:"id"`
	Name        string         `json:"name" `
	Description string         `json:"description"`
	HCL         string         `json:"hcl"` // Store HCL string
	PathID      string         `json:"path_id" pg:"fk:path_id"`
	Paths       []PolicyPath   `hcl:"path,block"`
	Secrets     []PolicySecret `hcl:"secret,block"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type PolicyPath struct {
	Name            string   `json:"name" hcl:"name,label"`
	Permissions     []string `hcl:"permissions,attr"`
	DenyPermissions []string `hcl:"deny_permissions,attr"`
	Users           []string `hcl:"users,attr"`
	Apps            []string `hcl:"apps,attr"`
	Groups          []string `hcl:"groups,attr"`
	Certificates    []string `hcl:"certificates,attr"`
}

type PolicySecret struct {
	Name             string    `json:"name" hcl:"name,label"`
	DenyPermissions  []string  `hcl:"deny_permissions,attr"`
	DenyUsers        *[]string `hcl:"deny_users,attr"`
	DenyApps         *[]string `hcl:"deny_apps,attr"`
	DenyGroups       *[]string `hcl:"deny_groups,attr"`
	DenyCertificates *[]string `hcl:"deny_certificates,attr"`
}

type UserPolicy struct {
	// ID           int64    `json:"id"`
	UserID       string   `json:"user_id" pg:"user_id,pk"`
	PolicyID     string   `json:"policy_id" pg:"policy_id,pk"`
	Capabilities []string `json:"capabilities"`
}

type GroupPolicy struct {
	// ID           int64    `json:"id"`
	GroupID      string   `json:"group_id" pg:"group_id,pk"`
	PolicyID     string   `json:"policy_id" pg:"policy_id,pk"`
	Capabilities []string `json:"capabilities"`
}

type AppPolicy struct {
	// ID           int64    `json:"id"`
	AppID        string   `json:"app_id" pg:"app_id,pk"`
	PolicyID     string   `json:"policy_id" pg:"policy_id,pk"`
	Capabilities []string `json:"capabilities"`
}

// type Policy struct {
// 	ID          int64         `json:"id"`
// 	Name        string        `json:"name"`
// 	Description string        `json:"description"`
// 	HCL         string        `json:"hcl"` // Store HCL string
// 	PathID      int64         `json:"path_id"`
// 	Paths       []PathRule    `hcl:"path,block"`
// 	Groups      []GroupRule   `hcl:"group,block"`
// 	Users       []UserRule    `hcl:"user,block"`
// 	AppRoles    []AppRoleRule `hcl:"approle,block"`
// 	CreatedAt   time.Time     `json:"created_at"`
// 	UpdatedAt   time.Time     `json:"updated_at"`
// }

// type PathRule struct {
// 	Path            string                 `hcl:"path,label"`
// 	Capabilities    []string               `hcl:"capabilities"`
// 	AllowParameters map[string]interface{} `hcl:"allowed_parameters,optional"`
// }

// type GroupRule struct {
// 	Name            string                 `hcl:"name,label"`
// 	Capabilities    []string               `hcl:"capabilities"`
// 	AllowParameters map[string]interface{} `hcl:"allowed_parameters,optional"`
// }

// type UserRule struct {
// 	Name            string                 `hcl:"name,label"`
// 	Capabilities    []string               `hcl:"capabilities"`
// 	AllowParameters map[string]interface{} `hcl:"allowed_parameters,optional"`
// }

// type AppRoleRule struct {
// 	Name            string                 `hcl:"name,label"`
// 	Capabilities    []string               `hcl:"capabilities"`
// 	AllowParameters map[string]interface{} `hcl:"allowed_parameters,optional"`
// }

// type Policy struct {
// 	ID          int64                 `json:"id"`
// 	Name        string                `json:"name"`
// 	Description string                `json:"description"`
// 	Rules       map[string]PolicyRule `json:"rules"` // JSON-encoded string representing the policy rules
// 	PathID      int64                 `json:"path_id"`
// 	CreatedAt   time.Time             `json:"created_at"`
// 	UpdatedAt   time.Time             `json:"updated_at"`
// }

// type PolicyRule struct {
// 	Groups   []string `json:"groups"`
// 	Users    []string `json:"users"`
// 	AppRoles []string `json:"app_roles"`
// }

// type AppRolesPolicy struct {
// 	ID           int64    `json:"id"`
// 	AppRoleID    int64    `json:"app_role_id"`
// 	PolicyID     int64    `json:"policy_id"`
// 	Capabilities []string `json:"capabilities"`
// }
