package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/evidence/domain"
)

type EvidenceFilter struct {
	TenantID  uuid.UUID
	Category  domain.Category
	Status    domain.Status
	OwnerID   uuid.UUID
	Expiring  bool // items expiring within the window
	ExpireWin time.Duration
	Expired   bool // already expired items
	Tag       string
	Search    string
	Limit     int
	Offset    int
}

type EvidenceRepository interface {
	Create(ctx context.Context, item *domain.EvidenceItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EvidenceItem, error)
	List(ctx context.Context, filter EvidenceFilter) ([]*domain.EvidenceItem, error)
	Count(ctx context.Context, filter EvidenceFilter) (int, error)
	Update(ctx context.Context, item *domain.EvidenceItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ApprovalRepository interface {
	Create(ctx context.Context, a *domain.Approval) error
	ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Approval, error)
}
