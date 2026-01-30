package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"agro-monitoring/internal/modules/clients/domain"
	"agro-monitoring/internal/modules/clients/dto"
	"agro-monitoring/internal/modules/clients/service"
	sharedErrors "agro-monitoring/internal/shared/errors"
)

// ClientUseCase define os casos de uso de clients
type ClientUseCase interface {
	CreateClient(ctx context.Context, req dto.CreateClientRequest) (*domain.Client, error)
	GetClient(ctx context.Context, id string) (*domain.Client, error)
	GetClientBySlug(ctx context.Context, slug string) (*domain.Client, error)
	ListClients(ctx context.Context, page, pageSize int) ([]*domain.Client, int, error)
	GetClientStats(ctx context.Context, clientID string) (*domain.ClientStats, error)

	RegisterUser(ctx context.Context, slug string, req dto.RegisterUserRequest) (*domain.ClientUser, error)
	CheckUserLimit(ctx context.Context, clientID string) (bool, error)
	ListClientUsers(ctx context.Context, clientID string, page, pageSize int) ([]*domain.ClientUser, int, error)
}

type clientUseCase struct {
	clientRepo     domain.ClientRepository
	clientUserRepo domain.ClientUserRepository
	keycloakSvc    service.KeycloakService
	uuidGen        func() string
}

// NewClientUseCase cria uma nova instância de ClientUseCase
func NewClientUseCase(
	clientRepo domain.ClientRepository,
	clientUserRepo domain.ClientUserRepository,
	keycloakSvc service.KeycloakService,
	uuidGen func() string,
) ClientUseCase {
	return &clientUseCase{
		clientRepo:     clientRepo,
		clientUserRepo: clientUserRepo,
		keycloakSvc:    keycloakSvc,
		uuidGen:        uuidGen,
	}
}

func (uc *clientUseCase) CreateClient(ctx context.Context, req dto.CreateClientRequest) (*domain.Client, error) {
	// Validar slug
	if err := validateSlug(req.Slug); err != nil {
		return nil, err
	}

	// Verificar se slug já existe
	existing, err := uc.clientRepo.GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar slug: %w", err)
	}
	if existing != nil {
		return nil, sharedErrors.ErrInvalidSlug
	}

	// Criar client
	client := domain.NewClient(uc.uuidGen(), req.Name, req.Slug, req.MaxUsers)
	if req.Metadata != nil {
		client.Metadata = req.Metadata
	}

	// Criar grupo no Keycloak
	groupName := fmt.Sprintf("/clients/%s", req.Slug)
	attrs := map[string][]string{
		"client_id": {client.ID},
		"max_users": {fmt.Sprintf("%d", req.MaxUsers)},
	}

	groupID, err := uc.keycloakSvc.CreateGroup(ctx, groupName, attrs)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar grupo no Keycloak: %w", err)
	}

	client.KeycloakGroupID = groupID

	// Salvar no banco
	if err := uc.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("erro ao salvar client: %w", err)
	}

	return client, nil
}

func (uc *clientUseCase) GetClient(ctx context.Context, id string) (*domain.Client, error) {
	client, err := uc.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, sharedErrors.ErrClientNotFound
	}
	return client, nil
}

func (uc *clientUseCase) GetClientBySlug(ctx context.Context, slug string) (*domain.Client, error) {
	client, err := uc.clientRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, sharedErrors.ErrClientNotFound
	}
	return client, nil
}

func (uc *clientUseCase) ListClients(ctx context.Context, page, pageSize int) ([]*domain.Client, int, error) {
	offset := (page - 1) * pageSize
	return uc.clientRepo.List(ctx, pageSize, offset)
}

func (uc *clientUseCase) GetClientStats(ctx context.Context, clientID string) (*domain.ClientStats, error) {
	stats, err := uc.clientRepo.GetStats(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, sharedErrors.ErrClientNotFound
	}
	return stats, nil
}

func (uc *clientUseCase) RegisterUser(ctx context.Context, slug string, req dto.RegisterUserRequest) (*domain.ClientUser, error) {
	// 1. Buscar client por slug
	client, err := uc.clientRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar client: %w", err)
	}
	if client == nil {
		return nil, sharedErrors.ErrClientNotFound
	}

	// 2. Validar se client está ativo
	if !client.Active {
		return nil, sharedErrors.ErrClientInactive
	}

	// 3. Contar usuários ativos
	currentUsers, err := uc.clientUserRepo.CountActiveByClient(ctx, client.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar usuários: %w", err)
	}

	// 4. Validar limite
	if currentUsers >= client.MaxUsers {
		return nil, sharedErrors.ErrClientUserLimitReached
	}

	// 5. Criar usuário no Keycloak
	username := strings.Split(req.Email, "@")[0] // username = parte antes do @
	kcUser := service.KeycloakUser{
		Username:  username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Enabled:   true,
	}

	userID, err := uc.keycloakSvc.CreateUser(ctx, kcUser)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário no Keycloak: %w", err)
	}

	// 6. Setar atributo client_id no usuário
	if err := uc.keycloakSvc.SetUserAttribute(ctx, userID, "client_id", client.ID); err != nil {
		return nil, fmt.Errorf("erro ao setar atributo client_id: %w", err)
	}

	// 7. Adicionar ao grupo do client
	if err := uc.keycloakSvc.AddUserToGroup(ctx, userID, client.KeycloakGroupID); err != nil {
		return nil, fmt.Errorf("erro ao adicionar ao grupo: %w", err)
	}

	// 8. Salvar em client_users
	clientUser := domain.NewClientUser(uc.uuidGen(), client.ID, userID, req.Email, "user")

	if err := uc.clientUserRepo.Create(ctx, clientUser); err != nil {
		// TODO: Compensação - tentar deletar do Keycloak
		return nil, fmt.Errorf("erro ao salvar client_user: %w", err)
	}

	return clientUser, nil
}

func (uc *clientUseCase) CheckUserLimit(ctx context.Context, clientID string) (bool, error) {
	client, err := uc.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return false, err
	}
	if client == nil {
		return false, sharedErrors.ErrClientNotFound
	}

	currentUsers, err := uc.clientUserRepo.CountActiveByClient(ctx, clientID)
	if err != nil {
		return false, err
	}

	return currentUsers < client.MaxUsers, nil
}

func (uc *clientUseCase) ListClientUsers(ctx context.Context, clientID string, page, pageSize int) ([]*domain.ClientUser, int, error) {
	offset := (page - 1) * pageSize
	return uc.clientUserRepo.ListByClient(ctx, clientID, pageSize, offset)
}

// validateSlug valida o formato do slug
func validateSlug(slug string) error {
	if len(slug) < 3 || len(slug) > 100 {
		return fmt.Errorf("slug deve ter entre 3 e 100 caracteres")
	}

	matched, err := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	if err != nil {
		return err
	}
	if !matched {
		return sharedErrors.ErrInvalidSlug
	}

	return nil
}
