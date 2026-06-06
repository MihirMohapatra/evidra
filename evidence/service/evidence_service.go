package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/evidra/evidra/evidence/domain"
	"github.com/evidra/evidra/evidence/events"
	"github.com/evidra/evidra/evidence/repository"
	"github.com/evidra/evidra/pkg/queue"
)

type EvidenceService struct {
	evidence  repository.EvidenceRepository
	approvals repository.ApprovalRepository
	bus       queue.EventBus
	v         *validator.Validate
}

func New(
	evidence repository.EvidenceRepository,
	approvals repository.ApprovalRepository,
	bus queue.EventBus,
) *EvidenceService {
	return &EvidenceService{
		evidence:  evidence,
		approvals: approvals,
		bus:       bus,
		v:         validator.New(),
	}
}

func (s *EvidenceService) Create(ctx context.Context, input domain.CreateEvidenceInput) (*domain.EvidenceItem, error) {
	if err := s.v.Var(input.Title, "required,max=500"); err != nil {
		return nil, fmt.Errorf("%w: title is required (max 500 chars)", domain.ErrInvalidInput)
	}
	if !input.Category.Valid() {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidCategory, input.Category)
	}

	item := domain.NewEvidence(input)
	if err := s.evidence.Create(ctx, item); err != nil {
		return nil, err
	}

	s.publish(ctx, events.EvidenceCreated{
		ID:        item.ID,
		TenantID:  item.TenantID,
		Title:     item.Title,
		Category:  item.Category,
		OwnerID:   item.OwnerID,
		ExpiresAt: item.ExpiresAt.Format(time.RFC3339),
	})

	slog.Info("evidence created",
		"id", item.ID,
		"category", item.Category,
		"tenant", item.TenantID,
	)
	return item, nil
}

func (s *EvidenceService) Get(ctx context.Context, id uuid.UUID) (*domain.EvidenceItem, error) {
	return s.evidence.GetByID(ctx, id)
}

type ListFilter struct {
	TenantID uuid.UUID
	Category string
	Status   string
	OwnerID  uuid.UUID
	Expiring bool
	Expired  bool
	Tag      string
	Search   string
	Limit    int
	Offset   int
}

func (s *EvidenceService) List(ctx context.Context, filter ListFilter) ([]*domain.EvidenceItem, int, error) {
	repoFilter := repository.EvidenceFilter{
		TenantID: filter.TenantID,
		OwnerID:  filter.OwnerID,
		Tag:      filter.Tag,
		Search:   filter.Search,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}

	if filter.Category != "" {
		cat, ok := domain.ParseCategory(filter.Category)
		if !ok {
			return nil, 0, fmt.Errorf("%w: %s", domain.ErrInvalidCategory, filter.Category)
		}
		repoFilter.Category = cat
	}
	if filter.Status != "" {
		repoFilter.Status = domain.Status(filter.Status)
		if !repoFilter.Status.Valid() {
			return nil, 0, fmt.Errorf("%w: %s", domain.ErrInvalidTransition, filter.Status)
		}
	}
	if filter.Expiring {
		repoFilter.Expiring = true
		repoFilter.ExpireWin = 30 * 24 * time.Hour
	}
	repoFilter.Expired = filter.Expired

	items, err := s.evidence.List(ctx, repoFilter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.evidence.Count(ctx, repoFilter)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (s *EvidenceService) Update(ctx context.Context, id uuid.UUID, input domain.CreateEvidenceInput) (*domain.EvidenceItem, error) {
	item, err := s.evidence.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		item.Title = input.Title
	}
	if input.Content != "" {
		item.Content = input.Content
	}
	if input.Category.Valid() {
		item.Category = input.Category
	}
	if input.OwnerID != uuid.Nil {
		item.OwnerID = input.OwnerID
	}
	if input.SourceURL != "" {
		item.SourceURL = input.SourceURL
	}
	if input.Tags != nil {
		item.Tags = input.Tags
	}
	if !input.ExpiresAt.IsZero() {
		item.ExpiresAt = input.ExpiresAt
	}
	item.Version++
	item.UpdatedAt = time.Now()

	if err := s.evidence.Update(ctx, item); err != nil {
		return nil, err
	}

	s.publish(ctx, events.EvidenceUpdated{
		ID:       item.ID,
		TenantID: item.TenantID,
	})

	return item, nil
}

func (s *EvidenceService) Delete(ctx context.Context, id uuid.UUID) error {
	item, err := s.evidence.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.evidence.Delete(ctx, id); err != nil {
		return err
	}
	s.publish(ctx, events.EvidenceDeleted{
		ID:       item.ID,
		TenantID: item.TenantID,
	})
	return nil
}

func (s *EvidenceService) Submit(ctx context.Context, id uuid.UUID) (*domain.EvidenceItem, error) {
	item, err := s.evidence.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := item.Submit(); err != nil {
		return nil, err
	}
	if err := s.evidence.Update(ctx, item); err != nil {
		return nil, err
	}
	s.publish(ctx, events.EvidenceStatusChanged{
		ID:       item.ID,
		TenantID: item.TenantID,
		Status:   item.Status,
	})
	return item, nil
}

func (s *EvidenceService) Approve(ctx context.Context, id, reviewerID uuid.UUID, comment string) (*domain.EvidenceItem, error) {
	return s.statusChange(ctx, id, reviewerID, comment, func(item *domain.EvidenceItem) error {
		return item.Approve()
	})
}

func (s *EvidenceService) Reject(ctx context.Context, id, reviewerID uuid.UUID, comment string) (*domain.EvidenceItem, error) {
	return s.statusChange(ctx, id, reviewerID, comment, func(item *domain.EvidenceItem) error {
		return item.Reject()
	})
}

func (s *EvidenceService) Export(ctx context.Context, id uuid.UUID) (*domain.EvidenceItem, error) {
	item, err := s.evidence.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := item.Export(); err != nil {
		return nil, err
	}
	if err := s.evidence.Update(ctx, item); err != nil {
		return nil, err
	}
	s.publish(ctx, events.EvidenceExported{
		ID:       item.ID,
		TenantID: item.TenantID,
		Title:    item.Title,
	})
	slog.Info("evidence exported", "id", item.ID, "tenant", item.TenantID)
	return item, nil
}

func (s *EvidenceService) GetApprovalHistory(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Approval, error) {
	return s.approvals.ListByEvidence(ctx, evidenceID)
}

// --- internal ---

func (s *EvidenceService) statusChange(ctx context.Context, id, reviewerID uuid.UUID, comment string, transition func(*domain.EvidenceItem) error) (*domain.EvidenceItem, error) {
	item, err := s.evidence.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := transition(item); err != nil {
		return nil, err
	}
	if err := s.evidence.Update(ctx, item); err != nil {
		return nil, err
	}

	approval := domain.NewApproval(id, reviewerID, item.Status, comment)
	if err := s.approvals.Create(ctx, approval); err != nil {
		slog.Error("failed to record approval", "evidence_id", id, "error", err)
	}

	s.publish(ctx, events.EvidenceStatusChanged{
		ID:         item.ID,
		TenantID:   item.TenantID,
		Status:     item.Status,
		ReviewerID: reviewerID,
		Comment:    comment,
	})

	return item, nil
}

func (s *EvidenceService) publish(ctx context.Context, event queue.Event) {
	if s.bus == nil {
		return
	}
	if err := s.bus.Publish(ctx, event); err != nil {
		slog.Warn("failed to publish event", "subject", event.Subject(), "error", err)
	}
}
