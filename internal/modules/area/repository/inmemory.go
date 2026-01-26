package repository

import (
	"context"
	"strings"
	"sync"

	"agro-monitoring/internal/modules/area/domain"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

// InMemoryRepository implementação em memória para testes
type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.AreaMonitoramento
}

// NewInMemoryRepository cria um novo repository em memória
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*domain.AreaMonitoramento),
	}
}

func (r *InMemoryRepository) CreateBatch(ctx context.Context, areas []*domain.AreaMonitoramento) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, a := range areas {
		r.items[a.ID] = a
	}
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*domain.AreaMonitoramento, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.items[id]
	if !ok {
		return nil, sharedErrors.ErrAreaMonitoramentoNotFound
	}
	return a, nil
}

func (r *InMemoryRepository) GetByMonitoramentoID(ctx context.Context, monitoramentoID string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.AreaMonitoramento
	for _, a := range r.items {
		if a.MonitoramentoID == monitoramentoID {
			result = append(result, a)
		}
	}

	total := len(result)

	if offset >= len(result) {
		return []*domain.AreaMonitoramento{}, total, nil
	}
	result = result[offset:]

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, total, nil
}

func (r *InMemoryRepository) SearchByFazenda(ctx context.Context, codFazenda string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.AreaMonitoramento
	for _, a := range r.items {
		if strings.Contains(strings.ToLower(a.CodFazenda), strings.ToLower(codFazenda)) {
			result = append(result, a)
		}
	}

	total := len(result)

	if offset >= len(result) {
		return []*domain.AreaMonitoramento{}, total, nil
	}
	result = result[offset:]

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, total, nil
}

func (r *InMemoryRepository) SearchByPraga(ctx context.Context, nomePraga string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.AreaMonitoramento
	for _, a := range r.items {
		if a.PragasData.HasPraga(nomePraga) {
			result = append(result, a)
		}
	}

	total := len(result)

	if offset >= len(result) {
		return []*domain.AreaMonitoramento{}, total, nil
	}
	result = result[offset:]

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, total, nil
}

func (r *InMemoryRepository) UpdatePragasData(ctx context.Context, id string, pragasData domain.PragasData) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	a, ok := r.items[id]
	if !ok {
		return sharedErrors.ErrAreaMonitoramentoNotFound
	}

	a.PragasData = pragasData
	return nil
}

// Clear limpa todos os dados (útil para testes)
func (r *InMemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*domain.AreaMonitoramento)
}
