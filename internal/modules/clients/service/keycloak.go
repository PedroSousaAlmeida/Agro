package service

import (
	"context"
	"fmt"

	"agro-monitoring/internal/config"
	"github.com/Nerzal/gocloak/v13"
)

// KeycloakService define operações de integração com Keycloak Admin API
type KeycloakService interface {
	CreateGroup(ctx context.Context, name string, attrs map[string][]string) (string, error)
	CreateUser(ctx context.Context, user KeycloakUser) (string, error)
	AddUserToGroup(ctx context.Context, userID, groupID string) error
	SetUserAttribute(ctx context.Context, userID, key, value string) error
}

// KeycloakUser representa dados para criar usuário no Keycloak
type KeycloakUser struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Password  string
	Enabled   bool
}

type keycloakService struct {
	client       *gocloak.GoCloak
	realm        string
	clientID     string
	clientSecret string
}

// NewKeycloakService cria uma nova instância do serviço Keycloak
func NewKeycloakService(env *config.Env) KeycloakService {
	client := gocloak.NewClient(env.KeycloakURL)
	return &keycloakService{
		client:       client,
		realm:        env.KeycloakRealm,
		clientID:     env.KeycloakAdminClientID,
		clientSecret: env.KeycloakAdminClientSecret,
	}
}

// getToken obtém token de acesso usando client credentials
func (s *keycloakService) getToken(ctx context.Context) (string, error) {
	token, err := s.client.LoginClient(ctx, s.clientID, s.clientSecret, s.realm)
	if err != nil {
		return "", fmt.Errorf("erro ao obter token admin: %w", err)
	}
	return token.AccessToken, nil
}

// CreateGroup cria um grupo no Keycloak
func (s *keycloakService) CreateGroup(ctx context.Context, name string, attrs map[string][]string) (string, error) {
	token, err := s.getToken(ctx)
	if err != nil {
		return "", err
	}

	group := gocloak.Group{
		Name:       gocloak.StringP(name),
		Attributes: &attrs,
	}

	groupID, err := s.client.CreateGroup(ctx, token, s.realm, group)
	if err != nil {
		return "", fmt.Errorf("erro ao criar grupo: %w", err)
	}

	return groupID, nil
}

// CreateUser cria um usuário no Keycloak
func (s *keycloakService) CreateUser(ctx context.Context, user KeycloakUser) (string, error) {
	token, err := s.getToken(ctx)
	if err != nil {
		return "", err
	}

	kcUser := gocloak.User{
		Username:      gocloak.StringP(user.Username),
		Email:         gocloak.StringP(user.Email),
		FirstName:     gocloak.StringP(user.FirstName),
		LastName:      gocloak.StringP(user.LastName),
		Enabled:       gocloak.BoolP(user.Enabled),
		EmailVerified: gocloak.BoolP(true),
	}

	userID, err := s.client.CreateUser(ctx, token, s.realm, kcUser)
	if err != nil {
		return "", fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Setar senha
	err = s.client.SetPassword(ctx, token, userID, s.realm, user.Password, false)
	if err != nil {
		// Tentar deletar usuário se falhou ao setar senha
		_ = s.client.DeleteUser(ctx, token, s.realm, userID)
		return "", fmt.Errorf("erro ao setar senha: %w", err)
	}

	return userID, nil
}

// AddUserToGroup adiciona um usuário a um grupo
func (s *keycloakService) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}

	err = s.client.AddUserToGroup(ctx, token, s.realm, userID, groupID)
	if err != nil {
		return fmt.Errorf("erro ao adicionar usuário ao grupo: %w", err)
	}

	return nil
}

// SetUserAttribute seta um atributo customizado no usuário
func (s *keycloakService) SetUserAttribute(ctx context.Context, userID, key, value string) error {
	token, err := s.getToken(ctx)
	if err != nil {
		return err
	}

	user, err := s.client.GetUserByID(ctx, token, s.realm, userID)
	if err != nil {
		return fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if user.Attributes == nil {
		attrs := make(map[string][]string)
		user.Attributes = &attrs
	}

	(*user.Attributes)[key] = []string{value}

	err = s.client.UpdateUser(ctx, token, s.realm, *user)
	if err != nil {
		return fmt.Errorf("erro ao atualizar atributo: %w", err)
	}

	return nil
}
