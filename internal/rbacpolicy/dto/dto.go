// Package dto — API contracts for RBAC policy administration.
package dto

import "time"

// RuleInput — a permissions matrix row record.
type RuleInput struct {
	Role     string `json:"role"     example:"vp"     binding:"required"`
	Resource string `json:"resource" example:"task"   binding:"required"`
	Action   string `json:"action"   example:"update" binding:"required"`
	Scope    string `json:"scope"    example:"parent" binding:"required"`
}

// RuleView — a permissions matrix row.
type RuleView struct {
	ID        int64     `json:"id"`
	Role      string    `json:"role"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	Scope     string    `json:"scope"`
	UpdatedBy *int64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutePolicyInput — a route policy record (kind + parameters).
type RoutePolicyInput struct {
	Name   string         `json:"name"   example:"task.create" binding:"required"`
	Kind   string         `json:"kind"   example:"create"      binding:"required"`
	Params map[string]any `json:"params"`
	Active *bool          `json:"active"`
}

// RoutePolicyView — a route policy view.
type RoutePolicyView struct {
	Name      string         `json:"name"`
	Kind      string         `json:"kind"`
	Params    map[string]any `json:"params"`
	Active    bool           `json:"active"`
	UpdatedBy *int64         `json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// MatrixCell — an effective matrix cell (with the admin bypass).
type MatrixCell struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

// ExplainInput — parameters of the "why allow/deny" check.
type ExplainInput struct {
	Role         string `form:"role"          binding:"required" example:"vp"`
	Resource     string `form:"resource"      binding:"required" example:"task"`
	Action       string `form:"action"        binding:"required" example:"update"`
	UserID       int64  `form:"user_id"                          example:"4"`
	ProjectOwner int64  `form:"project_owner"                    example:"3"`
	ProcessOwner int64  `form:"process_owner"                    example:"4"`
	Owner        int64  `form:"owner"`
}

// ExplainResult — a check result.
type ExplainResult struct {
	Scope   string `json:"scope"`
	Allowed bool   `json:"allowed"`
}

// Permission — a role's principal right: an action on a resource is allowed
// (the scope from the matrix; the frontend derives the ownership zone from it).
type Permission struct {
	Resource string `json:"resource" example:"project"`
	Action   string `json:"action"   example:"create"`
	Scope    string `json:"scope"    example:"own"`
}

// RoleUpsertInput — role creation (name = system access code).
type RoleUpsertInput struct {
	Name        string `json:"name"        example:"auditor"       binding:"required"`
	Description string `json:"description" example:"Внешний аудит"`
}

// RoleUpdateInput — role description update.
type RoleUpdateInput struct {
	Description string `json:"description" example:"Внешний аудит"`
}
