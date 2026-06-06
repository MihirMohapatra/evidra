package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	auditdomain "github.com/evidra/evidra/audit/domain"
	auditevents "github.com/evidra/evidra/audit/events"
	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/events"
	"github.com/evidra/evidra/orchestrator/repository"
	"github.com/evidra/evidra/pkg/queue"
)

type OrchestratorService struct {
	embedder    Embedder
	engine      AIEngine
	embeddings  repository.EmbeddingRepository
	drafts      repository.DraftRepository
	validator   *Validator
	bus         queue.EventBus
	v           *validator.Validate
	topK        int
}

func New(
	embedder Embedder,
	engine AIEngine,
	embeddings repository.EmbeddingRepository,
	drafts repository.DraftRepository,
	bus queue.EventBus,
	topK int,
) *OrchestratorService {
	if topK <= 0 {
		topK = 5
	}
	return &OrchestratorService{
		embedder:   embedder,
		engine:     engine,
		embeddings: embeddings,
		drafts:     drafts,
		validator:  NewValidator(),
		bus:        bus,
		v:          validator.New(),
		topK:       topK,
	}
}

type AnswerRequest struct {
	Question string    `validate:"required"`
	Context  string
	TenantID uuid.UUID
}

type AnswerResult struct {
	Draft     *domain.Draft
	Evidence  []domain.Evidence
}

func (s *OrchestratorService) Answer(ctx context.Context, req AnswerRequest) (*AnswerResult, error) {
	if err := s.v.Var(req.Question, "required"); err != nil {
		return nil, fmt.Errorf("%w: question is required", domain.ErrInvalidInput)
	}

	start := time.Now()
	slog.Info("processing question", "question", truncate(req.Question, 100))

	question := domain.NewQuestion(req.Question, req.TenantID)
	question.Context = req.Context

	embedding, err := s.embedder.GenerateEmbedding(ctx, req.Question)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrEmbeddingFailed, err)
	}
	slog.Info("embedding generated", "dims", len(embedding), "elapsed", time.Since(start))

	evidence, err := s.embeddings.SearchSimilar(ctx, req.TenantID, embedding, s.topK)
	if err != nil {
		return nil, fmt.Errorf("similarity search: %w", err)
	}
	slog.Info("evidence retrieved", "count", len(evidence), "elapsed", time.Since(start))

	draft, err := s.engine.GenerateAnswer(ctx, *question, evidence)
	if err != nil {
		return nil, err
	}
	slog.Info("answer generated", "model", draft.ModelUsed, "confidence", draft.Confidence, "elapsed", time.Since(start))

	if err := s.validator.ValidateDraft(&draft); err != nil {
		return nil, err
	}

	if err := s.drafts.Create(ctx, &draft); err != nil {
		return nil, fmt.Errorf("save draft: %w", err)
	}
	slog.Info("draft saved", "id", draft.ID, "elapsed", time.Since(start))

	s.publish(ctx, events.DraftCreated{
		ID:           draft.ID,
		QuestionID:   draft.QuestionID,
		QuestionText: draft.QuestionText,
		ModelUsed:    draft.ModelUsed,
		Confidence:   draft.Confidence,
		EvidenceIDs:  draft.EvidenceIDs,
	})
	s.auditPublish(ctx, auditdomain.ActionAIGenerated, req.TenantID, uuid.Nil, draft.ID.String())

	return &AnswerResult{Draft: &draft, Evidence: evidence}, nil
}

func (s *OrchestratorService) GetDraft(ctx context.Context, id uuid.UUID) (*domain.Draft, error) {
	return s.drafts.GetByID(ctx, id)
}

func (s *OrchestratorService) ListDrafts(ctx context.Context, limit, offset int) ([]*domain.Draft, error) {
	return s.drafts.List(ctx, limit, offset)
}

func (s *OrchestratorService) ApproveDraft(ctx context.Context, id uuid.UUID) error {
	d, err := s.drafts.GetByID(ctx, id)
	if err != nil {
		return err
	}
	d.Status = domain.DraftApproved
	d.UpdatedAt = time.Now()
	if err := s.drafts.UpdateStatus(ctx, id, domain.DraftApproved, d.Feedback); err != nil {
		return err
	}
	s.publish(ctx, events.DraftStatusChanged{ID: id, Status: string(domain.DraftApproved)})
	s.auditPublish(ctx, auditdomain.ActionDraftApproved, uuid.Nil, uuid.Nil, id.String())
	return nil
}

func (s *OrchestratorService) RejectDraft(ctx context.Context, id uuid.UUID, feedback string) error {
	if _, err := s.drafts.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.drafts.UpdateStatus(ctx, id, domain.DraftRejected, feedback); err != nil {
		return err
	}
	s.publish(ctx, events.DraftStatusChanged{ID: id, Status: string(domain.DraftRejected), Feedback: feedback})
	s.auditPublish(ctx, auditdomain.ActionDraftRejected, uuid.Nil, uuid.Nil, id.String())
	return nil
}

func (s *OrchestratorService) auditPublish(ctx context.Context, action auditdomain.Action, tenantID, actorID uuid.UUID, targetID string) {
	s.publish(ctx, auditevents.NewAuditRecorded(tenantID, actorID, string(action), targetID))
}

func (s *OrchestratorService) publish(ctx context.Context, event queue.Event) {
	if s.bus == nil {
		return
	}
	if err := s.bus.Publish(ctx, event); err != nil {
		slog.Warn("failed to publish event", "subject", event.Subject(), "error", err)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
