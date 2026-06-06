package domain

import (
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	KeyHash        string
	KeyPrefix      string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewAPIKey(orgID uuid.UUID, name, keyHash, keyPrefix string) *APIKey {
	now := time.Now()
	return &APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		KeyHash:        keyHash,
		KeyPrefix:      keyPrefix,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
