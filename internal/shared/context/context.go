package context

import (
	"context"

	"agro-monitoring/internal/shared/middleware"
)

// GetClientID extrai o client_id do context
func GetClientID(ctx context.Context) (string, bool) {
	clientID, ok := ctx.Value(middleware.ClientIDKey).(string)
	return clientID, ok
}

// GetUserID extrai o user_id (sub) do context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	return userID, ok
}
