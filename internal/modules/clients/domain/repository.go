package domain

import "context"

// ClientRepository define operações de persistência para clients
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id string) (*Client, error)
	GetBySlug(ctx context.Context, slug string) (*Client, error)
	List(ctx context.Context, limit, offset int) ([]*Client, int, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, clientID string) (*ClientStats, error)
}

// ClientUserRepository define operações de persistência para client_users
type ClientUserRepository interface {
	Create(ctx context.Context, cu *ClientUser) error
	GetByClientAndUserID(ctx context.Context, clientID, userID string) (*ClientUser, error)
	CountActiveByClient(ctx context.Context, clientID string) (int, error)
	ListByClient(ctx context.Context, clientID string, limit, offset int) ([]*ClientUser, int, error)
	Deactivate(ctx context.Context, id string) error
}
