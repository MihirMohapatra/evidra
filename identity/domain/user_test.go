package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	orgID := uuid.New()
	u := NewUser(orgID, "test@example.com", "hash123", RoleAdmin)
	require.NotNil(t, u)
	assert.Equal(t, orgID, u.OrganizationID)
	assert.Equal(t, "test@example.com", u.Email)
	assert.Equal(t, "hash123", u.PasswordHash)
	assert.Equal(t, RoleAdmin, u.Role)
	assert.True(t, u.IsActive)
	assert.NotEqual(t, uuid.Nil, u.ID)
	assert.WithinDuration(t, time.Now(), u.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), u.UpdatedAt, time.Second)
}
