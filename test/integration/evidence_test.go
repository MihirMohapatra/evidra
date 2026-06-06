package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/evidra/evidra/evidence/domain"
	evpg "github.com/evidra/evidra/evidence/repository/postgres"
	evservice "github.com/evidra/evidra/evidence/service"
	"github.com/evidra/evidra/pkg/queue"
)

type noopEventBus struct{}

func (n *noopEventBus) Publish(_ context.Context, _ queue.Event) error {
	return nil
}

func (n *noopEventBus) Subscribe(_ string, _ func(ctx context.Context, event []byte) error) error {
	return nil
}

func (n *noopEventBus) Close() error {
	return nil
}

type EvidenceSuite struct {
	suite.Suite
	ctx     context.Context
	service *evservice.EvidenceService
	evRepo  *evpg.EvidenceRepo
	appRepo *evpg.ApprovalRepo
}

func TestEvidenceSuite(t *testing.T) {
	if evidencePool == nil {
		t.Skip("evidence postgres container not available")
	}
	suite.Run(t, new(EvidenceSuite))
}

func (s *EvidenceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.evRepo = evpg.NewEvidenceRepo(evidencePool)
	s.appRepo = evpg.NewApprovalRepo(evidencePool)

	s.service = evservice.New(
		s.evRepo,
		s.appRepo,
		&noopEventBus{},
	)
}

func (s *EvidenceSuite) TearDownTest() {
	_, _ = evidencePool.Exec(s.ctx, "DELETE FROM approvals")
	_, _ = evidencePool.Exec(s.ctx, "DELETE FROM evidence_items")
}

func (s *EvidenceSuite) TestCreateEvidence() {
	tenantID := uuid.New()
	ownerID := uuid.New()

	input := domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     "Test Evidence Item",
		Content:   "This is the content of the evidence",
		Category:  domain.CategoryPolicy,
		OwnerID:   ownerID,
		SourceURL: "https://example.com/evidence",
		Tags:      []string{"compliance", "security"},
		ExpiresAt: time.Now().Add(365 * 24 * time.Hour),
	}

	item, err := s.service.Create(s.ctx, input)
	s.Require().NoError(err)
	s.Require().NotNil(item)
	s.Equal(input.Title, item.Title)
	s.Equal(input.Content, item.Content)
	s.Equal(domain.CategoryPolicy, item.Category)
	s.Equal(domain.StatusDraft, item.Status)
	s.Equal(ownerID, item.OwnerID)
	s.Equal(1, item.Version)
	s.NotEqual(uuid.Nil, item.ID)

	got, err := s.service.Get(s.ctx, item.ID)
	s.Require().NoError(err)
	s.Equal(item.Title, got.Title)
	s.Equal(item.Status, got.Status)
}

func (s *EvidenceSuite) TestApprovalWorkflow() {
	tenantID := uuid.New()
	ownerID := uuid.New()
	reviewerID := uuid.New()

	input := domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     "Approval Workflow Test",
		Content:   "Testing the full approval workflow",
		Category:  domain.CategoryClaim,
		OwnerID:   ownerID,
		Tags:      []string{"workflow"},
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	item, err := s.service.Create(s.ctx, input)
	s.Require().NoError(err)
	s.Equal(domain.StatusDraft, item.Status)

	submitted, err := s.service.Submit(s.ctx, item.ID)
	s.Require().NoError(err)
	s.Equal(domain.StatusReview, submitted.Status)

	approved, err := s.service.Approve(s.ctx, item.ID, reviewerID, "Looks good, approved")
	s.Require().NoError(err)
	s.Equal(domain.StatusApproved, approved.Status)

	history, err := s.service.GetApprovalHistory(s.ctx, item.ID)
	s.Require().NoError(err)
	s.Require().Len(history, 1)
	s.Equal(domain.StatusApproved, history[0].Status)
	s.Equal(reviewerID, history[0].ReviewerID)
	s.Equal("Looks good, approved", history[0].Comment)

	s.Equal(1, item.Version)
}

func (s *EvidenceSuite) TestApprovalWorkflowReject() {
	tenantID := uuid.New()
	ownerID := uuid.New()
	reviewerID := uuid.New()

	input := domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     "Rejection Workflow Test",
		Content:   "Testing the rejection workflow",
		Category:  domain.CategoryArchitecture,
		OwnerID:   ownerID,
	}

	item, err := s.service.Create(s.ctx, input)
	s.Require().NoError(err)
	s.Equal(domain.StatusDraft, item.Status)

	submitted, err := s.service.Submit(s.ctx, item.ID)
	s.Require().NoError(err)
	s.Equal(domain.StatusReview, submitted.Status)

	rejected, err := s.service.Reject(s.ctx, item.ID, reviewerID, "Needs more details")
	s.Require().NoError(err)
	s.Equal(domain.StatusDraft, rejected.Status)

	history, err := s.service.GetApprovalHistory(s.ctx, item.ID)
	s.Require().NoError(err)
	s.Require().Len(history, 1)
	s.Equal(domain.StatusDraft, history[0].Status)
	s.Equal(reviewerID, history[0].ReviewerID)
	s.Equal("Needs more details", history[0].Comment)
}

func (s *EvidenceSuite) TestInvalidTransition() {
	tenantID := uuid.New()
	ownerID := uuid.New()

	input := domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     "Invalid Transition Test",
		Content:   "Testing invalid transitions",
		Category:  domain.CategoryCertification,
		OwnerID:   ownerID,
	}

	item, err := s.service.Create(s.ctx, input)
	s.Require().NoError(err)
	s.Equal(domain.StatusDraft, item.Status)

	_, err = s.service.Approve(s.ctx, item.ID, uuid.New(), "")
	s.Error(err)
	s.ErrorContains(err, "invalid transition")

	_, err = s.service.Submit(s.ctx, item.ID)
	s.Require().NoError(err)

	_, err = s.service.Submit(s.ctx, item.ID)
	s.Error(err)
	s.ErrorContains(err, "invalid transition")
}

func (s *EvidenceSuite) TestListAndCount() {
	tenantID := uuid.New()
	ownerID := uuid.New()

	for i := 0; i < 3; i++ {
		input := domain.CreateEvidenceInput{
			TenantID:  tenantID,
			Title:     "List Test Item",
			Content:   "Content for listing test",
			Category:  domain.CategoryAnswer,
			OwnerID:   ownerID,
		}
		_, err := s.service.Create(s.ctx, input)
		s.Require().NoError(err)
	}

	filter := evservice.ListFilter{
		TenantID: tenantID,
		Limit:    10,
	}
	items, count, err := s.service.List(s.ctx, filter)
	s.Require().NoError(err)
	s.Equal(3, count)
	s.Len(items, 3)
}
