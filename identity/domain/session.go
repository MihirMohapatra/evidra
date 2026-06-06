package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Token        string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSession(userID uuid.UUID, token, refreshToken string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:           uuid.New(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(ttl),
		CreatedAt:    now,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
