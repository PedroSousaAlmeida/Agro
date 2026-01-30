package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"agro-monitoring/internal/modules/clients/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) domain.ClientRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, client *domain.Client) error {
	metadataJSON, err := json.Marshal(client.Metadata)
	if err != nil {
		return fmt.Errorf("erro ao serializar metadata: %w", err)
	}

	query := `
		INSERT INTO clients (id, name, slug, max_users, active, metadata, keycloak_group_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.ExecContext(ctx, query,
		client.ID,
		client.Name,
		client.Slug,
		client.MaxUsers,
		client.Active,
		metadataJSON,
		client.KeycloakGroupID,
		client.CreatedAt,
		client.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar client: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	query := `
		SELECT id, name, slug, max_users, active, metadata, keycloak_group_id, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	var client domain.Client
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Slug,
		&client.MaxUsers,
		&client.Active,
		&metadataJSON,
		&client.KeycloakGroupID,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar client: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &client.Metadata); err != nil {
		return nil, fmt.Errorf("erro ao deserializar metadata: %w", err)
	}

	return &client, nil
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Client, error) {
	query := `
		SELECT id, name, slug, max_users, active, metadata, keycloak_group_id, created_at, updated_at
		FROM clients
		WHERE slug = $1
	`

	var client domain.Client
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&client.ID,
		&client.Name,
		&client.Slug,
		&client.MaxUsers,
		&client.Active,
		&metadataJSON,
		&client.KeycloakGroupID,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar client por slug: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &client.Metadata); err != nil {
		return nil, fmt.Errorf("erro ao deserializar metadata: %w", err)
	}

	return &client, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, int, error) {
	// Contar total
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM clients").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar clients: %w", err)
	}

	// Buscar paginado
	query := `
		SELECT id, name, slug, max_users, active, metadata, keycloak_group_id, created_at, updated_at
		FROM clients
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar clients: %w", err)
	}
	defer rows.Close()

	var clients []*domain.Client
	for rows.Next() {
		var client domain.Client
		var metadataJSON []byte

		err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Slug,
			&client.MaxUsers,
			&client.Active,
			&metadataJSON,
			&client.KeycloakGroupID,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao escanear client: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &client.Metadata); err != nil {
			return nil, 0, fmt.Errorf("erro ao deserializar metadata: %w", err)
		}

		clients = append(clients, &client)
	}

	return clients, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, client *domain.Client) error {
	metadataJSON, err := json.Marshal(client.Metadata)
	if err != nil {
		return fmt.Errorf("erro ao serializar metadata: %w", err)
	}

	query := `
		UPDATE clients
		SET name = $1, slug = $2, max_users = $3, active = $4, metadata = $5,
		    keycloak_group_id = $6, updated_at = $7
		WHERE id = $8
	`

	_, err = r.db.ExecContext(ctx, query,
		client.Name,
		client.Slug,
		client.MaxUsers,
		client.Active,
		metadataJSON,
		client.KeycloakGroupID,
		client.UpdatedAt,
		client.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar client: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM clients WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar client: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetStats(ctx context.Context, clientID string) (*domain.ClientStats, error) {
	query := `
		SELECT id, name, slug, max_users, current_users, available_slots,
		       total_monitoramentos, total_areas, active, created_at
		FROM client_stats
		WHERE id = $1
	`

	var stats domain.ClientStats
	var metadata = make(map[string]interface{})

	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&stats.ID,
		&stats.Name,
		&stats.Slug,
		&stats.MaxUsers,
		&stats.CurrentUsers,
		&stats.AvailableSlots,
		&stats.TotalMonitoramentos,
		&stats.TotalAreas,
		&stats.Active,
		&stats.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %w", err)
	}

	stats.Metadata = metadata

	return &stats, nil
}
