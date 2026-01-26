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
)
