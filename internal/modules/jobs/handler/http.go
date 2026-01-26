package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"agro-monitoring/internal/modules/jobs/dto"
	"agro-monitoring/internal/modules/jobs/usecase"
	sharedErrors "agro-monitoring/internal/shared/errors"
	"agro-monitoring/internal/shared/response"
)

// Handler handler para jobs
type Handler struct {
	uc usecase.JobUseCase
}

// NewHandler cria novo handler
func NewHandler(uc usecase.JobUseCase) *Handler {
	return &Handler{uc: uc}
}

// RegisterRoutes registra as rotas de jobs
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/jobs", func(r chi.Router) {
		r.Post("/aplicacoes", h.CreateBulkAplicacoes)
		r.Get("/{id}", h.GetJobStatus)
	})
}

// CreateBulkAplicacoes cria job de aplicações em massa
func (h *Handler) CreateBulkAplicacoes(w http.ResponseWriter, r *http.Request) {
	var req dto.BulkAplicacoesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if len(req.Aplicacoes) == 0 {
		respondError(w, http.StatusBadRequest, "Lista de aplicações vazia")
		return
	}

	// Valida itens
	for i, item := range req.Aplicacoes {
		if item.AreaID == "" || item.Praga == "" || item.Posicao < 1 || item.Herbicida == "" || item.Dose <= 0 {
			respondError(w, http.StatusBadRequest, "Item "+string(rune(i+1))+" inválido: area_id, praga, posicao (>= 1), herbicida e dose são obrigatórios")
			return
		}
	}

	payload := req.ToPayload()
	job, err := h.uc.CreateBulkAplicacoesJob(r.Context(), payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao criar job: "+err.Error())
		return
	}

	respondJSON(w, http.StatusAccepted, dto.CreateJobResponse{
		ID:      job.ID,
		Status:  string(job.Status),
		Message: "Job criado com sucesso. Use GET /v1/jobs/" + job.ID + " para acompanhar o progresso.",
	})
}

// GetJobStatus retorna status de um job
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, err := h.uc.GetJobStatus(r.Context(), id)
	if err != nil {
		if err == sharedErrors.ErrJobNotFound {
			respondError(w, http.StatusNotFound, "Job não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToJobResponse(job))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, response.ErrorResponse{Message: message})
}
