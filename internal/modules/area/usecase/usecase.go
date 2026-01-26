package usecase

import (
	"context"

	"agro-monitoring/internal/modules/area/domain"
)

// AreaQueryUseCase interface para consultas de áreas
type AreaQueryUseCase interface {
	GetAreasByMonitoramento(ctx context.Context, monitoramentoID string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error)
	GetAreaByID(ctx context.Context, id string) (*domain.AreaMonitoramento, error)
	SearchByFazenda(ctx context.Context, codFazenda string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error)
	SearchByPraga(ctx context.Context, nomePraga string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error)
	AddAplicacaoHerbicida(ctx context.Context, areaID, praga string, posicao int, herbicida string, dose float64) error
}

type areaQueryUseCase struct {
	areaRepo domain.AreaMonitoramentoRepository
}

// NewAreaQueryUseCase cria um novo usecase de consulta de áreas
func NewAreaQueryUseCase(areaRepo domain.AreaMonitoramentoRepository) AreaQueryUseCase {
	return &areaQueryUseCase{
		areaRepo: areaRepo,
	}
}

func (uc *areaQueryUseCase) GetAreasByMonitoramento(ctx context.Context, monitoramentoID string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error) {
	offset, limit := uc.paginate(page, pageSize)
	return uc.areaRepo.GetByMonitoramentoID(ctx, monitoramentoID, limit, offset)
}

func (uc *areaQueryUseCase) GetAreaByID(ctx context.Context, id string) (*domain.AreaMonitoramento, error) {
	return uc.areaRepo.GetByID(ctx, id)
}

func (uc *areaQueryUseCase) SearchByFazenda(ctx context.Context, codFazenda string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error) {
	offset, limit := uc.paginate(page, pageSize)
	return uc.areaRepo.SearchByFazenda(ctx, codFazenda, limit, offset)
}

func (uc *areaQueryUseCase) SearchByPraga(ctx context.Context, nomePraga string, page, pageSize int) ([]*domain.AreaMonitoramento, int, error) {
	offset, limit := uc.paginate(page, pageSize)
	return uc.areaRepo.SearchByPraga(ctx, nomePraga, limit, offset)
}

func (uc *areaQueryUseCase) AddAplicacaoHerbicida(ctx context.Context, areaID, praga string, posicao int, herbicida string, dose float64) error {
	area, err := uc.areaRepo.GetByID(ctx, areaID)
	if err != nil {
		return err
	}

	if err := area.PragasData.AddAplicacao(praga, posicao, herbicida, dose); err != nil {
		return err
	}

	return uc.areaRepo.UpdatePragasData(ctx, areaID, area.PragasData)
}

func (uc *areaQueryUseCase) paginate(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return
}
