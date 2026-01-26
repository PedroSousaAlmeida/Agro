package dto

import (
	"time"

	"agro-monitoring/internal/modules/area/domain"
)

// AreaResponse resposta de área
type AreaResponse struct {
	ID              string                 `json:"id"`
	MonitoramentoID string                 `json:"monitoramento_id"`
	Setor           string                 `json:"setor"`
	Setor2          string                 `json:"setor2"`
	CodFazenda      string                 `json:"cod_fazenda"`
	DescFazenda     string                 `json:"desc_fazenda"`
	Quadra          string                 `json:"quadra"`
	Corte           int                    `json:"corte"`
	AreaTotal       float64                `json:"area_total"`
	DescTexturaSolo string                 `json:"desc_textura_solo"`
	CorteAtual      int                    `json:"corte_atual"`
	Reforma         string                 `json:"reforma"`
	MesColheita     string                 `json:"mes_colheita"`
	Restricao       string                 `json:"restricao"`
	PragasData      map[string]interface{} `json:"pragas_data"`
	CreatedAt       time.Time              `json:"created_at"`
}

// ListAreasResponse resposta paginada de áreas
type ListAreasResponse struct {
	Data       []AreaResponse `json:"data"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalCount int            `json:"total_count"`
}

// AddAplicacaoRequest request para adicionar/atualizar aplicação
type AddAplicacaoRequest struct {
	Praga     string  `json:"praga"`
	Posicao   int     `json:"posicao"`
	Herbicida string  `json:"herbicida"`
	Dose      float64 `json:"dose"`
}

// ToAreaResponse converte domain para DTO
func ToAreaResponse(a *domain.AreaMonitoramento) AreaResponse {
	pragasMap := make(map[string]interface{})
	for nome, info := range a.PragasData.Pragas {
		pragasMap[nome] = map[string]interface{}{
			"presente":   info.Presente,
			"aplicacoes": info.Aplicacoes,
		}
	}

	return AreaResponse{
		ID:              a.ID,
		MonitoramentoID: a.MonitoramentoID,
		Setor:           a.Setor,
		Setor2:          a.Setor2,
		CodFazenda:      a.CodFazenda,
		DescFazenda:     a.DescFazenda,
		Quadra:          a.Quadra,
		Corte:           a.Corte,
		AreaTotal:       a.AreaTotal,
		DescTexturaSolo: a.DescTexturaSolo,
		CorteAtual:      a.CorteAtual,
		Reforma:         a.Reforma,
		MesColheita:     a.MesColheita,
		Restricao:       a.Restricao,
		PragasData:      pragasMap,
		CreatedAt:       a.CreatedAt,
	}
}

// ToListAreasResponse converte lista para DTO
func ToListAreasResponse(items []*domain.AreaMonitoramento, page, pageSize, total int) ListAreasResponse {
	data := make([]AreaResponse, len(items))
	for i, a := range items {
		data[i] = ToAreaResponse(a)
	}

	return ListAreasResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	}
}
