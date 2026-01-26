package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"agro-monitoring/internal/modules/jobs/domain"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

type PostgresJobRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) domain.JobRepository {
	return &PostgresJobRepository{db: db}
}

func (r *PostgresJobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `
		INSERT INTO jobs (id, type, status, payload, total_items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, job.ID, job.Type, job.Status, job.Payload, job.TotalItems, job.CreatedAt, job.UpdatedAt)
	return err
}

func (r *PostgresJobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	query := `
		SELECT 
			id, type, status, payload, result, 
			progress, total_items, processed_items, error_count, error_details,
			started_at, completed_at, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`
	job := &domain.Job{}
	var payload, result, errorDetails sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Type, &job.Status, &payload, &result,
		&job.Progress, &job.TotalItems, &job.ProcessedItems, &job.ErrorCount, &errorDetails,
		&job.StartedAt, &job.CompletedAt, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sharedErrors.ErrJobNotFound
		}
		return nil, err
	}

	if payload.Valid {
		job.Payload = json.RawMessage(payload.String)
	}
	if result.Valid {
		job.Result = json.RawMessage(result.String)
	}
	if errorDetails.Valid {
		job.ErrorDetails = json.RawMessage(errorDetails.String)
	}

	return job, nil
}

func (r *PostgresJobRepository) Update(ctx context.Context, job *domain.Job) error {
	query := `
		UPDATE jobs
		SET
			status = $2,
			payload = $3,
			result = $4,
			progress = $5,
			total_items = $6,
			processed_items = $7,
			error_count = $8,
			error_details = $9,
			started_at = $10,
			completed_at = $11,
			updated_at = $12
		WHERE id = $1
	`

	// Trata campos JSON nulos
	var payload, result, errorDetails interface{}
	if len(job.Payload) > 0 {
		payload = job.Payload
	}
	if len(job.Result) > 0 {
		result = job.Result
	}
	if len(job.ErrorDetails) > 0 {
		errorDetails = job.ErrorDetails
	}

	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Status, payload, result,
		job.Progress, job.TotalItems, job.ProcessedItems, job.ErrorCount, errorDetails,
		job.StartedAt, job.CompletedAt, job.UpdatedAt,
	)
	return err
}

func (r *PostgresJobRepository) UpdateProgress(ctx context.Context, id string, processed, errorCount int) error {
	query := `
        UPDATE jobs
        SET 
            processed_items = $2,
            error_count = $3,
            progress = CASE
                WHEN total_items > 0 THEN ($2 * 100) / total_items
                ELSE 0
            END,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, id, processed, errorCount)
	return err
}

func (r *PostgresJobRepository) List(ctx context.Context, status *domain.JobStatus, limit, offset int) ([]*domain.Job, int, error) {
	// Implementação futura, se necessário
	return nil, 0, errors.New("not implemented")
}
