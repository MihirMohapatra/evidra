package domain

import (
	"time"

	"github.com/google/uuid"
)

type Format string

const (
	FormatPDF  Format = "PDF"
	FormatXLSX Format = "XLSX"
	FormatDOCX Format = "DOCX"
)

func (f Format) Valid() bool {
	switch f {
	case FormatPDF, FormatXLSX, FormatDOCX:
		return true
	}
	return false
}

func (f Format) ContentType() string {
	switch f {
	case FormatPDF:
		return "application/pdf"
	case FormatXLSX:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case FormatDOCX:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	return "application/octet-stream"
}

func (f Format) Extension() string {
	return string(f)
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Export struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	EvidenceID  uuid.UUID
	RequesterID uuid.UUID
	Format      Format
	FileURL     string
	FileSize    int64
	Status      Status
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewExport(tenantID, evidenceID, requesterID uuid.UUID, format Format) *Export {
	return &Export{
		ID:          uuid.New(),
		TenantID:    tenantID,
		EvidenceID:  evidenceID,
		RequesterID: requesterID,
		Format:      format,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
