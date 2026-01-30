package service

import (
	"context"
	"fmt"
	"sync"
)

type InMemoryKeycloakService struct {
	mu         sync.RWMutex
	groups     map[string]string            // groupID -> name
	users      map[string]KeycloakUser      // userID -> user
	userGroups map[string][]string          // userID -> []groupID
	userAttrs  map[string]map[string]string // userID -> {key: value}
	counter    int
}

func NewInMemoryKeycloakService() KeycloakService {
	return &InMemoryKeycloakService{
		groups:     make(map[string]string),
		users:      make(map[string]KeycloakUser),
		userGroups: make(map[string][]string),
		userAttrs:  make(map[string]map[string]string),
	}
}

func (s *InMemoryKeycloakService) CreateGroup(ctx context.Context, name string, attrs map[string][]string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	groupID := fmt.Sprintf("group-%d", s.counter)
	s.groups[groupID] = name
	return groupID, nil
}

func (s *InMemoryKeycloakService) CreateUser(ctx context.Context, user KeycloakUser) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	userID := fmt.Sprintf("user-%d", s.counter)
	s.users[userID] = user
	s.userAttrs[userID] = make(map[string]string)
	s.userGroups[userID] = make([]string, 0)
	return userID, nil
}

func (s *InMemoryKeycloakService) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.userGroups[userID]; !ok {
		s.userGroups[userID] = make([]string, 0)
	}

	s.userGroups[userID] = append(s.userGroups[userID], groupID)
	return nil
}

func (s *InMemoryKeycloakService) SetUserAttribute(ctx context.Context, userID, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.userAttrs[userID]; !ok {
		s.userAttrs[userID] = make(map[string]string)
	}

	s.userAttrs[userID][key] = value
	return nil
}

func (s *InMemoryKeycloakService) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups = make(map[string]string)
	s.users = make(map[string]KeycloakUser)
	s.userGroups = make(map[string][]string)
	s.userAttrs = make(map[string]map[string]string)
	s.counter = 0
}
