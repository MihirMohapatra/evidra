package domain

import (
	"time"

	"github.com/google/uuid"
)

type OIDCState struct {
	ID         uuid.UUID
	Provider   string
	State      string
	Nonce      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

func NewOIDCState(provider, state, nonce string, ttl time.Duration) *OIDCState {
	now := time.Now()
	return &OIDCState{
		ID:        uuid.New(),
		Provider:  provider,
		State:     state,
		Nonce:     nonce,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}
}

type LinkedAccount struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Provider    string
	Subject     string
	Email       string
	Name        string
	CreatedAt   time.Time
}

func NewLinkedAccount(userID uuid.UUID, provider, subject, email, name string) *LinkedAccount {
	return &LinkedAccount{
		ID:        uuid.New(),
		UserID:    userID,
		Provider:  provider,
		Subject:   subject,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}
}
