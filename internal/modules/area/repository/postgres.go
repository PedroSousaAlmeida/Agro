package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"agro-monitoring/internal/modules/area/domain"
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

func (r *PostgresRepository) CreateBatch(ctx context.Context, areas []*domain.AreaMonitoramento) error {
	if len(areas) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO areas_monitoramento (
			id, monitoramento_id, setor, setor2, cod_fazenda, desc_fazenda,
			quadra, corte, area_total, desc_textura_solo, corte_atual,
			reforma, mes_colheita, restricao, pragas_data, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range areas {
		pragasJSON, err := a.PragasData.Value()
		if err != nil {
			return fmt.Errorf("erro ao serializar pragas: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			a.ID,
			a.MonitoramentoID,
			a.Setor,
			a.Setor2,
			a.CodFazenda,
			a.DescFazenda,
			a.Quadra,
			a.Corte,
			a.AreaTotal,
			a.DescTexturaSolo,
			a.CorteAtual,
			a.Reforma,
			a.MesColheita,
			a.Restricao,
			pragasJSON,
			a.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.AreaMonitoramento, error) {
	query := `
		SELECT id, monitoramento_id, setor, setor2, cod_fazenda, desc_fazenda,
			quadra, corte, area_total, desc_textura_solo, corte_atual,
			reforma, mes_colheita, restricao, pragas_data, created_at
		FROM areas_monitoramento
		WHERE id = $1
	`

	a := &domain.AreaMonitoramento{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.MonitoramentoID,
		&a.Setor,
		&a.Setor2,
		&a.CodFazenda,
		&a.DescFazenda,
		&a.Quadra,
		&a.Corte,
		&a.AreaTotal,
		&a.DescTexturaSolo,
		&a.CorteAtual,
		&a.Reforma,
		&a.MesColheita,
		&a.Restricao,
		&a.PragasData,
		&a.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedErrors.ErrAreaMonitoramentoNotFound
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (r *PostgresRepository) GetByMonitoramentoID(ctx context.Context, monitoramentoID string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM areas_monitoramento WHERE monitoramento_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, monitoramentoID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, monitoramento_id, setor, setor2, cod_fazenda, desc_fazenda,
			quadra, corte, area_total, desc_textura_solo, corte_atual,
			reforma, mes_colheita, restricao, pragas_data, created_at
		FROM areas_monitoramento
		WHERE monitoramento_id = $1
		ORDER BY created_at
		LIMIT $2 OFFSET $3
	`

	return r.queryAreas(ctx, query, total, monitoramentoID, limit, offset)
}

func (r *PostgresRepository) SearchByFazenda(ctx context.Context, codFazenda string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	search := "%" + strings.ToLower(codFazenda) + "%"

	var total int
	countQuery := `SELECT COUNT(*) FROM areas_monitoramento WHERE LOWER(cod_fazenda) LIKE $1`
	if err := r.db.QueryRowContext(ctx, countQuery, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, monitoramento_id, setor, setor2, cod_fazenda, desc_fazenda,
			quadra, corte, area_total, desc_textura_solo, corte_atual,
			reforma, mes_colheita, restricao, pragas_data, created_at
		FROM areas_monitoramento
		WHERE LOWER(cod_fazenda) LIKE $1
		ORDER BY cod_fazenda
		LIMIT $2 OFFSET $3
	`

	return r.queryAreas(ctx, query, total, search, limit, offset)
}

func (r *PostgresRepository) SearchByPraga(ctx context.Context, nomePraga string, limit, offset int) ([]*domain.AreaMonitoramento, int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) FROM areas_monitoramento
		WHERE pragas_data->'pragas' ? $1
		AND pragas_data->'pragas'->$1->>'presente' = 'true'
	`
	if err := r.db.QueryRowContext(ctx, countQuery, nomePraga).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, monitoramento_id, setor, setor2, cod_fazenda, desc_fazenda,
			quadra, corte, area_total, desc_textura_solo, corte_atual,
			reforma, mes_colheita, restricao, pragas_data, created_at
		FROM areas_monitoramento
		WHERE pragas_data->'pragas' ? $1
		AND pragas_data->'pragas'->$1->>'presente' = 'true'
		ORDER BY created_at
		LIMIT $2 OFFSET $3
	`

	return r.queryAreas(ctx, query, total, nomePraga, limit, offset)
}

func (r *PostgresRepository) UpdatePragasData(ctx context.Context, id string, pragasData domain.PragasData) error {
	pragasJSON, err := pragasData.Value()
	if err != nil {
		return err
	}

	query := `UPDATE areas_monitoramento SET pragas_data = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, pragasJSON, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sharedErrors.ErrAreaMonitoramentoNotFound
	}

	return nil
}

func (r *PostgresRepository) queryAreas(ctx context.Context, query string, total int, args ...interface{}) ([]*domain.AreaMonitoramento, int, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*domain.AreaMonitoramento
	for rows.Next() {
		a := &domain.AreaMonitoramento{}
		if err := rows.Scan(
			&a.ID,
			&a.MonitoramentoID,
			&a.Setor,
			&a.Setor2,
			&a.CodFazenda,
			&a.DescFazenda,
			&a.Quadra,
			&a.Corte,
			&a.AreaTotal,
			&a.DescTexturaSolo,
			&a.CorteAtual,
			&a.Reforma,
			&a.MesColheita,
			&a.Restricao,
			&a.PragasData,
			&a.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, a)
	}

	return result, total, rows.Err()
}
