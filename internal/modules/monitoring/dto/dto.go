package dto

import (
	"time"

	"agro-monitoring/internal/modules/monitoring/domain"
)

// MonitoramentoResponse resposta de monitoramento
type MonitoramentoResponse struct {
	ID          string    `json:"id"`
	DataUpload  time.Time `json:"data_upload"`
	NomeArquivo string    `json:"nome_arquivo"`
	Status      string    `json:"status"`
	TotalLinhas int       `json:"total_linhas"`
	CreatedAt   time.Time `json:"created_at"`
}

// ListMonitoramentosResponse resposta paginada
type ListMonitoramentosResponse struct {
	Data       []MonitoramentoResponse `json:"data"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalCount int                     `json:"total_count"`
}

// ToMonitoramentoResponse converte domain para DTO
func ToMonitoramentoResponse(m *domain.Monitoramento) MonitoramentoResponse {
	return MonitoramentoResponse{
		ID:          m.ID,
		DataUpload:  m.DataUpload,
		NomeArquivo: m.NomeArquivo,
		Status:      string(m.Status),
		TotalLinhas: m.TotalLinhas,
		CreatedAt:   m.CreatedAt,
	}
}

// ToListMonitoramentosResponse converte lista para DTO
func ToListMonitoramentosResponse(items []*domain.Monitoramento, page, pageSize, total int) ListMonitoramentosResponse {
	data := make([]MonitoramentoResponse, len(items))
	for i, m := range items {
		data[i] = ToMonitoramentoResponse(m)
	}

	return ListMonitoramentosResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	}
}
