package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPragasData(t *testing.T) {
	pd := NewPragasData()

	assert.NotNil(t, pd.Pragas)
	assert.Empty(t, pd.Pragas)
}

func TestPragasData_AddPraga(t *testing.T) {
	pd := NewPragasData()

	pd.AddPraga("Camalote")
	pd.AddPraga("Vassoura")

	assert.Len(t, pd.Pragas, 2)
	assert.True(t, pd.Pragas["Camalote"].Presente)
	assert.True(t, pd.Pragas["Vassoura"].Presente)
	assert.Empty(t, pd.Pragas["Camalote"].Aplicacoes)
}

func TestPragasData_AddPraga_NilMap(t *testing.T) {
	pd := PragasData{}

	pd.AddPraga("Camalote")

	assert.NotNil(t, pd.Pragas)
	assert.True(t, pd.HasPraga("Camalote"))
}

func TestPragasData_AddAplicacao(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")

	err := pd.AddAplicacao("Camalote", 1, "Boral", 1.40)

	require.NoError(t, err)
	assert.Len(t, pd.Pragas["Camalote"].Aplicacoes, 1)
	assert.Equal(t, 1, pd.Pragas["Camalote"].Aplicacoes[0].Posicao)
	assert.Equal(t, "Boral", pd.Pragas["Camalote"].Aplicacoes[0].Herbicida)
	assert.Equal(t, 1.40, pd.Pragas["Camalote"].Aplicacoes[0].Dose)
}

func TestPragasData_AddAplicacao_MultipleApplications(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")

	pd.AddAplicacao("Camalote", 1, "Boral", 1.40)
	pd.AddAplicacao("Camalote", 2, "Roundup", 2.00)

	assert.Len(t, pd.Pragas["Camalote"].Aplicacoes, 2)
}

func TestPragasData_AddAplicacao_Upsert(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")

	// Adiciona posição 1
	pd.AddAplicacao("Camalote", 1, "Boral", 1.40)
	// Atualiza posição 1 (upsert)
	pd.AddAplicacao("Camalote", 1, "Hexagon", 2.50)

	assert.Len(t, pd.Pragas["Camalote"].Aplicacoes, 1)
	assert.Equal(t, "Hexagon", pd.Pragas["Camalote"].Aplicacoes[0].Herbicida)
	assert.Equal(t, 2.50, pd.Pragas["Camalote"].Aplicacoes[0].Dose)
}

func TestPragasData_AddAplicacao_PragaNotFound(t *testing.T) {
	pd := NewPragasData()

	err := pd.AddAplicacao("Inexistente", 1, "Boral", 1.40)

	assert.Error(t, err)
}

func TestPragasData_GetPragasPresentes(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")
	pd.AddPraga("Vassoura")
	pd.AddPraga("Tiririca")

	pragas := pd.GetPragasPresentes()

	assert.Len(t, pragas, 3)
	assert.Contains(t, pragas, "Camalote")
	assert.Contains(t, pragas, "Vassoura")
	assert.Contains(t, pragas, "Tiririca")
}

func TestPragasData_GetPragasPresentes_Empty(t *testing.T) {
	pd := NewPragasData()

	pragas := pd.GetPragasPresentes()

	assert.Empty(t, pragas)
}

func TestPragasData_HasPraga(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")

	assert.True(t, pd.HasPraga("Camalote"))
	assert.False(t, pd.HasPraga("Vassoura"))
	assert.False(t, pd.HasPraga(""))
}

func TestPragasData_Value(t *testing.T) {
	pd := NewPragasData()
	pd.AddPraga("Camalote")

	value, err := pd.Value()

	require.NoError(t, err)
	assert.NotNil(t, value)

	// Verifica se é JSON válido
	var result map[string]interface{}
	err = json.Unmarshal(value.([]byte), &result)
	require.NoError(t, err)
	assert.Contains(t, result, "pragas")
}

func TestPragasData_Scan(t *testing.T) {
	jsonData := []byte(`{"pragas":{"Camalote":{"presente":true,"aplicacoes":[{"posicao":1,"herbicida":"Boral","dose":1.4}]}}}`)

	pd := &PragasData{}
	err := pd.Scan(jsonData)

	require.NoError(t, err)
	assert.True(t, pd.HasPraga("Camalote"))
	assert.Len(t, pd.Pragas["Camalote"].Aplicacoes, 1)
	assert.Equal(t, 1, pd.Pragas["Camalote"].Aplicacoes[0].Posicao)
	assert.Equal(t, "Boral", pd.Pragas["Camalote"].Aplicacoes[0].Herbicida)
}

func TestPragasData_Scan_Nil(t *testing.T) {
	pd := &PragasData{}
	err := pd.Scan(nil)

	require.NoError(t, err)
	assert.NotNil(t, pd.Pragas)
	assert.Empty(t, pd.Pragas)
}

func TestPragasData_Scan_InvalidType(t *testing.T) {
	pd := &PragasData{}
	err := pd.Scan("invalid")

	assert.Error(t, err)
}

func TestPragasData_Scan_InvalidJSON(t *testing.T) {
	pd := &PragasData{}
	err := pd.Scan([]byte("invalid json"))

	assert.Error(t, err)
}
