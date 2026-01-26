package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	sharedErrors "agro-monitoring/internal/shared/errors"
)

// PragasData representa o campo JSONB com informações das pragas
type PragasData struct {
	Pragas map[string]PragaInfo `json:"pragas"`
}

// PragaInfo contém informações de uma praga específica
type PragaInfo struct {
	Presente   bool                     `json:"presente"`
	Aplicacoes []AplicacaoHerbicidaJson `json:"aplicacoes,omitempty"`
}

// AplicacaoHerbicida representa uma aplicação de herbicida
type AplicacaoHerbicida struct {
	Posicao   int     `json:"posicao"`
	Herbicida string  `json:"herbicida"`
	Dose      float64 `json:"dose"`
}

// NewPragasData cria uma nova instância de PragasData
func NewPragasData() PragasData {
	return PragasData{
		Pragas: make(map[string]PragaInfo),
	}
}

// AddPraga adiciona uma praga presente
func (p *PragasData) AddPraga(nome string) {
	if p.Pragas == nil {
		p.Pragas = make(map[string]PragaInfo)
	}
	p.Pragas[nome] = PragaInfo{
		Presente:   true,
		Aplicacoes: []AplicacaoHerbicidaJson{},
	}
}

// AddAplicacao adiciona ou atualiza uma aplicação de herbicida numa posição específica (upsert)
func (p *PragasData) AddAplicacao(praga string, posicao int, herbicida string, dose float64) error {
	info, exists := p.Pragas[praga]
	if !exists {
		return sharedErrors.ErrPragaNotFound
	}

	// Procura aplicação existente na posição
	found := false
	for i, app := range info.Aplicacoes {
		if app.Posicao == posicao {
			// Atualiza aplicação existente
			info.Aplicacoes[i].Herbicida = herbicida
			info.Aplicacoes[i].Dose = dose
			found = true
			break
		}
	}

	// Se não encontrou, adiciona nova
	if !found {
		info.Aplicacoes = append(info.Aplicacoes, AplicacaoHerbicidaJson{
			Posicao:   posicao,
			Herbicida: herbicida,
			Dose:      dose,
		})
	}

	p.Pragas[praga] = info
	return nil
}

// GetPragasPresentes retorna lista de pragas presentes
func (p *PragasData) GetPragasPresentes() []string {
	var pragas []string
	for nome, info := range p.Pragas {
		if info.Presente {
			pragas = append(pragas, nome)
		}
	}
	return pragas
}

// HasPraga verifica se uma praga está presente
func (p *PragasData) HasPraga(nome string) bool {
	info, exists := p.Pragas[nome]
	return exists && info.Presente
}

// Value implementa driver.Valuer para PostgreSQL
func (p PragasData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implementa sql.Scanner para PostgreSQL
func (p *PragasData) Scan(value interface{}) error {
	if value == nil {
		p.Pragas = make(map[string]PragaInfo)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan PragasData: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, p)
}
