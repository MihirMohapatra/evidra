package dto

import (
	"time"

	"github.com/evidra/evidra/export/domain"
)

type ExportResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	EvidenceID  string `json:"evidence_id"`
	RequesterID string `json:"requester_id"`
	Format      string `json:"format"`
	FileURL     string `json:"file_url,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ToExportResponse(e *domain.Export) ExportResponse {
	return ExportResponse{
		ID:          e.ID.String(),
		TenantID:    e.TenantID.String(),
		EvidenceID:  e.EvidenceID.String(),
		RequesterID: e.RequesterID.String(),
		Format:      string(e.Format),
		FileURL:     e.FileURL,
		FileSize:    e.FileSize,
		Status:      string(e.Status),
		Error:       e.Error,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   e.UpdatedAt.Format(time.RFC3339),
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
