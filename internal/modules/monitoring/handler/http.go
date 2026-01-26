package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"agro-monitoring/internal/modules/monitoring/dto"
	"agro-monitoring/internal/modules/monitoring/usecase"
	sharedErrors "agro-monitoring/internal/shared/errors"
	"agro-monitoring/internal/shared/response"
)

// Handler handler para monitoramentos
type Handler struct {
	uc usecase.MonitoringUseCase
}

// NewHandler cria novo handler
func NewHandler(uc usecase.MonitoringUseCase) *Handler {
	return &Handler{uc: uc}
}

// RegisterRoutes registra as rotas de monitoramento
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/monitoramentos", func(r chi.Router) {
		r.Post("/", h.Upload)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
	})
}

// Upload processa upload de CSV
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Erro ao processar formulário: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Arquivo 'file' não encontrado")
		return
	}
	defer file.Close()

	mon, err := h.uc.UploadAndProcessCSV(r.Context(), file, header.Filename)
	if err != nil {
		if err == sharedErrors.ErrInvalidCSV {
			respondError(w, http.StatusBadRequest, "CSV inválido")
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro ao processar CSV")
		return
	}

	respondJSON(w, http.StatusCreated, dto.ToMonitoramentoResponse(mon))
}

// GetByID retorna um monitoramento
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mon, err := h.uc.GetMonitoramento(r.Context(), id)
	if err != nil {
		if err == sharedErrors.ErrMonitoramentoNotFound {
			respondError(w, http.StatusNotFound, "Monitoramento não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToMonitoramentoResponse(mon))
}

// List lista monitoramentos
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 10)

	items, total, err := h.uc.ListMonitoramentos(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao listar")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToListMonitoramentosResponse(items, page, pageSize, total))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, response.ErrorResponse{Message: message})
}

func getQueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
