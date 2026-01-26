package repository

import (
	"context"
	"sync"

	"agro-monitoring/internal/modules/monitoring/domain"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

// InMemoryRepository implementação em memória para testes
type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Monitoramento
}

// NewInMemoryRepository cria um novo repository em memória
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*domain.Monitoramento),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, m *domain.Monitoramento) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[m.ID] = m
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*domain.Monitoramento, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.items[id]
	if !ok {
		return nil, sharedErrors.ErrMonitoramentoNotFound
	}
	return m, nil
}

func (r *InMemoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.Monitoramento, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]*domain.Monitoramento, 0, len(r.items))
	for _, m := range r.items {
		all = append(all, m)
	}

	total := len(all)

	if offset >= len(all) {
		return []*domain.Monitoramento{}, total, nil
	}
	all = all[offset:]

	if limit > 0 && limit < len(all) {
		all = all[:limit]
	}

	return all, total, nil
}

func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id string, status domain.MonitoramentoStatus, totalLinhas int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.items[id]
	if !ok {
		return sharedErrors.ErrMonitoramentoNotFound
	}

	m.Status = status
	m.TotalLinhas = totalLinhas
	return nil
}

// Clear limpa todos os dados (útil para testes)
func (r *InMemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.Monitoramento)
}
