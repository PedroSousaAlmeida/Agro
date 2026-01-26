package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitoramentoStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   MonitoramentoStatus
		expected bool
	}{
		{StatusProcessando, true},
		{StatusConcluido, true},
		{StatusErro, true},
		{MonitoramentoStatus("invalido"), false},
		{MonitoramentoStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestNewMonitoramento(t *testing.T) {
	m := NewMonitoramento("test-id", "arquivo.csv")

	assert.Equal(t, "test-id", m.ID)
	assert.Equal(t, "arquivo.csv", m.NomeArquivo)
	assert.Equal(t, StatusProcessando, m.Status)
	assert.Equal(t, 0, m.TotalLinhas)
	assert.False(t, m.DataUpload.IsZero())
	assert.False(t, m.CreatedAt.IsZero())
	assert.False(t, m.UpdatedAt.IsZero())
}

func TestMonitoramento_MarkAsCompleted(t *testing.T) {
	m := NewMonitoramento("test-id", "arquivo.csv")
	oldUpdatedAt := m.UpdatedAt

	m.MarkAsCompleted(100)

	assert.Equal(t, StatusConcluido, m.Status)
	assert.Equal(t, 100, m.TotalLinhas)
	assert.True(t, m.UpdatedAt.After(oldUpdatedAt) || m.UpdatedAt.Equal(oldUpdatedAt))
}

func TestMonitoramento_MarkAsError(t *testing.T) {
	m := NewMonitoramento("test-id", "arquivo.csv")
	oldUpdatedAt := m.UpdatedAt

	m.MarkAsError()

	assert.Equal(t, StatusErro, m.Status)
	assert.True(t, m.UpdatedAt.After(oldUpdatedAt) || m.UpdatedAt.Equal(oldUpdatedAt))
}
