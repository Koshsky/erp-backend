// Package domain — entities of configurable RBAC policies (stored in Postgres).
package domain

import "time"

// Preset — a preset catalog entry (a named set of permissions).
type Preset struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PresetRule — a permissions matrix row (preset × resource × action → ownership scope).
type PresetRule struct {
	ID        int64     `json:"id"`
	Preset    string    `json:"preset"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	Scope     string    `json:"scope"`
	UpdatedBy *int64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserPermission — a per-user override of the preset rule for
// (resource, action): an explicit grant (Granted=true, Scope) or an explicit
// revoke (Granted=false, Scope ignored).
type UserPermission struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	Scope     string    `json:"scope"`
	Granted   bool      `json:"granted"`
	UpdatedBy *int64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutePolicy — a route policy definition (kind + parameters).
type RoutePolicy struct {
	Name      string         `json:"name"`
	Kind      string         `json:"kind"`
	Params    map[string]any `json:"params"`
	Active    bool           `json:"active"`
	UpdatedBy *int64         `json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at"`
}
