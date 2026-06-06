package service

import (
	"context"

	"github.com/evidra/evidra/orchestrator/domain"
)

type AIEngine interface {
	GenerateAnswer(ctx context.Context, question domain.Question, evidence []domain.Evidence) (domain.Draft, error)
	Name() string
}
