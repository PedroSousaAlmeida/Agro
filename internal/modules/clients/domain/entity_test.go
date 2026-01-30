package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("client-1", "Usina ABC", "usina-abc", 10)

	assert.Equal(t, "client-1", client.ID)
	assert.Equal(t, "Usina ABC", client.Name)
	assert.Equal(t, "usina-abc", client.Slug)
	assert.Equal(t, 10, client.MaxUsers)
	assert.True(t, client.Active)
	assert.NotNil(t, client.Metadata)
	assert.False(t, client.CreatedAt.IsZero())
	assert.False(t, client.UpdatedAt.IsZero())
}

func TestNewClientUser(t *testing.T) {
	cu := NewClientUser("cu-1", "client-1", "user-1", "user@test.com", "user")

	assert.Equal(t, "cu-1", cu.ID)
	assert.Equal(t, "client-1", cu.ClientID)
	assert.Equal(t, "user-1", cu.UserID)
	assert.Equal(t, "user@test.com", cu.Email)
	assert.Equal(t, "user", cu.Role)
	assert.True(t, cu.Active)
	assert.False(t, cu.CreatedAt.IsZero())
}
