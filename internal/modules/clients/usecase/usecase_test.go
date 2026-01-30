package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"agro-monitoring/internal/modules/clients/domain"
	"agro-monitoring/internal/modules/clients/dto"
	"agro-monitoring/internal/modules/clients/repository"
	"agro-monitoring/internal/modules/clients/service"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

func mockUUID() func() string {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("uuid-%d", counter)
	}
}

func setupClientTest() (ClientUseCase, domain.ClientRepository, domain.ClientUserRepository, service.KeycloakService) {
	clientRepo := repository.NewInMemoryRepository()
	clientUserRepo := repository.NewInMemoryClientUserRepository()
	keycloakSvc := service.NewInMemoryKeycloakService()
	uuidGen := mockUUID()

	uc := NewClientUseCase(clientRepo, clientUserRepo, keycloakSvc, uuidGen)
	return uc, clientRepo, clientUserRepo, keycloakSvc
}

func TestClientUseCase_CreateClient(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.CreateClientRequest{
		Name:     "Usina Santa Clara",
		Slug:     "usina-santa-clara",
		MaxUsers: 15,
		Metadata: map[string]interface{}{"cidade": "Ribeirão Preto"},
	}

	client, err := uc.CreateClient(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "Usina Santa Clara", client.Name)
	assert.Equal(t, "usina-santa-clara", client.Slug)
	assert.Equal(t, 15, client.MaxUsers)
	assert.True(t, client.Active)
	assert.NotEmpty(t, client.KeycloakGroupID)
	assert.Equal(t, "Ribeirão Preto", client.Metadata["cidade"])
}

func TestClientUseCase_CreateClient_InvalidSlug(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.CreateClientRequest{
		Name:     "Test",
		Slug:     "Invalid Slug!",
		MaxUsers: 10,
	}

	_, err := uc.CreateClient(context.Background(), req)

	assert.Error(t, err)
}

func TestClientUseCase_CreateClient_DuplicateSlug(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.CreateClientRequest{Name: "Test 1", Slug: "test-slug", MaxUsers: 10}
	_, err := uc.CreateClient(context.Background(), req)
	require.NoError(t, err)

	// Tentar criar com mesmo slug
	req2 := dto.CreateClientRequest{Name: "Test 2", Slug: "test-slug", MaxUsers: 10}
	_, err = uc.CreateClient(context.Background(), req2)

	assert.Error(t, err)
	assert.Equal(t, sharedErrors.ErrInvalidSlug, err)
}

func TestClientUseCase_GetClient(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.CreateClientRequest{Name: "Test", Slug: "test", MaxUsers: 10}
	created, _ := uc.CreateClient(context.Background(), req)

	client, err := uc.GetClient(context.Background(), created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, client.ID)
	assert.Equal(t, "Test", client.Name)
}

func TestClientUseCase_GetClient_NotFound(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	_, err := uc.GetClient(context.Background(), "nonexistent-id")

	assert.Error(t, err)
	assert.Equal(t, sharedErrors.ErrClientNotFound, err)
}

func TestClientUseCase_GetClientBySlug(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.CreateClientRequest{Name: "Test", Slug: "test-slug", MaxUsers: 10}
	created, _ := uc.CreateClient(context.Background(), req)

	client, err := uc.GetClientBySlug(context.Background(), "test-slug")

	require.NoError(t, err)
	assert.Equal(t, created.ID, client.ID)
	assert.Equal(t, "test-slug", client.Slug)
}

func TestClientUseCase_ListClients(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	// Criar 3 clients
	for i := 1; i <= 3; i++ {
		req := dto.CreateClientRequest{
			Name:     fmt.Sprintf("Client %d", i),
			Slug:     fmt.Sprintf("client-%d", i),
			MaxUsers: 10,
		}
		_, err := uc.CreateClient(context.Background(), req)
		require.NoError(t, err)
	}

	// Listar
	clients, total, err := uc.ListClients(context.Background(), 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, clients, 3)
}

func TestClientUseCase_RegisterUser(t *testing.T) {
	uc, _, clientUserRepo, _ := setupClientTest()

	// Criar client
	reqClient := dto.CreateClientRequest{Name: "Usina", Slug: "usina-test", MaxUsers: 10}
	client, _ := uc.CreateClient(context.Background(), reqClient)

	// Registrar usuário
	reqUser := dto.RegisterUserRequest{
		Email:     "joao@usina.com",
		Password:  "senha123",
		FirstName: "João",
		LastName:  "Silva",
	}

	clientUser, err := uc.RegisterUser(context.Background(), "usina-test", reqUser)

	require.NoError(t, err)
	assert.Equal(t, client.ID, clientUser.ClientID)
	assert.Equal(t, "joao@usina.com", clientUser.Email)
	assert.Equal(t, "user", clientUser.Role)
	assert.True(t, clientUser.Active)

	// Verificar que foi salvo no repositório
	count, _ := clientUserRepo.CountActiveByClient(context.Background(), client.ID)
	assert.Equal(t, 1, count)
}

