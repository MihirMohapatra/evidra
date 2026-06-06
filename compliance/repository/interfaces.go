package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/evidra/evidra/compliance/domain"
)

type FrameworkRepository interface {
	Create(ctx context.Context, f *domain.Framework) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Framework, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Framework, error)
	List(ctx context.Context) ([]*domain.Framework, error)
	Update(ctx context.Context, f *domain.Framework) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ControlRepository interface {
	Create(ctx context.Context, c *domain.Control) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Control, error)
	ListByFramework(ctx context.Context, frameworkID uuid.UUID) ([]*domain.Control, error)
	Update(ctx context.Context, c *domain.Control) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByFramework(ctx context.Context, frameworkID uuid.UUID) (int, error)
}

type EvidenceMappingRepository interface {
	Create(ctx context.Context, m *domain.EvidenceMapping) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EvidenceMapping, error)
	ListByControl(ctx context.Context, controlID uuid.UUID) ([]*domain.EvidenceMapping, error)
	ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.EvidenceMapping, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.EvidenceMapping, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEvidenceAndControl(ctx context.Context, evidenceID, controlID uuid.UUID) (*domain.EvidenceMapping, error)
}

type QuestionMappingRepository interface {
	Create(ctx context.Context, m *domain.QuestionMapping) error
	ListByControl(ctx context.Context, controlID uuid.UUID) ([]*domain.QuestionMapping, error)
	ListByQuestion(ctx context.Context, questionID uuid.UUID) ([]*domain.QuestionMapping, error)
}
