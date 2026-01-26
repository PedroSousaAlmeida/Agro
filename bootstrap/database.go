package bootstrap

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// NewDatabase cria conexão com PostgreSQL
func NewDatabase(env *Env) (*sql.DB, error) {
	db, err := sql.Open("postgres", env.DSN())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	// Configurações do pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Testa conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar banco: %w", err)
	}

	return db, nil
}
