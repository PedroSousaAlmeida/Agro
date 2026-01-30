package repository

import (
	"context"
	"sort"
	"sync"

	"agro-monitoring/internal/modules/clients/domain"
)

type InMemoryClientUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.ClientUser
}

func NewInMemoryClientUserRepository() domain.ClientUserRepository {
	return &InMemoryClientUserRepository{
		users: make(map[string]*domain.ClientUser),
	}
}

func (r *InMemoryClientUserRepository) Create(ctx context.Context, cu *domain.ClientUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[cu.ID] = cu
	return nil
}

func (r *InMemoryClientUserRepository) GetByClientAndUserID(ctx context.Context, clientID, userID string) (*domain.ClientUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, cu := range r.users {
		if cu.ClientID == clientID && cu.UserID == userID {
			return cu, nil
		}
	}

	return nil, nil
}

func (r *InMemoryClientUserRepository) CountActiveByClient(ctx context.Context, clientID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, cu := range r.users {
		if cu.ClientID == clientID && cu.Active {
			count++
		}
	}

	return count, nil
}

func (r *InMemoryClientUserRepository) ListByClient(ctx context.Context, clientID string, limit, offset int) ([]*domain.ClientUser, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter by client
	filtered := make([]*domain.ClientUser, 0)
	for _, cu := range r.users {
		if cu.ClientID == clientID {
			filtered = append(filtered, cu)
		}
	}

	// Sort by created_at DESC
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)

	// Apply pagination
	start := offset
	if start > total {
		return []*domain.ClientUser{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *InMemoryClientUserRepository) Deactivate(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cu, ok := r.users[id]; ok {
		cu.Active = false
	}

	return nil
}

func (r *InMemoryClientUserRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[string]*domain.ClientUser)
}
