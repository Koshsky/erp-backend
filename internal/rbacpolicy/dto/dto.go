// Package dto — контракты API администрирования RBAC-политик.
package dto

import "time"

// RuleInput — запись строки матрицы прав.
type RuleInput struct {
	Role     string `json:"role"     example:"vp"     binding:"required"`
	Resource string `json:"resource" example:"task"   binding:"required"`
	Action   string `json:"action"   example:"update" binding:"required"`
	Scope    string `json:"scope"    example:"parent" binding:"required"`
}

// RuleView — строка матрицы прав.
type RuleView struct {
	ID        int64     `json:"id"`
	Role      string    `json:"role"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	Scope     string    `json:"scope"`
	UpdatedBy *int64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutePolicyInput — запись маршрутной проверки (kind + параметры).
type RoutePolicyInput struct {
	Name   string         `json:"name"   example:"task.create" binding:"required"`
	Kind   string         `json:"kind"   example:"create"      binding:"required"`
	Params map[string]any `json:"params"`
	Active *bool          `json:"active"`
}

// RoutePolicyView — маршрутная проверка.
type RoutePolicyView struct {
	Name      string         `json:"name"`
	Kind      string         `json:"kind"`
	Params    map[string]any `json:"params"`
	Active    bool           `json:"active"`
	UpdatedBy *int64         `json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// MatrixCell — клетка эффективной матрицы (с admin-байпасом).
type MatrixCell struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

// ExplainInput — параметры проверки «почему allow/deny».
type ExplainInput struct {
	Role         string `form:"role"          binding:"required" example:"vp"`
	Resource     string `form:"resource"      binding:"required" example:"task"`
	Action       string `form:"action"        binding:"required" example:"update"`
	UserID       int64  `form:"user_id"                          example:"4"`
	ProjectOwner int64  `form:"project_owner"                    example:"3"`
	ProcessOwner int64  `form:"process_owner"                    example:"4"`
	Owner        int64  `form:"owner"`
}

// ExplainResult — результат проверки.
type ExplainResult struct {
	Scope   string `json:"scope"`
	Allowed bool   `json:"allowed"`
}

// Permission — принципиальное право роли: действие над ресурсом разрешено
// (скоуп доступа из матрицы; по нему фронт понимает зону владения).
type Permission struct {
	Resource string `json:"resource" example:"project"`
	Action   string `json:"action" example:"create"`
	Scope    string `json:"scope" example:"own"`
}
