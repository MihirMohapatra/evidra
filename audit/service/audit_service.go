package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/evidra/evidra/audit/domain"
	"github.com/evidra/evidra/audit/repository"
)

type AuditService struct {
	repo repository.AuditRepository
}

func New(repo repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

type RecordInput struct {
	TenantID uuid.UUID
	ActorID  uuid.UUID
	Action   domain.Action
	TargetID string
	Metadata map[string]any
}

func (s *AuditService) Record(ctx context.Context, input RecordInput) (*domain.AuditEvent, error) {
	if input.Action == "" {
		return nil, fmt.Errorf("%w: action is required", domain.ErrInvalidInput)
	}
	if input.TenantID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id is required", domain.ErrInvalidInput)
	}

	event := domain.NewAuditEvent(input.TenantID, input.ActorID, input.Action, input.TargetID, input.Metadata)

	if err := s.repo.Append(ctx, event); err != nil {
		return nil, fmt.Errorf("record audit: %w", err)
	}

	slog.Info("audit event recorded",
		"action", event.Action,
		"tenant", event.TenantID,
		"actor", event.ActorID,
	)

	return event, nil
}

func (s *AuditService) List(ctx context.Context, filter repository.AuditFilter) ([]*domain.AuditEvent, int, error) {
	events, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}
