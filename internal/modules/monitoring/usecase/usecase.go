package usecase

import (
	"context"
	"io"

	areaDomain "agro-monitoring/internal/modules/area/domain"
	"agro-monitoring/internal/modules/monitoring/domain"
	"agro-monitoring/internal/services/csv"
)

// MonitoringUseCase interface para operações de monitoramento
type MonitoringUseCase interface {
	UploadAndProcessCSV(ctx context.Context, file io.Reader, filename string) (*domain.Monitoramento, error)
	GetMonitoramento(ctx context.Context, id string) (*domain.Monitoramento, error)
	ListMonitoramentos(ctx context.Context, page, pageSize int) ([]*domain.Monitoramento, int, error)
}

type monitoringUseCase struct {
	monitoramentoRepo domain.MonitoramentoRepository
	areaRepo          areaDomain.AreaMonitoramentoRepository
	csvParser         *csv.Parser
	uuidGenerator     func() string
}

// NewMonitoringUseCase cria um novo usecase de monitoramento
func NewMonitoringUseCase(
	monitoramentoRepo domain.MonitoramentoRepository,
	areaRepo areaDomain.AreaMonitoramentoRepository,
	csvParser *csv.Parser,
	uuidGenerator func() string,
) MonitoringUseCase {
	return &monitoringUseCase{
		monitoramentoRepo: monitoramentoRepo,
		areaRepo:          areaRepo,
		csvParser:         csvParser,
		uuidGenerator:     uuidGenerator,
	}
}

func (uc *monitoringUseCase) UploadAndProcessCSV(ctx context.Context, file io.Reader, filename string) (*domain.Monitoramento, error) {
	monitoramento := domain.NewMonitoramento(uc.uuidGenerator(), filename)

	if err := uc.monitoramentoRepo.Create(ctx, monitoramento); err != nil {
		return nil, err
	}

	result, err := uc.csvParser.Parse(file, monitoramento.ID)
	if err != nil {
		uc.monitoramentoRepo.UpdateStatus(ctx, monitoramento.ID, domain.StatusErro, 0)
		return nil, err
	}

	if len(result.Areas) > 0 {
		if err := uc.areaRepo.CreateBatch(ctx, result.Areas); err != nil {
			uc.monitoramentoRepo.UpdateStatus(ctx, monitoramento.ID, domain.StatusErro, 0)
			return nil, err
		}
	}

	if err := uc.monitoramentoRepo.UpdateStatus(ctx, monitoramento.ID, domain.StatusConcluido, result.TotalLinhas); err != nil {
		return nil, err
	}

	return uc.monitoramentoRepo.GetByID(ctx, monitoramento.ID)
}

func (uc *monitoringUseCase) GetMonitoramento(ctx context.Context, id string) (*domain.Monitoramento, error) {
	return uc.monitoramentoRepo.GetByID(ctx, id)
}

func (uc *monitoringUseCase) ListMonitoramentos(ctx context.Context, page, pageSize int) ([]*domain.Monitoramento, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return uc.monitoramentoRepo.List(ctx, pageSize, offset)
}
