package dto

import (
	"time"

	"agro-monitoring/internal/modules/clients/domain"
)

// CreateClientRequest representa a requisição para criar um client
type CreateClientRequest struct {
	Name     string                 `json:"name"`
	Slug     string                 `json:"slug"`
	MaxUsers int                    `json:"max_users"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ClientResponse representa a resposta com dados de um client
type ClientResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	MaxUsers    int                    `json:"max_users"`
	RegisterURL string                 `json:"register_url"`
	Active      bool                   `json:"active"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

// RegisterUserRequest representa a requisição para registrar um usuário
type RegisterUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// RegisterUserResponse representa a resposta após registrar usuário
type RegisterUserResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	ClientID string `json:"client_id"`
	Message  string `json:"message"`
}

// ClientStatsResponse representa estatísticas de um client
type ClientStatsResponse struct {
	ClientResponse
	CurrentUsers        int `json:"current_users"`
	AvailableSlots      int `json:"available_slots"`
	TotalMonitoramentos int `json:"total_monitoramentos"`
	TotalAreas          int `json:"total_areas"`
}

// ClientUserResponse representa um usuário de um client
type ClientUserResponse struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToClientResponse converte domain.Client para ClientResponse
func ToClientResponse(client *domain.Client, baseURL string) *ClientResponse {
	registerURL := baseURL + "/v1/register/" + client.Slug
	return &ClientResponse{
		ID:          client.ID,
		Name:        client.Name,
		Slug:        client.Slug,
		MaxUsers:    client.MaxUsers,
		RegisterURL: registerURL,
		Active:      client.Active,
		Metadata:    client.Metadata,
		CreatedAt:   client.CreatedAt,
	}
}

// ToClientStatsResponse converte domain.ClientStats para ClientStatsResponse
func ToClientStatsResponse(stats *domain.ClientStats, baseURL string) *ClientStatsResponse {
	return &ClientStatsResponse{
		ClientResponse:      *ToClientResponse(&stats.Client, baseURL),
		CurrentUsers:        stats.CurrentUsers,
		AvailableSlots:      stats.AvailableSlots,
		TotalMonitoramentos: stats.TotalMonitoramentos,
		TotalAreas:          stats.TotalAreas,
	}
}

// ToClientUserResponse converte domain.ClientUser para ClientUserResponse
func ToClientUserResponse(cu *domain.ClientUser) *ClientUserResponse {
	return &ClientUserResponse{
		ID:        cu.ID,
		ClientID:  cu.ClientID,
		UserID:    cu.UserID,
		Email:     cu.Email,
		Role:      cu.Role,
		Active:    cu.Active,
		CreatedAt: cu.CreatedAt,
	}
}
