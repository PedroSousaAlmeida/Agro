package domain

import "context"

// JobRepository interface de persistência
type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	UpdateProgress(ctx context.Context, id string, processed, errorCount int) error
	List(ctx context.Context, status *JobStatus, limit, offset int) ([]*Job, int, error)
}
