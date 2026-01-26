package repository

import (
	"context"
	"database/sql"
	"time"

	"agro-monitoring/internal/modules/monitoring/domain"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

// PostgresRepository implementação PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository cria um novo repository PostgreSQL
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, m *domain.Monitoramento) error {
	query := `
		INSERT INTO monitoramentos (id, data_upload, nome_arquivo, status, total_linhas, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.DataUpload,
		m.NomeArquivo,
		m.Status,
		m.TotalLinhas,
		m.CreatedAt,
		m.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Monitoramento, error) {
	query := `
		SELECT id, data_upload, nome_arquivo, status, total_linhas, created_at, updated_at
		FROM monitoramentos
		WHERE id = $1
	`

	m := &domain.Monitoramento{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.DataUpload,
		&m.NomeArquivo,
		&m.Status,
		&m.TotalLinhas,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedErrors.ErrMonitoramentoNotFound
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.Monitoramento, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM monitoramentos`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, data_upload, nome_arquivo, status, total_linhas, created_at, updated_at
		FROM monitoramentos
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*domain.Monitoramento
	for rows.Next() {
		m := &domain.Monitoramento{}
		if err := rows.Scan(
			&m.ID,
			&m.DataUpload,
			&m.NomeArquivo,
			&m.Status,
			&m.TotalLinhas,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, m)
	}

	return result, total, rows.Err()
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.MonitoramentoStatus, totalLinhas int) error {
	query := `
		UPDATE monitoramentos
		SET status = $1, total_linhas = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, status, totalLinhas, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sharedErrors.ErrMonitoramentoNotFound
	}

	return nil
}
