package ctx

import "context"

const (
	ContextKeyRole   = "role"
	ContextKeyUserID = "user_id"
)

// GetRole извлекает роль из контекста запроса (для использования в репозиториях/сервисах)
func GetRole(ctx context.Context) string {
	role, ok := ctx.Value(ContextKeyRole).(string)
	if !ok {
		return ""
	}
	return role
}

// GetUserID извлекает user_id из контекста запроса (для использования в репозиториях/сервисах)
func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(ContextKeyUserID).(int64)
	if !ok {
		return 0
	}
	return userID
}
