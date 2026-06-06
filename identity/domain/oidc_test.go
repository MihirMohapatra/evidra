package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOIDCState(t *testing.T) {
	s := NewOIDCState("google", "state123", "nonce456", 10*time.Minute)
	require.NotNil(t, s)
	assert.Equal(t, "google", s.Provider)
	assert.Equal(t, "state123", s.State)
	assert.Equal(t, "nonce456", s.Nonce)
	assert.True(t, s.ExpiresAt.After(time.Now()))
	assert.False(t, s.ExpiresAt.Before(time.Now().Add(5*time.Minute)))
	assert.NotEqual(t, uuid.Nil, s.ID)
	assert.WithinDuration(t, time.Now(), s.CreatedAt, time.Second)
}

func TestNewLinkedAccount(t *testing.T) {
	uid := uuid.New()
	a := NewLinkedAccount(uid, "google", "sub123", "user@example.com", "Test User")
	require.NotNil(t, a)
	assert.Equal(t, uid, a.UserID)
	assert.Equal(t, "google", a.Provider)
	assert.Equal(t, "sub123", a.Subject)
	assert.Equal(t, "user@example.com", a.Email)
	assert.Equal(t, "Test User", a.Name)
	assert.NotEqual(t, uuid.Nil, a.ID)
	assert.WithinDuration(t, time.Now(), a.CreatedAt, time.Second)
}
