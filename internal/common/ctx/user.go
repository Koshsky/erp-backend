package ctx

import "context"

type contextKey string

const (
	ContextKeyRole   contextKey = "role"
	ContextKeyUserID contextKey = "user_id"
)

// GetRole extracts the role from the request context (for use in repositories/services).
func GetRole(ctx context.Context) string {
	role, ok := ctx.Value(ContextKeyRole).(string)
	if !ok {
		return ""
	}
	return role
}

// GetUserID extracts the user_id from the request context (for use in repositories/services).
func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(ContextKeyUserID).(int64)
	if !ok {
		return 0
	}
	return userID
}

// SetRole writes the role to the context.
func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyRole, role)
}

// SetUserID writes the user_id to the context.
func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}
