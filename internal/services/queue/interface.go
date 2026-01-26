package queue

import (
	"context"

	"agro-monitoring/internal/modules/jobs/domain"
)

const (
	QueueDefault = "default"
)

// EnqueueOptions opções para enfileirar
type EnqueueOptions struct {
	QueueName string
	Delay     int
}

// Job a ser processado
type Job struct {
	ID        string
	Queue     string
	Payload   []byte
	JobEntity *domain.Job
}

// Service define a interface do serviço de fila
type Service interface {
	Enqueue(ctx context.Context, job *Job, opts *EnqueueOptions) error
	Dequeue(ctx context.Context, queueName string) (*Job, error)
	Close() error
}
