package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIKey(t *testing.T) {
	orgID := uuid.New()
	k := NewAPIKey(orgID, "my-key", "hash123", "ev_abc")
	require.NotNil(t, k)
	assert.Equal(t, orgID, k.OrganizationID)
	assert.Equal(t, "my-key", k.Name)
	assert.Equal(t, "hash123", k.KeyHash)
	assert.Equal(t, "ev_abc", k.KeyPrefix)
	assert.True(t, k.IsActive)
	assert.NotEqual(t, uuid.Nil, k.ID)
	assert.WithinDuration(t, time.Now(), k.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), k.UpdatedAt, time.Second)
}
