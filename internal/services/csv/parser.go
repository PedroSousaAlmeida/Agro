package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"agro-monitoring/internal/modules/area/domain"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

// Campos fixos esperados no CSV (em ordem)
var camposFixos = []string{
	"Id",
	"Setor",
	"Setor2",
	"Cod.Fazenda",
	"Desc.Fazenda",
	"Quadra",
	"Corte",
	"Área Total",
	"Desc. Textura Solo",
	"Corte Atual",
	"Reforma",
	"Mês Colheita",
	"Restrição",
}

// Parser processa arquivos CSV de monitoramento
type Parser struct {
	uuidGenerator func() string
}

// NewParser cria um novo parser
func NewParser(uuidGenerator func() string) *Parser {
	return &Parser{
		uuidGenerator: uuidGenerator,
	}
}

// ParseResult contém o resultado do parsing
type ParseResult struct {
	Areas       []*domain.AreaMonitoramento
	TotalLinhas int
	Errors      []ParseError
}

// ParseError representa um erro em uma linha específica
type ParseError struct {
	Linha int
	Erro  string
}

// Parse processa o CSV e retorna as áreas de monitoramento
func (p *Parser) Parse(reader io.Reader, monitoramentoID string) (*ParseResult, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: erro ao ler header: %v", sharedErrors.ErrInvalidCSV, err)
	}

	colIndex, pragaColumns, err := p.mapColumns(header)
	if err != nil {
		return nil, err
	}

	result := &ParseResult{
		Areas:  make([]*domain.AreaMonitoramento, 0),
		Errors: make([]ParseError, 0),
	}

	linha := 1
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		linha++

		if err != nil {
			result.Errors = append(result.Errors, ParseError{
				Linha: linha,
				Erro:  fmt.Sprintf("erro ao ler linha: %v", err),
			})
			continue
		}

		area, err := p.parseRecord(record, colIndex, pragaColumns, monitoramentoID)
		if err != nil {
			result.Errors = append(result.Errors, ParseError{
				Linha: linha,
				Erro:  err.Error(),
			})
			continue
		}

		result.Areas = append(result.Areas, area)
	}

	result.TotalLinhas = len(result.Areas)

	if result.TotalLinhas == 0 && len(result.Errors) > 0 {
		return nil, sharedErrors.ErrInvalidCSV
	}

	return result, nil
}

func (p *Parser) mapColumns(header []string) (map[string]int, []string, error) {
	colIndex := make(map[string]int)
	pragaColumns := make([]string, 0)

	restricaoIndex := -1
	herbIndex := -1

	for i, col := range header {
		col = strings.TrimSpace(col)
		header[i] = col
		colIndex[col] = i

		if col == "Restrição" || col == "Restricao" {
			restricaoIndex = i
		}
		// Encontra onde começam os herbicidas (fim das pragas)
		if herbIndex == -1 && strings.HasPrefix(col, "Herb ") {
			herbIndex = i
		}
	}

	for _, campo := range camposFixos[:4] {
		if _, ok := colIndex[campo]; !ok {
			found := false
			for headerCol := range colIndex {
				if strings.EqualFold(headerCol, campo) || strings.Contains(strings.ToLower(headerCol), strings.ToLower(campo)) {
					found = true
					break
				}
			}
			if !found {
				return nil, nil, fmt.Errorf("%w: campo obrigatório não encontrado: %s", sharedErrors.ErrInvalidCSV, campo)
			}
		}
	}

	if restricaoIndex >= 0 {
		endIndex := len(header)
		if herbIndex > restricaoIndex {
			endIndex = herbIndex
		}

		for i := restricaoIndex + 1; i < endIndex; i++ {
			colName := header[i]
			// Ignora colunas vazias e colunas "ColunaX"
			if colName == "" || strings.HasPrefix(colName, "Coluna") {
				continue
			}
			pragaColumns = append(pragaColumns, colName)
		}
	}

	return colIndex, pragaColumns, nil
}

func (p *Parser) parseRecord(record []string, colIndex map[string]int, pragaColumns []string, monitoramentoID string) (*domain.AreaMonitoramento, error) {
	area := domain.NewAreaMonitoramento(p.uuidGenerator(), monitoramentoID)

	setor := p.getString(record, colIndex, "Setor")
	setor2 := p.getString(record, colIndex, "Setor2")
	codFazenda := p.getString(record, colIndex, "Cod.Fazenda")
	descFazenda := p.getString(record, colIndex, "Desc.Fazenda")
	quadra := p.getString(record, colIndex, "Quadra")
	corte := p.getInt(record, colIndex, "Corte")
	areaTotal := p.getFloat(record, colIndex, "Área Total")
	descTexturaSolo := p.getString(record, colIndex, "Desc. Textura Solo")
	corteAtual := p.getInt(record, colIndex, "Corte Atual")
	reforma := p.getString(record, colIndex, "Reforma")
	mesColheita := p.getString(record, colIndex, "Mês Colheita")
	restricao := p.getString(record, colIndex, "Restrição")

	area.SetDadosCampo(
		setor, setor2, codFazenda, descFazenda, quadra,
		corte, areaTotal, descTexturaSolo,
		corteAtual, reforma, mesColheita, restricao,
	)

	for _, pragaName := range pragaColumns {
		idx, ok := colIndex[pragaName]
		if !ok || idx >= len(record) {
			continue
		}

		valor := strings.TrimSpace(strings.ToUpper(record[idx]))
		if valor == "S" || valor == "SIM" || valor == "1" || valor == "X" {
			area.AddPraga(pragaName)
		}
	}

	return area, nil
}

func (p *Parser) getString(record []string, colIndex map[string]int, colName string) string {
	idx, ok := colIndex[colName]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func (p *Parser) getInt(record []string, colIndex map[string]int, colName string) int {
	str := p.getString(record, colIndex, colName)
	if str == "" {
		return 0
	}
	val, _ := strconv.Atoi(str)
	return val
}

func (p *Parser) getFloat(record []string, colIndex map[string]int, colName string) float64 {
	str := p.getString(record, colIndex, colName)
	if str == "" {
		return 0
	}
	str = strings.Replace(str, ",", ".", 1)
	val, _ := strconv.ParseFloat(str, 64)
	return val
}
