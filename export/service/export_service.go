package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	auditdomain "github.com/evidra/evidra/audit/domain"
	auditevents "github.com/evidra/evidra/audit/events"
	evdomain "github.com/evidra/evidra/evidence/domain"
	evrepo "github.com/evidra/evidra/evidence/repository"
	"github.com/evidra/evidra/export/domain"
	"github.com/evidra/evidra/export/events"
	"github.com/evidra/evidra/export/repository"
	"github.com/evidra/evidra/pkg/queue"
	"github.com/evidra/evidra/pkg/storage"
)

type ExportService struct {
	exports  repository.ExportRepository
	evidence evrepo.EvidenceRepository
	store    storage.FileStorage
	bus      queue.EventBus
	v        *validator.Validate
	bucket   string
}

func New(
	exports repository.ExportRepository,
	evidence evrepo.EvidenceRepository,
	store storage.FileStorage,
	bus queue.EventBus,
	bucket string,
) *ExportService {
	return &ExportService{
		exports:  exports,
		evidence: evidence,
		store:    store,
		bus:      bus,
		v:        validator.New(),
		bucket:   bucket,
	}
}

type ExportInput struct {
	TenantID    uuid.UUID `validate:"required"`
	EvidenceID  uuid.UUID `validate:"required"`
	RequesterID uuid.UUID `validate:"required"`
	Format      string    `validate:"required"`
}

func (s *ExportService) Export(ctx context.Context, input ExportInput) (*domain.Export, error) {
	if err := s.v.Struct(input); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	fmtVal := domain.Format(input.Format)
	if !fmtVal.Valid() {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidFormat, input.Format)
	}

	ev, err := s.evidence.GetByID(ctx, input.EvidenceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrEvidenceNotFound, err)
	}

	exp := domain.NewExport(input.TenantID, input.EvidenceID, input.RequesterID, fmtVal)
	exp.Status = domain.StatusProcessing
	if err := s.exports.Create(ctx, exp); err != nil {
		return nil, err
	}

	s.publish(ctx, events.ExportRequested{
		ID: exp.ID, TenantID: exp.TenantID, EvidenceID: exp.EvidenceID, Format: string(exp.Format),
	})

	data, err := s.generate(ctx, ev, fmtVal)
	if err != nil {
		exp.Status = domain.StatusFailed
		exp.Error = err.Error()
		_ = s.exports.Update(ctx, exp)
		s.publish(ctx, events.ExportFailed{ID: exp.ID, TenantID: exp.TenantID, EvidenceID: exp.EvidenceID, Format: string(exp.Format), Error: err.Error()})
		return nil, fmt.Errorf("%w: %v", domain.ErrExportFailed, err)
	}

	key := fmt.Sprintf("exports/%s/%s/%s", exp.TenantID, exp.ID, fmtVal.Extension())
	if err := s.store.Upload(ctx, s.bucket, key, bytes.NewReader(data), int64(len(data)), fmtVal.ContentType()); err != nil {
		exp.Status = domain.StatusFailed
		exp.Error = err.Error()
		_ = s.exports.Update(ctx, exp)
		return nil, fmt.Errorf("storage upload: %w", err)
	}

	exp.FileURL = key
	exp.FileSize = int64(len(data))
	exp.Status = domain.StatusCompleted
	exp.UpdatedAt = time.Now()
	if err := s.exports.Update(ctx, exp); err != nil {
		return nil, err
	}

	s.publish(ctx, events.ExportCompleted{
		ID: exp.ID, TenantID: exp.TenantID, EvidenceID: exp.EvidenceID,
		Format: string(exp.Format), FileURL: exp.FileURL, FileSize: exp.FileSize,
	})
	s.auditPublish(ctx, auditdomain.ActionDocumentExported, exp.TenantID, exp.RequesterID, exp.ID.String())

	return exp, nil
}

func (s *ExportService) Get(ctx context.Context, id uuid.UUID) (*domain.Export, error) {
	return s.exports.GetByID(ctx, id)
}

func (s *ExportService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Export, error) {
	return s.exports.ListByTenant(ctx, tenantID)
}

func (s *ExportService) ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Export, error) {
	return s.exports.ListByEvidence(ctx, evidenceID)
}

func (s *ExportService) generate(ctx context.Context, ev *evdomain.EvidenceItem, format domain.Format) ([]byte, error) {
	switch format {
	case domain.FormatPDF:
		return s.generatePDF(ev)
	case domain.FormatXLSX:
		return s.generateXLSX(ev)
	case domain.FormatDOCX:
		return s.generateDOCX(ev)
	}
	return nil, fmt.Errorf("unsupported format: %s", format)
}

func (s *ExportService) generatePDF(ev *evdomain.EvidenceItem) ([]byte, error) {
	pdf := NewPDFGenerator()
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, ev.Title)
	pdf.Ln(12)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, "Category: "+string(ev.Category))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Status: "+string(ev.Status))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Created: "+ev.CreatedAt.Format(time.RFC3339))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Expires: "+ev.ExpiresAt.Format(time.RFC3339))
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 11)
	pdf.MultiCell(0, 5, ev.Content, "", "L", false)
	return PDFBytes(pdf)
}

func (s *ExportService) generateXLSX(ev *evdomain.EvidenceItem) ([]byte, error) {
	f := NewXLSXGenerator()
	_ = f.SetCellValue("Sheet1", "A1", "Title")
	_ = f.SetCellValue("Sheet1", "B1", ev.Title)
	_ = f.SetCellValue("Sheet1", "A2", "Category")
	_ = f.SetCellValue("Sheet1", "B2", string(ev.Category))
	_ = f.SetCellValue("Sheet1", "A3", "Status")
	_ = f.SetCellValue("Sheet1", "B3", string(ev.Status))
	_ = f.SetCellValue("Sheet1", "A4", "Content")
	_ = f.SetCellValue("Sheet1", "B4", ev.Content)
	if ev.SourceURL != "" {
		_ = f.SetCellValue("Sheet1", "A5", "Source URL")
		_ = f.SetCellValue("Sheet1", "B5", ev.SourceURL)
	}
	return XLSXBytes(f)
}

func (s *ExportService) generateDOCX(ev *evdomain.EvidenceItem) ([]byte, error) {
	doc := NewDOCXGenerator()
	doc.AddTitle(ev.Title)
	doc.AddParagraph("Category: " + string(ev.Category))
	doc.AddParagraph("Status: " + string(ev.Status))
	doc.AddParagraph("Created: " + ev.CreatedAt.Format(time.RFC3339))
	doc.AddParagraph("Expires: " + ev.ExpiresAt.Format(time.RFC3339))
	doc.AddParagraph("")
	doc.AddParagraph(ev.Content)
	return doc.Bytes()
}

func (s *ExportService) publish(ctx context.Context, event queue.Event) {
	if s.bus == nil {
		return
	}
	if err := s.bus.Publish(ctx, event); err != nil {
		slog.Warn("failed to publish event", "subject", event.Subject(), "error", err)
	}
}

func (s *ExportService) auditPublish(ctx context.Context, action auditdomain.Action, tenantID, actorID uuid.UUID, targetID string) {
	s.publish(ctx, auditevents.NewAuditRecorded(tenantID, actorID, string(action), targetID))
}
