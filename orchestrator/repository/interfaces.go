package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/evidra/evidra/orchestrator/domain"
)

type EmbeddingRepository interface {
	SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]domain.Evidence, error)
	Upsert(ctx context.Context, chunk *domain.EvidenceChunk) error
}

type DraftRepository interface {
	Create(ctx context.Context, draft *domain.Draft) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Draft, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Draft, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DraftStatus, feedback string) error
}
