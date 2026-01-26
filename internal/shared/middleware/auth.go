package middleware

import (
	"agro-monitoring/internal/config"
	"agro-monitoring/internal/shared/response"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// ClaimsKey is the key for storing user claims in the context.
const ClaimsKey contextKey = "userClaims"

// Authenticator holds the OIDC provider and verifier.
type Authenticator struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator(env *config.Env) (*Authenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), fmt.Sprintf("%s/realms/%s", env.KeycloakURL, env.KeycloakRealm))
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: env.KeycloakClientID})

	return &Authenticator{
		provider: provider,
		verifier: verifier,
	}, nil
}

// Auth is the middleware that validates JWT tokens.
func (a *Authenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondError(w, http.StatusUnauthorized, "Authorization header must be in 'Bearer <token>' format")
			return
		}
		tokenString := parts[1]

		idToken, err := a.verifier.Verify(r.Context(), tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Extract claims and add to context
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to parse claims")
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.ErrorResponse{Message: message})
}
