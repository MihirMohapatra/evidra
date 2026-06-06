package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	auditdomain "github.com/evidra/evidra/audit/domain"
	auditevents "github.com/evidra/evidra/audit/events"
	"github.com/evidra/evidra/compliance/domain"
	"github.com/evidra/evidra/compliance/events"
	"github.com/evidra/evidra/compliance/repository"
	"github.com/evidra/evidra/pkg/queue"
)

type ComplianceService struct {
	frameworks   repository.FrameworkRepository
	controls     repository.ControlRepository
	evMappings   repository.EvidenceMappingRepository
	qMappings    repository.QuestionMappingRepository
	bus          queue.EventBus
	v            *validator.Validate
}

func New(
	frameworks repository.FrameworkRepository,
	controls repository.ControlRepository,
	evMappings repository.EvidenceMappingRepository,
	qMappings repository.QuestionMappingRepository,
	bus queue.EventBus,
) *ComplianceService {
	return &ComplianceService{
		frameworks: frameworks,
		controls:   controls,
		evMappings: evMappings,
		qMappings:  qMappings,
		bus:        bus,
		v:          validator.New(),
	}
}

// --- Frameworks ---

func (s *ComplianceService) CreateFramework(ctx context.Context, name, slug, description, version string) (*domain.Framework, error) {
	if err := s.v.Var(name, "required,max=255"); err != nil {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if err := s.v.Var(slug, "required,alphanum,max=100"); err != nil {
		return nil, fmt.Errorf("%w: slug must be alphanumeric", domain.ErrInvalidInput)
	}

	existing, err := s.frameworks.GetBySlug(ctx, slug)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: framework with slug %q already exists", domain.ErrAlreadyExists, slug)
	}

	f := domain.NewFramework(name, slug, description, version)
	if err := s.frameworks.Create(ctx, f); err != nil {
		return nil, err
	}

	s.publish(ctx, events.FrameworkCreated{ID: f.ID, Name: f.Name, Slug: f.Slug})
	s.auditPublish(ctx, auditdomain.ActionComplianceFrameworkCreated, uuid.Nil, uuid.Nil, f.ID.String())
	return f, nil
}

func (s *ComplianceService) GetFramework(ctx context.Context, id uuid.UUID) (*domain.Framework, error) {
	return s.frameworks.GetByID(ctx, id)
}

func (s *ComplianceService) ListFrameworks(ctx context.Context) ([]*domain.Framework, error) {
	return s.frameworks.List(ctx)
}

func (s *ComplianceService) DeleteFramework(ctx context.Context, id uuid.UUID) error {
	count, err := s.controls.CountByFramework(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%w: framework has %d controls", domain.ErrFrameworkInUse, count)
	}
	return s.frameworks.Delete(ctx, id)
}

// --- Controls ---

type CreateControlInput struct {
	FrameworkID     uuid.UUID `validate:"required"`
	ControlID       string    `validate:"required,max=50"`
	Name            string    `validate:"required,max=255"`
	Description     string
	Category        string `validate:"required"`
	RiskDescription string
}

func (s *ComplianceService) CreateControl(ctx context.Context, input CreateControlInput) (*domain.Control, error) {
	if err := s.v.Struct(input); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	if _, err := s.frameworks.GetByID(ctx, input.FrameworkID); err != nil {
		return nil, err
	}

	cat := domain.ControlCategory(input.Category)
	if cat != domain.ControlCategoryAdministrative && cat != domain.ControlCategoryTechnical && cat != domain.ControlCategoryPhysical {
		return nil, fmt.Errorf("%w: invalid control category", domain.ErrInvalidInput)
	}

	c := domain.NewControl(input.FrameworkID, input.ControlID, input.Name, input.Description, cat, input.RiskDescription)
	if err := s.controls.Create(ctx, c); err != nil {
		return nil, err
	}

	s.publish(ctx, events.ControlCreated{ID: c.ID, FrameworkID: c.FrameworkID, ControlID: c.ControlID, Name: c.Name})
	s.auditPublish(ctx, auditdomain.ActionComplianceControlCreated, uuid.Nil, uuid.Nil, c.ID.String())
	return c, nil
}

