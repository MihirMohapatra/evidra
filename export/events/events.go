package events

import (
	"github.com/google/uuid"
)

const (
	SubjectExportRequested = "export.requested"
	SubjectExportCompleted = "export.completed"
	SubjectExportFailed    = "export.failed"
)

type ExportRequested struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	Format     string    `json:"format"`
}

func (e ExportRequested) Subject() string { return SubjectExportRequested }

type ExportCompleted struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	Format     string    `json:"format"`
	FileURL    string    `json:"file_url"`
	FileSize   int64     `json:"file_size"`
}

func (e ExportCompleted) Subject() string { return SubjectExportCompleted }

type ExportFailed struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	Format     string    `json:"format"`
	Error      string    `json:"error"`
}

func (e ExportFailed) Subject() string { return SubjectExportFailed }
