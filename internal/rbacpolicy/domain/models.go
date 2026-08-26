// Package domain — сущности конфигурируемых RBAC-политик (хранятся в Postgres).
package domain

import "time"

// Role — каталог ролей.
type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Rule — строка матрицы прав (роль × ресурс × действие → зона владения).
type Rule struct {
	ID        int64     `json:"id"`
	Role      string    `json:"role"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	Scope     string    `json:"scope"`
	UpdatedBy *int64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutePolicy — определение маршрутной проверки (kind + параметры).
type RoutePolicy struct {
	Name      string         `json:"name"`
	Kind      string         `json:"kind"`
	Params    map[string]any `json:"params"`
	Active    bool           `json:"active"`
	UpdatedBy *int64         `json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at"`
}
