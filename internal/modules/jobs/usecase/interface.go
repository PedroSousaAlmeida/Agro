package usecase

import (
	"context"

	areaDomain "agro-monitoring/internal/modules/area/domain"
	"agro-monitoring/internal/modules/jobs/domain"
	queue "agro-monitoring/internal/services/queue"
)

// JobUseCase define a interface para os casos de uso de jobs
type JobUseCase interface {
	CreateBulkAplicacoesJob(ctx context.Context, payload domain.BulkAplicacoesPayload) (*domain.Job, error)
	GetJobStatus(ctx context.Context, jobID string) (*domain.Job, error)
	RegisterAndProcessJobs(ctx context.Context)
}

// Config contém as dependências para o usecase
type Config struct {
	UUIDGenerator func() string
	JobRepo       domain.JobRepository
	AreaRepo      areaDomain.AreaMonitoramentoRepository
	Queue         queue.Service
}
