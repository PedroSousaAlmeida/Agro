package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"agro-monitoring/internal/modules/area/dto"
	"agro-monitoring/internal/modules/area/usecase"
	sharedErrors "agro-monitoring/internal/shared/errors"
	"agro-monitoring/internal/shared/response"
)

// Handler handler para áreas
type Handler struct {
	uc usecase.AreaQueryUseCase
}

// NewHandler cria novo handler
func NewHandler(uc usecase.AreaQueryUseCase) *Handler {
	return &Handler{uc: uc}
}

// RegisterRoutes registra as rotas de áreas
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/areas", func(r chi.Router) {
		r.Get("/", h.ListByMonitoramento)
		r.Get("/search/fazenda", h.SearchByFazenda)
		r.Get("/search/praga", h.SearchByPraga)
		r.Get("/{id}", h.GetByID)
		r.Post("/{id}/aplicacao", h.AddAplicacao)
	})
}

// GetByID retorna uma área
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	area, err := h.uc.GetAreaByID(r.Context(), id)
	if err != nil {
		if err == sharedErrors.ErrAreaMonitoramentoNotFound {
			respondError(w, http.StatusNotFound, "Área não encontrada")
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToAreaResponse(area))
}

// ListByMonitoramento lista áreas de um monitoramento
func (h *Handler) ListByMonitoramento(w http.ResponseWriter, r *http.Request) {
	monitoramentoID := r.URL.Query().Get("monitoramento_id")
	if monitoramentoID == "" {
		respondError(w, http.StatusBadRequest, "monitoramento_id é obrigatório")
		return
	}

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 10)

	items, total, err := h.uc.GetAreasByMonitoramento(r.Context(), monitoramentoID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao listar")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToListAreasResponse(items, page, pageSize, total))
}

// SearchByFazenda busca áreas por fazenda
func (h *Handler) SearchByFazenda(w http.ResponseWriter, r *http.Request) {
	cod := r.URL.Query().Get("cod")
	if cod == "" {
		respondError(w, http.StatusBadRequest, "cod é obrigatório")
		return
	}

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 10)

	items, total, err := h.uc.SearchByFazenda(r.Context(), cod, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao buscar")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToListAreasResponse(items, page, pageSize, total))
}

// SearchByPraga busca áreas por praga
func (h *Handler) SearchByPraga(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	if nome == "" {
		respondError(w, http.StatusBadRequest, "nome é obrigatório")
		return
	}

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 10)

	items, total, err := h.uc.SearchByPraga(r.Context(), nome, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao buscar")
		return
	}

	respondJSON(w, http.StatusOK, dto.ToListAreasResponse(items, page, pageSize, total))
}

// AddAplicacao adiciona aplicação de herbicida
func (h *Handler) AddAplicacao(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "id")

	var req dto.AddAplicacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.Praga == "" || req.Posicao < 1 || req.Herbicida == "" || req.Dose <= 0 {
		respondError(w, http.StatusBadRequest, "praga, posicao (>= 1), herbicida e dose são obrigatórios")
		return
	}

	err := h.uc.AddAplicacaoHerbicida(r.Context(), areaID, req.Praga, req.Posicao, req.Herbicida, req.Dose)
	if err != nil {
		if err == sharedErrors.ErrAreaMonitoramentoNotFound {
			respondError(w, http.StatusNotFound, "Área não encontrada")
			return
		}
		if err == sharedErrors.ErrPragaNotFound {
			respondError(w, http.StatusBadRequest, "Praga não encontrada nesta área")
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro ao adicionar aplicação")
		return
	}

	area, _ := h.uc.GetAreaByID(r.Context(), areaID)
	respondJSON(w, http.StatusOK, dto.ToAreaResponse(area))
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
