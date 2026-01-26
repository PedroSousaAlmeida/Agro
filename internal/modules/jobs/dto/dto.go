package dto

import (
	"encoding/json"
	"time"

	"agro-monitoring/internal/modules/jobs/domain"
)

// BulkAplicacoesRequest request para criar job de aplicações em massa
type BulkAplicacoesRequest struct {
	Aplicacoes []AplicacaoItemRequest `json:"aplicacoes"`
}

// AplicacaoItemRequest item de aplicação no request
type AplicacaoItemRequest struct {
	AreaID    string  `json:"area_id"`
	Praga     string  `json:"praga"`
	Posicao   int     `json:"posicao"`
	Herbicida string  `json:"herbicida"`
	Dose      float64 `json:"dose"`
}

// ToPayload converte request para payload de domínio
func (r *BulkAplicacoesRequest) ToPayload() domain.BulkAplicacoesPayload {
	items := make([]domain.AplicacaoItem, len(r.Aplicacoes))
	for i, a := range r.Aplicacoes {
		items[i] = domain.AplicacaoItem{
			AreaID:    a.AreaID,
			Praga:     a.Praga,
			Posicao:   a.Posicao,
			Herbicida: a.Herbicida,
			Dose:      a.Dose,
		}
	}
	return domain.BulkAplicacoesPayload{Aplicacoes: items}
}

// JobResponse resposta de job
type JobResponse struct {
	ID             string           `json:"id"`
	Type           string           `json:"type"`
	Status         string           `json:"status"`
	Progress       int              `json:"progress"`
	TotalItems     int              `json:"total_items"`
	ProcessedItems int              `json:"processed_items"`
	ErrorCount     int              `json:"error_count"`
	Errors         []JobErrorDetail `json:"errors,omitempty"`
	StartedAt      *time.Time       `json:"started_at,omitempty"`
	CompletedAt    *time.Time       `json:"completed_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}

// JobErrorDetail detalhe de erro
type JobErrorDetail struct {
	Line    int    `json:"line,omitempty"`
	ItemID  string `json:"item_id,omitempty"`
	Message string `json:"message"`
}

// ToJobResponse converte domain para DTO
func ToJobResponse(j *domain.Job) JobResponse {
	resp := JobResponse{
		ID:             j.ID,
		Type:           string(j.Type),
		Status:         string(j.Status),
		Progress:       j.Progress,
		TotalItems:     j.TotalItems,
		ProcessedItems: j.ProcessedItems,
		ErrorCount:     j.ErrorCount,
		StartedAt:      j.StartedAt,
		CompletedAt:    j.CompletedAt,
		CreatedAt:      j.CreatedAt,
	}

	// Parse error details se existir
	if len(j.ErrorDetails) > 0 {
		var errors []JobErrorDetail
		if err := json.Unmarshal(j.ErrorDetails, &errors); err == nil {
			resp.Errors = errors
		}
	}

	return resp
}

// CreateJobResponse resposta simplificada ao criar job
type CreateJobResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