func (s *ComplianceService) ListControls(ctx context.Context, frameworkID uuid.UUID) ([]*domain.Control, error) {
	return s.controls.ListByFramework(ctx, frameworkID)
}

// --- Evidence Mapping ---

type MapEvidenceInput struct {
	TenantID    uuid.UUID `validate:"required"`
	EvidenceID  uuid.UUID `validate:"required"`
	ControlID   uuid.UUID `validate:"required"`
	MappedBy    uuid.UUID
	Notes       string
}

func (s *ComplianceService) MapEvidence(ctx context.Context, input MapEvidenceInput) (*domain.EvidenceMapping, error) {
	if err := s.v.Struct(input); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	if _, err := s.controls.GetByID(ctx, input.ControlID); err != nil {
		return nil, err
	}

	existing, err := s.evMappings.GetByEvidenceAndControl(ctx, input.EvidenceID, input.ControlID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: evidence already mapped to this control", domain.ErrMappingExists)
	}

	m := domain.NewEvidenceMapping(input.TenantID, input.EvidenceID, input.ControlID, input.MappedBy, input.Notes)
	if err := s.evMappings.Create(ctx, m); err != nil {
		return nil, err
	}

	s.publish(ctx, events.EvidenceMapped{ID: m.ID, TenantID: m.TenantID, EvidenceID: m.EvidenceID, ControlID: m.ControlID})
	s.auditPublish(ctx, auditdomain.ActionComplianceEvidenceMapped, m.TenantID, m.MappedBy, m.ID.String())
	return m, nil
}

func (s *ComplianceService) UnmapEvidence(ctx context.Context, id uuid.UUID) error {
	m, err := s.evMappings.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.evMappings.Delete(ctx, id); err != nil {
		return err
	}
	s.publish(ctx, events.EvidenceUnmapped{ID: m.ID, TenantID: m.TenantID, EvidenceID: m.EvidenceID, ControlID: m.ControlID})
	return nil
}

func (s *ComplianceService) ListMappingsByControl(ctx context.Context, controlID uuid.UUID) ([]*domain.EvidenceMapping, error) {
	return s.evMappings.ListByControl(ctx, controlID)
}

func (s *ComplianceService) ListMappingsByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.EvidenceMapping, error) {
	return s.evMappings.ListByEvidence(ctx, evidenceID)
}

// --- Coverage ---

func (s *ComplianceService) GetFrameworkCoverage(ctx context.Context, frameworkID uuid.UUID) (*domain.FrameworkCoverage, error) {
	f, err := s.frameworks.GetByID(ctx, frameworkID)
	if err != nil {
		return nil, err
	}

	controls, err := s.controls.ListByFramework(ctx, frameworkID)
	if err != nil {
		return nil, err
	}

	coverage := &domain.FrameworkCoverage{
		Framework: *f,
		Controls:  make([]domain.ControlCoverage, 0, len(controls)),
		Total:     len(controls),
	}

	for _, c := range controls {
		mappings, err := s.evMappings.ListByControl(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		cc := domain.ControlCoverage{
			Control: *c,
			Status:  domain.MappingStatusUnmapped,
		}

		if len(mappings) > 0 {
			cc.Status = domain.MappingStatusMapped
			cc.EvidenceIDs = make([]uuid.UUID, len(mappings))
			for i, m := range mappings {
				cc.EvidenceIDs[i] = m.EvidenceID
			}
			coverage.Mapped++
		}

		coverage.Controls = append(coverage.Controls, cc)
	}

	return coverage, nil
}

// --- publish ---

func (s *ComplianceService) publish(ctx context.Context, event queue.Event) {
	if s.bus == nil {
		return
	}
	if err := s.bus.Publish(ctx, event); err != nil {
		slog.Warn("failed to publish event", "subject", event.Subject(), "error", err)
	}
}

func (s *ComplianceService) auditPublish(ctx context.Context, action auditdomain.Action, tenantID, actorID uuid.UUID, targetID string) {
	s.publish(ctx, auditevents.NewAuditRecorded(tenantID, actorID, string(action), targetID))
}
