package handler

import (
	"encoding/json"
	"net/http"

	"agro-monitoring/internal/shared/middleware"
	"agro-monitoring/internal/shared/response"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles user-related requests.
type UserHandler struct{}

// NewUserHandler creates a new UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// RegisterRoutes registers the user routes.
// The grouping and middleware are handled in bootstrap/routes.go.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/me", h.Me)
}

// Me is a handler to get the current user's information.
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Retrieve claims from context
	claims, ok := r.Context().Value(middleware.ClaimsKey).(map[string]interface{})
	if !ok {
		respondError(w, http.StatusUnauthorized, "No user claims found in context")
		return
	}

	respondJSON(w, http.StatusOK, response.NewSuccessResponse(claims))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, response.ErrorResponse{Message: message})
}
