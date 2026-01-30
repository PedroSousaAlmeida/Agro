package domain

import "time"

// Client representa uma usina/empresa que contrata o sistema
type Client struct {
	ID              string
	Name            string
	Slug            string
	MaxUsers        int
	Active          bool
	Metadata        map[string]interface{}
	KeycloakGroupID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ClientUser representa a relação entre um client e um usuário do Keycloak
type ClientUser struct {
	ID        string
	ClientID  string
	UserID    string // Sub claim do Keycloak
	Email     string
	Role      string // "user" ou "admin"
	Active    bool
	CreatedAt time.Time
}

// ClientStats contém estatísticas de um client
type ClientStats struct {
	Client
	CurrentUsers        int
	AvailableSlots      int
	TotalMonitoramentos int
	TotalAreas          int
}

// NewClient cria um novo client
func NewClient(id, name, slug string, maxUsers int) *Client {
	now := time.Now()
	return &Client{
		ID:        id,
		Name:      name,
		Slug:      slug,
		MaxUsers:  maxUsers,
		Active:    true,
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewClientUser cria um novo ClientUser
func NewClientUser(id, clientID, userID, email, role string) *ClientUser {
	return &ClientUser{
		ID:        id,
		ClientID:  clientID,
		UserID:    userID,
		Email:     email,
		Role:      role,
		Active:    true,
		CreatedAt: time.Now(),
	}
}
