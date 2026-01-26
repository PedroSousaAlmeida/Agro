package bootstrap

import (
	"agro-monitoring/internal/config"
)

// Re-exporta tipos do config para manter compatibilidade
type Env = config.Env

// NewEnv cria um novo Env
func NewEnv() *Env {
	return config.NewEnv()
}
