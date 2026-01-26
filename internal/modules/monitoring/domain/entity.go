package domain

import (
	"time"
)

// MonitoramentoStatus representa o status do processamento
type MonitoramentoStatus string

const (
	StatusProcessando MonitoramentoStatus = "processando"
	StatusConcluido   MonitoramentoStatus = "concluido"
	StatusErro        MonitoramentoStatus = "erro"
)

// IsValid verifica se o status é válido
func (s MonitoramentoStatus) IsValid() bool {
	switch s {
	case StatusProcessando, StatusConcluido, StatusErro:
		return true
	}
	return false
}

// Monitoramento representa um upload de CSV
type Monitoramento struct {
	ID          string
	DataUpload  time.Time
	NomeArquivo string
	Status      MonitoramentoStatus
	TotalLinhas int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMonitoramento cria um novo monitoramento com status processando
func NewMonitoramento(id, nomeArquivo string) *Monitoramento {
	now := time.Now()
	return &Monitoramento{
		ID:          id,
		DataUpload:  now,
		NomeArquivo: nomeArquivo,
		Status:      StatusProcessando,
		TotalLinhas: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsCompleted marca o monitoramento como concluído
func (m *Monitoramento) MarkAsCompleted(totalLinhas int) {
	m.Status = StatusConcluido
	m.TotalLinhas = totalLinhas
	m.UpdatedAt = time.Now()
}

// MarkAsError marca o monitoramento como erro
func (m *Monitoramento) MarkAsError() {
	m.Status = StatusErro
	m.UpdatedAt = time.Now()
}
