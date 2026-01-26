package usecase

import (
	"context"
	"encoding/json"
	"log"

	areaDomain "agro-monitoring/internal/modules/area/domain"
	"agro-monitoring/internal/modules/jobs/domain"
	"agro-monitoring/internal/services/queue"
)

const (
	QueueBulkAplicacoes = "jobs:bulk_aplicacoes"
)

type jobUseCase struct {
	uuidGen  func() string
	jobRepo  domain.JobRepository
	areaRepo areaDomain.AreaMonitoramentoRepository
	queue    queue.Service
}

// NewJobUseCase cria um novo usecase de jobs
func NewJobUseCase(cfg Config) JobUseCase {
	return &jobUseCase{
		uuidGen:  cfg.UUIDGenerator,
		jobRepo:  cfg.JobRepo,
		areaRepo: cfg.AreaRepo,
		queue:    cfg.Queue,
	}
}

// CreateBulkAplicacoesJob cria um job para processar aplicações em massa
func (uc *jobUseCase) CreateBulkAplicacoesJob(ctx context.Context, payload domain.BulkAplicacoesPayload) (*domain.Job, error) {
	// Cria o job
	job, err := domain.NewJob(uc.uuidGen(), domain.JobTypeBulkAplicacoes, payload)
	if err != nil {
		return nil, err
	}

	job.TotalItems = len(payload.Aplicacoes)

	// Salva no banco
	if err := uc.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	// Enfileira para processamento
	queueJob := &queue.Job{
		ID:        job.ID,
		Queue:     QueueBulkAplicacoes,
		JobEntity: job,
	}

	if err := uc.queue.Enqueue(ctx, queueJob, &queue.EnqueueOptions{QueueName: QueueBulkAplicacoes}); err != nil {
		// Atualiza status para falha se não conseguir enfileirar
		job.Fail([]domain.JobError{{Message: "Falha ao enfileirar job: " + err.Error()}})
		uc.jobRepo.Update(ctx, job)
		return nil, err
	}

	return job, nil
}

// GetJobStatus retorna o status de um job
func (uc *jobUseCase) GetJobStatus(ctx context.Context, jobID string) (*domain.Job, error) {
	return uc.jobRepo.GetByID(ctx, jobID)
}

// RegisterAndProcessJobs inicia o worker para processar jobs (chamado pelo cmd/worker)
func (uc *jobUseCase) RegisterAndProcessJobs(ctx context.Context) {
	log.Println("Worker iniciado, aguardando jobs na fila:", QueueBulkAplicacoes)

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker encerrado")
			return
		default:
			// Bloqueia esperando job na fila
			queueJob, err := uc.queue.Dequeue(ctx, QueueBulkAplicacoes)
			if err != nil {
				if ctx.Err() != nil {
					return // Context cancelado
				}
				log.Printf("Erro ao buscar job da fila: %v", err)
				continue
			}

			if queueJob == nil {
				log.Println("Nothing in queue")
				continue
			}

			if queueJob.JobEntity != nil {
				uc.processJob(ctx, queueJob.JobEntity)
			}
		}
	}
}

// processJob processa um job de aplicações
func (uc *jobUseCase) processJob(ctx context.Context, job *domain.Job) {
	log.Printf("Processando job %s tipo %s", job.ID, job.Type)

	// Busca job atualizado do banco
	job, err := uc.jobRepo.GetByID(ctx, job.ID)
	if err != nil {
		log.Printf("Erro ao buscar job %s: %v", job.ID, err)
		return
	}

	// Parse do payload
	var payload domain.BulkAplicacoesPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		log.Printf("Erro ao parsear payload do job %s: %v", job.ID, err)
		job.Fail([]domain.JobError{{Message: "Payload inválido: " + err.Error()}})
		uc.jobRepo.Update(ctx, job)
		return
	}

	// Marca como processando
	job.Start(len(payload.Aplicacoes))
	if err := uc.jobRepo.Update(ctx, job); err != nil {
		log.Printf("Erro ao atualizar status do job %s para processing: %v", job.ID, err)
	}

	var errors []domain.JobError
	processed := 0

	// Processa cada aplicação
	for i, item := range payload.Aplicacoes {
		err := uc.processAplicacao(ctx, item)
		if err != nil {
			errors = append(errors, domain.JobError{
				Line:    i + 1,
				ItemID:  item.AreaID,
				Message: err.Error(),
			})
			job.AddError()
		} else {
			processed++
		}

		// Atualiza progresso a cada 100 itens ou no final
		if (i+1)%100 == 0 || i == len(payload.Aplicacoes)-1 {
			job.UpdateProgress(processed + len(errors))
			uc.jobRepo.UpdateProgress(ctx, job.ID, job.ProcessedItems, job.ErrorCount)
		}
	}

	// Finaliza job
	if len(errors) > 0 && processed == 0 {
		job.Fail(errors)
	} else {
		result := domain.BulkAplicacoesResult{
			Processed: processed,
			Errors:    len(errors),
		}
		job.Complete(result)
	}

	if err := uc.jobRepo.Update(ctx, job); err != nil {
		log.Printf("Erro ao atualizar job %s para status final: %v", job.ID, err)
	}
	log.Printf("Job %s concluído: %d processados, %d erros", job.ID, processed, len(errors))
}

// processAplicacao processa uma aplicação individual
func (uc *jobUseCase) processAplicacao(ctx context.Context, item domain.AplicacaoItem) error {
	// Busca a área
	area, err := uc.areaRepo.GetByID(ctx, item.AreaID)
	if err != nil {
		return err
	}

	// Adiciona/atualiza aplicação na praga (upsert por posição)
	if err := area.PragasData.AddAplicacao(item.Praga, item.Posicao, item.Herbicida, item.Dose); err != nil {
		return err
	}

	// Salva atualização
	return uc.areaRepo.UpdatePragasData(ctx, item.AreaID, area.PragasData)
}
