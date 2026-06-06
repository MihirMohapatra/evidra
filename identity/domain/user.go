package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	PasswordHash   string
	Role           Role
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUser(orgID uuid.UUID, email, passwordHash string, role Role) *User {
	now := time.Now()
	return &User{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          email,
		PasswordHash:   passwordHash,
		Role:           role,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
