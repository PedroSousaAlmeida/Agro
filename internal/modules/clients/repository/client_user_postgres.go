package repository

import (
	"context"
	"database/sql"
	"fmt"

	"agro-monitoring/internal/modules/clients/domain"
)

type ClientUserPostgresRepository struct {
	db *sql.DB
}

func NewClientUserPostgresRepository(db *sql.DB) domain.ClientUserRepository {
	return &ClientUserPostgresRepository{db: db}
}

func (r *ClientUserPostgresRepository) Create(ctx context.Context, cu *domain.ClientUser) error {
	query := `
		INSERT INTO client_users (id, client_id, user_id, email, role, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		cu.ID,
		cu.ClientID,
		cu.UserID,
		cu.Email,
		cu.Role,
		cu.Active,
		cu.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar client_user: %w", err)
	}

	return nil
}

func (r *ClientUserPostgresRepository) GetByClientAndUserID(ctx context.Context, clientID, userID string) (*domain.ClientUser, error) {
	query := `
		SELECT id, client_id, user_id, email, role, active, created_at
		FROM client_users
		WHERE client_id = $1 AND user_id = $2
	`

	var cu domain.ClientUser
	err := r.db.QueryRowContext(ctx, query, clientID, userID).Scan(
		&cu.ID,
		&cu.ClientID,
		&cu.UserID,
		&cu.Email,
		&cu.Role,
		&cu.Active,
		&cu.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar client_user: %w", err)
	}

	return &cu, nil
}

func (r *ClientUserPostgresRepository) CountActiveByClient(ctx context.Context, clientID string) (int, error) {
	query := "SELECT COUNT(*) FROM client_users WHERE client_id = $1 AND active = true"

	var count int
	err := r.db.QueryRowContext(ctx, query, clientID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar usuários: %w", err)
	}

	return count, nil
}

func (r *ClientUserPostgresRepository) ListByClient(ctx context.Context, clientID string, limit, offset int) ([]*domain.ClientUser, int, error) {
	// Contar total
	var total int
	countQuery := "SELECT COUNT(*) FROM client_users WHERE client_id = $1"
	err := r.db.QueryRowContext(ctx, countQuery, clientID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar client_users: %w", err)
	}

	// Buscar paginado
	query := `
		SELECT id, client_id, user_id, email, role, active, created_at
		FROM client_users
		WHERE client_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar client_users: %w", err)
	}
	defer rows.Close()

	var users []*domain.ClientUser
	for rows.Next() {
		var cu domain.ClientUser
		err := rows.Scan(
			&cu.ID,
			&cu.ClientID,
			&cu.UserID,
			&cu.Email,
			&cu.Role,
			&cu.Active,
			&cu.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao escanear client_user: %w", err)
		}
		users = append(users, &cu)
	}

	return users, total, nil
}

func (r *ClientUserPostgresRepository) Deactivate(ctx context.Context, id string) error {
	query := "UPDATE client_users SET active = false WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao desativar client_user: %w", err)
	}
	return nil
}