func TestClientUseCase_RegisterUser_ClientNotFound(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	req := dto.RegisterUserRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	_, err := uc.RegisterUser(context.Background(), "nonexistent-slug", req)

	assert.Error(t, err)
	assert.Equal(t, sharedErrors.ErrClientNotFound, err)
}

func TestClientUseCase_RegisterUser_ClientInactive(t *testing.T) {
	uc, clientRepo, _, _ := setupClientTest()

	// Criar client
	reqClient := dto.CreateClientRequest{Name: "Test", Slug: "test", MaxUsers: 10}
	client, _ := uc.CreateClient(context.Background(), reqClient)

	// Desativar client
	client.Active = false
	clientRepo.Update(context.Background(), client)

	// Tentar registrar
	reqUser := dto.RegisterUserRequest{Email: "test@test.com", Password: "pass", FirstName: "Test", LastName: "User"}
	_, err := uc.RegisterUser(context.Background(), "test", reqUser)

	assert.Error(t, err)
	assert.Equal(t, sharedErrors.ErrClientInactive, err)
}

func TestClientUseCase_RegisterUser_UserLimitReached(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	// Criar client com limite de 2 usuários
	reqClient := dto.CreateClientRequest{Name: "Test", Slug: "test", MaxUsers: 2}
	_, _ = uc.CreateClient(context.Background(), reqClient)

	// Registrar 2 usuários (limite)
	reqUser1 := dto.RegisterUserRequest{Email: "user1@test.com", Password: "pass", FirstName: "User", LastName: "One"}
	_, err := uc.RegisterUser(context.Background(), "test", reqUser1)
	require.NoError(t, err)

	reqUser2 := dto.RegisterUserRequest{Email: "user2@test.com", Password: "pass", FirstName: "User", LastName: "Two"}
	_, err = uc.RegisterUser(context.Background(), "test", reqUser2)
	require.NoError(t, err)

	// Tentar registrar o 3º (deve falhar)
	reqUser3 := dto.RegisterUserRequest{Email: "user3@test.com", Password: "pass", FirstName: "User", LastName: "Three"}
	_, err = uc.RegisterUser(context.Background(), "test", reqUser3)

	assert.Error(t, err)
	assert.Equal(t, sharedErrors.ErrClientUserLimitReached, err)
}

func TestClientUseCase_CheckUserLimit(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	// Criar client com limite 2
	reqClient := dto.CreateClientRequest{Name: "Test", Slug: "test", MaxUsers: 2}
	client, _ := uc.CreateClient(context.Background(), reqClient)

	// Nenhum usuário - deve permitir
	canAdd, err := uc.CheckUserLimit(context.Background(), client.ID)
	require.NoError(t, err)
	assert.True(t, canAdd)

	// Adicionar 2 usuários
	reqUser1 := dto.RegisterUserRequest{Email: "user1@test.com", Password: "pass", FirstName: "User", LastName: "One"}
	uc.RegisterUser(context.Background(), "test", reqUser1)
	reqUser2 := dto.RegisterUserRequest{Email: "user2@test.com", Password: "pass", FirstName: "User", LastName: "Two"}
	uc.RegisterUser(context.Background(), "test", reqUser2)

	// Limite atingido - não deve permitir
	canAdd, err = uc.CheckUserLimit(context.Background(), client.ID)
	require.NoError(t, err)
	assert.False(t, canAdd)
}

func TestClientUseCase_ListClientUsers(t *testing.T) {
	uc, _, _, _ := setupClientTest()

	// Criar client
	reqClient := dto.CreateClientRequest{Name: "Test", Slug: "test", MaxUsers: 10}
	client, _ := uc.CreateClient(context.Background(), reqClient)

	// Registrar 2 usuários
	reqUser1 := dto.RegisterUserRequest{Email: "user1@test.com", Password: "pass", FirstName: "User", LastName: "One"}
	uc.RegisterUser(context.Background(), "test", reqUser1)
	reqUser2 := dto.RegisterUserRequest{Email: "user2@test.com", Password: "pass", FirstName: "User", LastName: "Two"}
	uc.RegisterUser(context.Background(), "test", reqUser2)

	// Listar
	users, total, err := uc.ListClientUsers(context.Background(), client.ID, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
}
