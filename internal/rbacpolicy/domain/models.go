// Package domain — сущности конфигурируемых RBAC-политик (хранятся в Postgres).
package domain

import "time"

// Role — каталог ролей.
type Role struct {
	ID          int64
	Name        string
	Description string
}

// Rule — строка матрицы прав (роль × ресурс × действие → зона владения).
type Rule struct {
	ID        int64
	Role      string
	Resource  string
	Action    string
	Scope     string
	UpdatedBy *int64
	UpdatedAt time.Time
}

// RoutePolicy — определение маршрутной проверки (kind + параметры).
type RoutePolicy struct {
	Name      string
	Kind      string
	Params    map[string]any
	Active    bool
	UpdatedBy *int64
	UpdatedAt time.Time
}
