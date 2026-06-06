package domain

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uuid.UUID
	Text      string
	Context   string
	TenantID  uuid.UUID
	CreatedAt time.Time
}

func NewQuestion(text string, tenantID uuid.UUID) *Question {
	return &Question{
		ID:        uuid.New(),
		Text:      text,
		TenantID:  tenantID,
		CreatedAt: time.Now(),
	}
}
