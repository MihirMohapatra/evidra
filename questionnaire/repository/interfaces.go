package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/evidra/evidra/questionnaire/domain"
)

type QuestionnaireRepository interface {
	Create(ctx context.Context, q *domain.Questionnaire) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Questionnaire, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Questionnaire, error)
	Update(ctx context.Context, q *domain.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type QuestionRepository interface {
	Create(ctx context.Context, q *domain.Question) error
	CreateBatch(ctx context.Context, questions []*domain.Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	ListByQuestionnaire(ctx context.Context, questionnaireID uuid.UUID) ([]*domain.Question, error)
	Update(ctx context.Context, q *domain.Question) error
	DeleteByQuestionnaire(ctx context.Context, questionnaireID uuid.UUID) error
}
