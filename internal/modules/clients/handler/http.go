package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"agro-monitoring/internal/config"
	"agro-monitoring/internal/modules/clients/dto"
	"agro-monitoring/internal/modules/clients/usecase"
	sharedContext "agro-monitoring/internal/shared/context"
	sharedErrors "agro-monitoring/internal/shared/errors"
	"agro-monitoring/internal/shared/response"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	clientUC usecase.ClientUseCase
	env      *config.Env
}

func NewHandler(clientUC usecase.ClientUseCase, env *config.Env) *Handler {
	return &Handler{
		clientUC: clientUC,
		env:      env,
	}
}

// RegisterRoutes registra as rotas do módulo clients
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Rotas protegidas (/v1/clients/me)
	r.Route("/clients", func(r chi.Router) {
		r.Get("/me", h.GetMyClient)
		r.Get("/me/stats", h.GetMyStats)
		r.Get("/me/users", h.ListMyUsers)
	})
}

// RegisterAdminRoutes registra rotas admin (/v1/admin/clients)
func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.Route("/clients", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Get("/{id}/stats", h.GetStats)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// Create cria um novo client (admin)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := h.clientUC.CreateClient(r.Context(), req)
	if err != nil {
		if err == sharedErrors.ErrInvalidSlug {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		fmt.Printf("[ERROR] Failed to create client: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Failed to create client")
		return
	}

	resp := dto.ToClientResponse(client, h.env.AppBaseURL)
	respondJSON(w, http.StatusCreated, response.NewSuccessResponse(resp))
}

// List lista todos os clients (admin)
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	clients, total, err := h.clientUC.ListClients(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list clients")
		return
	}

	var resp []dto.ClientResponse
	for _, client := range clients {
		resp = append(resp, *dto.ToClientResponse(client, h.env.AppBaseURL))
	}

	data := map[string]interface{}{
		"clients": resp,
		"total":   total,
		"page":    page,
	}
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(data))
}

// GetByID busca client por ID (admin)
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	client, err := h.clientUC.GetClient(r.Context(), id)
	if err != nil {
		if err == sharedErrors.ErrClientNotFound {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get client")
		return
	}

	resp := dto.ToClientResponse(client, h.env.AppBaseURL)
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(resp))
}

// GetStats retorna estatísticas de um client (admin)
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	stats, err := h.clientUC.GetClientStats(r.Context(), id)
	if err != nil {
		if err == sharedErrors.ErrClientNotFound {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	resp := dto.ToClientStatsResponse(stats, h.env.AppBaseURL)
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(resp))
}

// Update atualiza um client (admin)
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar update
	respondError(w, http.StatusNotImplemented, "Not implemented")
}

// Delete deleta um client (admin)
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar delete
	respondError(w, http.StatusNotImplemented, "Not implemented")
}

// RegisterUser registra um novo usuário para um client (público)
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	clientUser, err := h.clientUC.RegisterUser(r.Context(), slug, req)
	if err != nil {
		switch err {
		case sharedErrors.ErrClientNotFound:
			respondError(w, http.StatusNotFound, "URL de registro inválida")
		case sharedErrors.ErrClientInactive:
			respondError(w, http.StatusForbidden, "Client não está aceitando novos cadastros")
		case sharedErrors.ErrClientUserLimitReached:
			respondError(w, http.StatusForbidden, "Limite de usuários atingido")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	resp := dto.RegisterUserResponse{
		UserID:   clientUser.UserID,
		Email:    clientUser.Email,
		ClientID: clientUser.ClientID,
		Message:  "Usuário registrado com sucesso",
	}
	respondJSON(w, http.StatusCreated, response.NewSuccessResponse(resp))
}

// GetMyClient retorna o client do usuário autenticado
func (h *Handler) GetMyClient(w http.ResponseWriter, r *http.Request) {
	clientID, ok := sharedContext.GetClientID(r.Context())
	if !ok {
		respondError(w, http.StatusForbidden, "No client association")
		return
	}

	client, err := h.clientUC.GetClient(r.Context(), clientID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get client")
		return
	}

	resp := dto.ToClientResponse(client, h.env.AppBaseURL)
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(resp))
}

// GetMyStats retorna estatísticas do client do usuário autenticado
func (h *Handler) GetMyStats(w http.ResponseWriter, r *http.Request) {
	clientID, ok := sharedContext.GetClientID(r.Context())
	if !ok {
		respondError(w, http.StatusForbidden, "No client association")
		return
	}

	stats, err := h.clientUC.GetClientStats(r.Context(), clientID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	resp := dto.ToClientStatsResponse(stats, h.env.AppBaseURL)
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(resp))
}

// ListMyUsers lista usuários do client do usuário autenticado
func (h *Handler) ListMyUsers(w http.ResponseWriter, r *http.Request) {
	clientID, ok := sharedContext.GetClientID(r.Context())
	if !ok {
		respondError(w, http.StatusForbidden, "No client association")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.clientUC.ListClientUsers(r.Context(), clientID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	var resp []dto.ClientUserResponse
	for _, user := range users {
		resp = append(resp, *dto.ToClientUserResponse(user))
	}

	data := map[string]interface{}{
		"users": resp,
		"total": total,
		"page":  page,
	}
	respondJSON(w, http.StatusOK, response.NewSuccessResponse(data))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, response.ErrorResponse{Message: message})
}
