package domain

import (
	"time"
)

// AreaMonitoramento representa uma linha do CSV (uma área monitorada)
type AreaMonitoramento struct {
	ID              string
	MonitoramentoID string
	Setor           string
	Setor2          string
	CodFazenda      string
	DescFazenda     string
	Quadra          string
	Corte           int
	AreaTotal       float64
	DescTexturaSolo string
	CorteAtual      int
	Reforma         string
	MesColheita     string
	Restricao       string
	PragasData      PragasData
	Aplicacoes      []AplicacaoHerbicidaJson
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// AplicacaoHerbicidaJson representa uma aplicação de herbicida em uma área
type AplicacaoHerbicidaJson struct {
	Posicao   int       `json:"posicao"`
	Praga     string    `json:"praga,omitempty"`
	Herbicida string    `json:"herbicida"`
	Dose      float64   `json:"dose"`
	AppliedAt time.Time `json:"applied_at,omitempty"`
}

// NewAreaMonitoramento cria uma nova área de monitoramento
func NewAreaMonitoramento(id, monitoramentoID string) *AreaMonitoramento {
	now := time.Now()
	return &AreaMonitoramento{
		ID:              id,
		MonitoramentoID: monitoramentoID,
		PragasData:      NewPragasData(),
		Aplicacoes:      make([]AplicacaoHerbicidaJson, 0),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// SetDadosCampo define os dados dos campos fixos
func (a *AreaMonitoramento) SetDadosCampo(
	setor, setor2, codFazenda, descFazenda, quadra string,
	corte int, areaTotal float64, descTexturaSolo string,
	corteAtual int, reforma, mesColheita, restricao string,
) {
	a.Setor = setor
	a.Setor2 = setor2
	a.CodFazenda = codFazenda
	a.DescFazenda = descFazenda
	a.Quadra = quadra
	a.Corte = corte
	a.AreaTotal = areaTotal
	a.DescTexturaSolo = descTexturaSolo
	a.CorteAtual = corteAtual
	a.Reforma = reforma
	a.MesColheita = mesColheita
	a.Restricao = restricao
}

// AddPraga adiciona uma praga presente na área
func (a *AreaMonitoramento) AddPraga(nome string) {
	a.PragasData.AddPraga(nome)
}
