package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAreaMonitoramento(t *testing.T) {
	area := NewAreaMonitoramento("area-id", "mon-id")

	assert.Equal(t, "area-id", area.ID)
	assert.Equal(t, "mon-id", area.MonitoramentoID)
	assert.NotNil(t, area.PragasData.Pragas)
	assert.False(t, area.CreatedAt.IsZero())
}

func TestAreaMonitoramento_SetDadosCampo(t *testing.T) {
	area := NewAreaMonitoramento("area-id", "mon-id")

	area.SetDadosCampo(
		"Norte", "Sub1", "FAZ001", "Fazenda A", "Q1",
		3, 150.5, "Argiloso",
		2, "2020", "Agosto", "Nenhuma",
	)

	assert.Equal(t, "Norte", area.Setor)
	assert.Equal(t, "Sub1", area.Setor2)
	assert.Equal(t, "FAZ001", area.CodFazenda)
	assert.Equal(t, "Fazenda A", area.DescFazenda)
	assert.Equal(t, "Q1", area.Quadra)
	assert.Equal(t, 3, area.Corte)
	assert.Equal(t, 150.5, area.AreaTotal)
	assert.Equal(t, "Argiloso", area.DescTexturaSolo)
	assert.Equal(t, 2, area.CorteAtual)
	assert.Equal(t, "2020", area.Reforma)
	assert.Equal(t, "Agosto", area.MesColheita)
	assert.Equal(t, "Nenhuma", area.Restricao)
}

func TestAreaMonitoramento_AddPraga(t *testing.T) {
	area := NewAreaMonitoramento("area-id", "mon-id")

	area.AddPraga("Camalote")
	area.AddPraga("Vassoura")

	assert.True(t, area.PragasData.HasPraga("Camalote"))
	assert.True(t, area.PragasData.HasPraga("Vassoura"))
	assert.False(t, area.PragasData.HasPraga("Tiririca"))
}
