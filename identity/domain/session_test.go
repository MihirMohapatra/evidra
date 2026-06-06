package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	uid := uuid.New()
	s := NewSession(uid, "token123", "refresh123", time.Hour)
	require.NotNil(t, s)
	assert.Equal(t, uid, s.UserID)
	assert.Equal(t, "token123", s.Token)
	assert.Equal(t, "refresh123", s.RefreshToken)
	assert.True(t, s.ExpiresAt.After(time.Now()))
	assert.WithinDuration(t, time.Now().Add(time.Hour), s.ExpiresAt, time.Second)
	assert.NotEqual(t, uuid.Nil, s.ID)
}

func TestSession_IsExpired(t *testing.T) {
	s := NewSession(uuid.New(), "token", "refresh", -time.Hour)
	assert.True(t, s.IsExpired())

	s = NewSession(uuid.New(), "token", "refresh", time.Hour)
	assert.False(t, s.IsExpired())
}
