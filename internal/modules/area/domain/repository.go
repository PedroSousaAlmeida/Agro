package domain

import "context"

// AreaMonitoramentoRepository define as operações de persistência
type AreaMonitoramentoRepository interface {
	CreateBatch(ctx context.Context, areas []*AreaMonitoramento) error
	GetByID(ctx context.Context, id string) (*AreaMonitoramento, error)
	GetByMonitoramentoID(ctx context.Context, monitoramentoID string, limit, offset int) ([]*AreaMonitoramento, int, error)
	SearchByFazenda(ctx context.Context, codFazenda string, limit, offset int) ([]*AreaMonitoramento, int, error)
	SearchByPraga(ctx context.Context, nomePraga string, limit, offset int) ([]*AreaMonitoramento, int, error)
	UpdatePragasData(ctx context.Context, id string, pragasData PragasData) error
}
