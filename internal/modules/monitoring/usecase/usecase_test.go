package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"

	areaRepo "agro-monitoring/internal/modules/area/repository"
	"agro-monitoring/internal/services/csv"
	"agro-monitoring/internal/modules/monitoring/repository"

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

func TestMonitoringUseCase_UploadAndProcessCSV(t *testing.T) {
	monRepo := repository.NewInMemoryRepository()
	areaRepository := areaRepo.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	uc := NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote;Vassoura
1;Norte;Sub1;FAZ001;Fazenda A;Q1;3;150,5;Argiloso;2;2020;Agosto;Nenhuma;S;N
2;Sul;Sub2;FAZ002;Fazenda B;Q2;4;200,75;Arenoso;3;2019;Setembro;APP;N;S`

	result, err := uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")

	require.NoError(t, err)
	assert.Equal(t, "concluido", string(result.Status))
	assert.Equal(t, 2, result.TotalLinhas)
	assert.Equal(t, "teste.csv", result.NomeArquivo)

	// Verifica se as áreas foram salvas
	areas, total, _ := areaRepository.GetByMonitoramentoID(context.Background(), result.ID, 10, 0)
	assert.Equal(t, 2, total)

	// Verifica fazendas
	fazendas := []string{areas[0].CodFazenda, areas[1].CodFazenda}
	assert.Contains(t, fazendas, "FAZ001")
	assert.Contains(t, fazendas, "FAZ002")
}

func TestMonitoringUseCase_UploadAndProcessCSV_InvalidCSV(t *testing.T) {
	monRepo := repository.NewInMemoryRepository()
	areaRepository := areaRepo.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	uc := NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)

	csvContent := `Campo1;Campo2
1;2`

	_, err := uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "invalido.csv")

	assert.Error(t, err)

	mons, _, _ := monRepo.List(context.Background(), 10, 0)
	if len(mons) > 0 {
		assert.Equal(t, "erro", string(mons[0].Status))
	}
}

func TestMonitoringUseCase_GetMonitoramento(t *testing.T) {
	monRepo := repository.NewInMemoryRepository()
	areaRepository := areaRepo.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	uc := NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N`

	created, err := uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	found, err := uc.GetMonitoramento(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestMonitoringUseCase_ListMonitoramentos(t *testing.T) {
	monRepo := repository.NewInMemoryRepository()
	areaRepository := areaRepo.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	uc := NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N`

	uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste1.csv")
	uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste2.csv")
	uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste3.csv")

	list, total, err := uc.ListMonitoramentos(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, list, 3)
}

func TestMonitoringUseCase_Pagination(t *testing.T) {
	monRepo := repository.NewInMemoryRepository()
	areaRepository := areaRepo.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	uc := NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;F1;Fazenda;Q;1;100;Arg;1;2020;Jan;N`

	for i := 0; i < 5; i++ {
		uc.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), fmt.Sprintf("teste%d.csv", i))
	}

	list, total, _ := uc.ListMonitoramentos(context.Background(), 1, 2)
	assert.Equal(t, 5, total)
	assert.Len(t, list, 2)

	list2, _, _ := uc.ListMonitoramentos(context.Background(), 2, 2)
	assert.Len(t, list2, 2)

	list3, _, _ := uc.ListMonitoramentos(context.Background(), 3, 2)
	assert.Len(t, list3, 1)
}
