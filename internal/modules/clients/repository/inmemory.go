package repository

import (
	"context"
	"sort"
	"sync"

	"agro-monitoring/internal/modules/clients/domain"
)

type InMemoryRepository struct {
	mu      sync.RWMutex
	clients map[string]*domain.Client
	slugs   map[string]string // slug -> id
}

func NewInMemoryRepository() domain.ClientRepository {
	return &InMemoryRepository{
		clients: make(map[string]*domain.Client),
		slugs:   make(map[string]string),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, client *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[client.ID] = client
	r.slugs[client.Slug] = client.ID
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.clients[id], nil
}

func (r *InMemoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.slugs[slug]
	if !ok {
		return nil, nil
	}
	return r.clients[id], nil
}

func (r *InMemoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	all := make([]*domain.Client, 0, len(r.clients))
	for _, client := range r.clients {
		all = append(all, client)
	}

	// Sort by created_at DESC
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := len(all)

	// Apply pagination
	start := offset
	if start > total {
		return []*domain.Client{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return all[start:end], total, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, client *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Remove old slug mapping if slug changed
	if existing, ok := r.clients[client.ID]; ok {
		if existing.Slug != client.Slug {
			delete(r.slugs, existing.Slug)
			r.slugs[client.Slug] = client.ID
		}
	}

	r.clients[client.ID] = client
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if client, ok := r.clients[id]; ok {
		delete(r.slugs, client.Slug)
		delete(r.clients, id)
	}

	return nil
}

func (r *InMemoryRepository) GetStats(ctx context.Context, clientID string) (*domain.ClientStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client := r.clients[clientID]
	if client == nil {
		return nil, nil
	}

	// Mock stats (real implementation would query database)
	return &domain.ClientStats{
		Client:              *client,
		CurrentUsers:        0, // Would be computed from client_users
		AvailableSlots:      client.MaxUsers,
		TotalMonitoramentos: 0,
		TotalAreas:          0,
	}, nil
}

func (r *InMemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients = make(map[string]*domain.Client)
	r.slugs = make(map[string]string)
}
