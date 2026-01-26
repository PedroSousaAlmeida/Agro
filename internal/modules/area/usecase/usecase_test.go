package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"agro-monitoring/internal/modules/area/repository"
	"agro-monitoring/internal/services/csv"
	monitoringRepo "agro-monitoring/internal/modules/monitoring/repository"
	monitoringUsecase "agro-monitoring/internal/modules/monitoring/usecase"

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

func setupAreaTest() (monitoringUsecase.MonitoringUseCase, AreaQueryUseCase, *repository.InMemoryRepository) {
	monRepo := monitoringRepo.NewInMemoryRepository()
	areaRepository := repository.NewInMemoryRepository()
	uuidGen := mockUUID()
	parser := csv.NewParser(uuidGen)

	monUC := monitoringUsecase.NewMonitoringUseCase(monRepo, areaRepository, parser, uuidGen)
	areaUC := NewAreaQueryUseCase(areaRepository)

	return monUC, areaUC, areaRepository
}

func TestAreaQueryUseCase_GetAreasByMonitoramento(t *testing.T) {
	monUC, areaUC, _ := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote
1;N;S;FAZ001;Fazenda A;Q1;1;100;Arg;1;2020;Jan;N;S
2;N;S;FAZ002;Fazenda B;Q2;2;200;Are;2;2021;Fev;N;N`

	mon, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	areas, total, err := areaUC.GetAreasByMonitoramento(context.Background(), mon.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, areas, 2)
}

func TestAreaQueryUseCase_GetAreaByID(t *testing.T) {
	monUC, areaUC, areaRepository := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;Norte;Sub1;FAZ001;Fazenda A;Q1;3;150,5;Argiloso;2;2020;Agosto;Nenhuma`

	mon, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	areas, _, _ := areaRepository.GetByMonitoramentoID(context.Background(), mon.ID, 10, 0)

	found, err := areaUC.GetAreaByID(context.Background(), areas[0].ID)
	require.NoError(t, err)
	assert.Equal(t, "FAZ001", found.CodFazenda)
}

func TestAreaQueryUseCase_SearchByFazenda(t *testing.T) {
	monUC, areaUC, _ := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição
1;N;S;FAZ001;Fazenda A;Q1;1;100;Arg;1;2020;Jan;N
2;N;S;FAZ001;Fazenda A;Q2;2;200;Are;2;2021;Fev;N
3;N;S;FAZ002;Fazenda B;Q3;3;300;Arg;3;2022;Mar;N`

	_, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	areas, total, err := areaUC.SearchByFazenda(context.Background(), "FAZ001", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)

	for _, a := range areas {
		assert.Equal(t, "FAZ001", a.CodFazenda)
	}
}

func TestAreaQueryUseCase_SearchByPraga(t *testing.T) {
	monUC, areaUC, _ := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote;Vassoura
1;N;S;FAZ001;Fazenda A;Q1;1;100;Arg;1;2020;Jan;N;S;N
2;N;S;FAZ002;Fazenda B;Q2;2;200;Are;2;2021;Fev;N;S;S
3;N;S;FAZ003;Fazenda C;Q3;3;300;Arg;3;2022;Mar;N;N;S`

	_, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	// Busca por Camalote - deve retornar 2 áreas
	areas, total, err := areaUC.SearchByPraga(context.Background(), "Camalote", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)

	for _, a := range areas {
		assert.True(t, a.PragasData.HasPraga("Camalote"))
	}

	// Busca por Vassoura - deve retornar 2 áreas
	areasV, totalV, _ := areaUC.SearchByPraga(context.Background(), "Vassoura", 1, 10)
	assert.Equal(t, 2, totalV)

	for _, a := range areasV {
		assert.True(t, a.PragasData.HasPraga("Vassoura"))
	}
}

func TestAreaQueryUseCase_AddAplicacaoHerbicida(t *testing.T) {
	monUC, areaUC, areaRepository := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote
1;N;S;FAZ001;Fazenda A;Q1;1;100;Arg;1;2020;Jan;N;S`

	mon, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	areas, _, _ := areaRepository.GetByMonitoramentoID(context.Background(), mon.ID, 10, 0)
	areaID := areas[0].ID

	err = areaUC.AddAplicacaoHerbicida(context.Background(), areaID, "Camalote", 1, "Boral", 1.40)
	require.NoError(t, err)

	updated, _ := areaUC.GetAreaByID(context.Background(), areaID)
	info := updated.PragasData.Pragas["Camalote"]

	assert.Len(t, info.Aplicacoes, 1)
	assert.Equal(t, 1, info.Aplicacoes[0].Posicao)
	assert.Equal(t, "Boral", info.Aplicacoes[0].Herbicida)
	assert.Equal(t, 1.40, info.Aplicacoes[0].Dose)
}

func TestAreaQueryUseCase_AddAplicacaoHerbicida_PragaNotFound(t *testing.T) {
	monUC, areaUC, areaRepository := setupAreaTest()

	csvContent := `Id;Setor;Setor2;Cod.Fazenda;Desc.Fazenda;Quadra;Corte;Área Total;Desc. Textura Solo;Corte Atual;Reforma;Mês Colheita;Restrição;Camalote
1;N;S;FAZ001;Fazenda A;Q1;1;100;Arg;1;2020;Jan;N;S`

	mon, err := monUC.UploadAndProcessCSV(context.Background(), strings.NewReader(csvContent), "teste.csv")
	require.NoError(t, err)

	areas, _, _ := areaRepository.GetByMonitoramentoID(context.Background(), mon.ID, 10, 0)
	areaID := areas[0].ID

	err = areaUC.AddAplicacaoHerbicida(context.Background(), areaID, "PragaInexistente", 1, "Boral", 1.40)
	assert.Error(t, err)
}
