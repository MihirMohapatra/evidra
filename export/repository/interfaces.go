package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/evidra/evidra/export/domain"
)

type ExportRepository interface {
	Create(ctx context.Context, exp *domain.Export) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Export, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Export, error)
	ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Export, error)
	Update(ctx context.Context, exp *domain.Export) error
}
