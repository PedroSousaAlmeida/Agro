package csv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockUUID() func() string {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("uuid-%d", counter)
	}
}

func TestParser_Parse_Basic(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote;Vassoura;Tiririca
1;Norte;Sub1;FAZ001;Fazenda São João;Q1;3;150,5;Argiloso;2;2020;Agosto;Nenhuma;S;N;S
2;Sul;Sub2;FAZ002;Fazenda Boa Vista;Q2;4;200,75;Arenoso;3;2019;Setembro;APP;N;S;N`

	parser := NewParser(mockUUID())
	result, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalLinhas)
	assert.Empty(t, result.Errors)

	// Verifica primeira área
	area1 := result.Areas[0]
	assert.Equal(t, "Norte", area1.Setor)
	assert.Equal(t, "FAZ001", area1.CodFazenda)
	assert.Equal(t, 150.5, area1.AreaTotal)
	assert.Equal(t, 3, area1.Corte)

	// Verifica pragas da primeira área
	assert.True(t, area1.PragasData.HasPraga("Camalote"))
	assert.False(t, area1.PragasData.HasPraga("Vassoura"))
	assert.True(t, area1.PragasData.HasPraga("Tiririca"))

	// Verifica segunda área
	area2 := result.Areas[1]
	assert.Equal(t, "FAZ002", area2.CodFazenda)
	assert.True(t, area2.PragasData.HasPraga("Vassoura"))
	assert.False(t, area2.PragasData.HasPraga("Camalote"))
}

func TestParser_Parse_VariousPragaFormats(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Praga1;Praga2;Praga3;Praga4
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N;S;SIM;1;X
2;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N;N;NAO;0;`

	parser := NewParser(mockUUID())
	result, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	require.NoError(t, err)

	// Primeira linha: todas as pragas devem estar presentes
	area1 := result.Areas[0]
	pragas := area1.PragasData.GetPragasPresentes()
	assert.Len(t, pragas, 4)

	// Segunda linha: nenhuma praga
	area2 := result.Areas[1]
	pragas2 := area2.PragasData.GetPragasPresentes()
	assert.Empty(t, pragas2)
}

func TestParser_Parse_EmptyCSV(t *testing.T) {
	csvContent := ``

	parser := NewParser(mockUUID())
	_, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	assert.Error(t, err)
}

func TestParser_Parse_OnlyHeader(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Praga1`

	parser := NewParser(mockUUID())
	result, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalLinhas)
}

func TestParser_Parse_FloatBrazilianFormat(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;1234,56;Arg;1;2020;Jan;N`

	parser := NewParser(mockUUID())
	result, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	require.NoError(t, err)
	assert.Equal(t, 1234.56, result.Areas[0].AreaTotal)
}

func TestParser_Parse_ManyPragas(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote;Grama seda;Coloniao;Vassoura;Braquiária;Mamona;Mucuna;Corda Viola;Tiririca
1;Norte;Sub1;FAZ001;Fazenda A;Q1;3;150,5;Argiloso;2;2020;Agosto;Nenhuma;S;S;N;S;N;N;S;N;S`

	parser := NewParser(mockUUID())
	result, err := parser.Parse(strings.NewReader(csvContent), "mon-123")

	require.NoError(t, err)

	area := result.Areas[0]
	pragas := area.PragasData.GetPragasPresentes()

	expectedPragas := []string{"Camalote", "Grama seda", "Vassoura", "Mucuna", "Tiririca"}
	assert.Len(t, pragas, 5)

	for _, expected := range expectedPragas {
		assert.True(t, area.PragasData.HasPraga(expected), "Deveria ter praga: %s", expected)
	}

	notExpected := []string{"Coloniao", "Braquiária", "Mamona", "Corda Viola"}
	for _, ne := range notExpected {
		assert.False(t, area.PragasData.HasPraga(ne), "NÃO deveria ter praga: %s", ne)
	}
}

func TestParser_Parse_MonitoramentoID(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N`

	parser := NewParser(mockUUID())
	result, _ := parser.Parse(strings.NewReader(csvContent), "meu-monitoramento-id")

	assert.Equal(t, "meu-monitoramento-id", result.Areas[0].MonitoramentoID)
}

func TestParser_Parse_UUIDGeneration(t *testing.T) {
	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N
2;N;S;F2;Fazenda;Q;1;100;Arg;1;2020;Jan;N`

	parser := NewParser(mockUUID())
	result, _ := parser.Parse(strings.NewReader(csvContent), "mon-123")

	assert.NotEqual(t, result.Areas[0].ID, result.Areas[1].ID)
	assert.Equal(t, "uuid-1", result.Areas[0].ID)
	assert.Equal(t, "uuid-2", result.Areas[1].ID)
}
