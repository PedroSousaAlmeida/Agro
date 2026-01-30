package errors

import "errors"

var (
	ErrMonitoramentoNotFound     = errors.New("monitoramento não encontrado")
	ErrAreaMonitoramentoNotFound = errors.New("área de monitoramento não encontrada")
	ErrJobNotFound               = errors.New("job não encontrado")
	ErrInvalidCSV                = errors.New("arquivo CSV inválido")
	ErrEmptyCSV                  = errors.New("arquivo CSV vazio")
	ErrInvalidStatus             = errors.New("status inválido")
	ErrPragaNotFound             = errors.New("praga não encontrada")
	ErrInvalidPragaData          = errors.New("dados de praga inválidos")

	// Clients
	ErrClientNotFound         = errors.New("client não encontrado")
	ErrClientInactive         = errors.New("client inativo")
	ErrClientUserLimitReached = errors.New("limite de usuários atingido")
	ErrDuplicateEmail         = errors.New("email já cadastrado para este client")
	ErrInvalidSlug            = errors.New("slug inválido")
)
