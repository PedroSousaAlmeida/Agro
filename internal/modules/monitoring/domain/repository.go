package domain

import "context"

// MonitoramentoRepository define as operações de persistência
type MonitoramentoRepository interface {
	Create(ctx context.Context, m *Monitoramento) error
	GetByID(ctx context.Context, id string) (*Monitoramento, error)
	List(ctx context.Context, limit, offset int) ([]*Monitoramento, int, error)
	UpdateStatus(ctx context.Context, id string, status MonitoramentoStatus, totalLinhas int) error
}
