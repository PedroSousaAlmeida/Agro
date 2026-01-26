package domain

import (
	"encoding/json"
	"time"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type JobType string

const (
	JobTypeBulkAplicacoes JobType = "bulk_aplicacoes"
	JobTypeCSVImport      JobType = "csv_import"
)

// Job representa um trabalho em background
type Job struct {
	ID             string
	Type           JobType
	Status         JobStatus
	Payload        json.RawMessage
	Result         json.RawMessage
	Progress       int
	TotalItems     int
	ProcessedItems int
	ErrorCount     int
	ErrorDetails   json.RawMessage
	StartedAt      *time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewJob cria um novo job
func NewJob(id string, jobType JobType, payload interface{}) (*Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:        id,
		Type:      jobType,
		Status:    JobStatusPending,
		Payload:   payloadBytes,
		Progress:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Start marca o job como em processamento
func (j *Job) Start(totalItems int) {
	now := time.Now()
	j.Status = JobStatusProcessing
	j.StartedAt = &now
	j.TotalItems = totalItems
	j.UpdatedAt = now
}

// UpdateProgress atualiza o progresso
func (j *Job) UpdateProgress(processed int) {
	j.ProcessedItems = processed
	if j.TotalItems > 0 {
		j.Progress = (processed * 100) / j.TotalItems
	}
	j.UpdatedAt = time.Now()
}

// Complete marca como concluído
func (j *Job) Complete(result interface{}) error {
	now := time.Now()
	j.Status = JobStatusCompleted
	j.CompletedAt = &now
	j.Progress = 100
	j.UpdatedAt = now

	if result != nil {
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return err
		}
		j.Result = resultBytes
	}
	return nil
}

// Fail marca como falhou
func (j *Job) Fail(errors []JobError) error {
	now := time.Now()
	j.Status = JobStatusFailed
	j.CompletedAt = &now
	j.ErrorCount = len(errors)
	j.UpdatedAt = now

	if len(errors) > 0 {
		errBytes, err := json.Marshal(errors)
		if err != nil {
			return err
		}
		j.ErrorDetails = errBytes
	}
	return nil
}

// AddError incrementa contador de erros
func (j *Job) AddError() {
	j.ErrorCount++
	j.UpdatedAt = time.Now()
}

// JobError representa um erro durante processamento
type JobError struct {
	Line    int    `json:"line,omitempty"`
	ItemID  string `json:"item_id,omitempty"`
	Message string `json:"message"`
}

// BulkAplicacoesPayload payload para job de aplicações em massa
type BulkAplicacoesPayload struct {
	Aplicacoes []AplicacaoItem `json:"aplicacoes"`
}

// AplicacaoItem item de aplicação
type AplicacaoItem struct {
	AreaID    string  `json:"area_id"`
	Praga     string  `json:"praga"`
	Posicao   int     `json:"posicao"`
	Herbicida string  `json:"herbicida"`
	Dose      float64 `json:"dose"`
}

// BulkAplicacoesResult resultado do processamento
type BulkAplicacoesResult struct {
	Processed int `json:"processed"`
	Errors    int `json:"errors"`
}
