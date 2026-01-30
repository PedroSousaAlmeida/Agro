package middleware

import (
	"context"
	"net/http"
)

const (
	ClientIDKey contextKey = "clientID"
	UserIDKey   contextKey = "userID"
)

// ExtractTenancy extrai client_id e user_id do token JWT e adiciona ao context
func ExtractTenancy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(map[string]interface{})
		if !ok {
			http.Error(w, "No claims found", http.StatusUnauthorized)
			return
		}

		// Extrair client_id e user_id (sub) do token
		clientID, _ := claims["client_id"].(string)
		userID, _ := claims["sub"].(string)

		// Adicionar ao context
		ctx := r.Context()
		if clientID != "" {
			ctx = context.WithValue(ctx, ClientIDKey, clientID)
		}
		if userID != "" {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireClient valida que o usuário tem client_id no token
func RequireClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID, ok := r.Context().Value(ClientIDKey).(string)
		if !ok || clientID == "" {
			http.Error(w, "No client association", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
