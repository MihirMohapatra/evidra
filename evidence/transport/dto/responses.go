package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/evidence/domain"
)

type EvidenceResponse struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	OwnerID   uuid.UUID `json:"owner_id"`
	SourceURL string    `json:"source_url"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	ExpiresAt *string   `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToEvidenceResponse(item *domain.EvidenceItem) EvidenceResponse {
	resp := EvidenceResponse{
		ID:        item.ID,
		TenantID:  item.TenantID,
		Title:     item.Title,
		Content:   item.Content,
		Category:  string(item.Category),
		Status:    string(item.Status),
		OwnerID:   item.OwnerID,
		SourceURL: item.SourceURL,
		Tags:      item.Tags,
		Version:   item.Version,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	if !item.ExpiresAt.IsZero() {
		s := item.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &s
	}
	return resp
}

type ApprovalResponse struct {
	ID         uuid.UUID `json:"id"`
	EvidenceID uuid.UUID `json:"evidence_id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	Status     string    `json:"status"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToApprovalResponse(a *domain.Approval) ApprovalResponse {
	return ApprovalResponse{
		ID:         a.ID,
		EvidenceID: a.EvidenceID,
		ReviewerID: a.ReviewerID,
		Status:     string(a.Status),
		Comment:    a.Comment,
		CreatedAt:  a.CreatedAt,
	}
}

type PaginatedResponse struct {
	Items      any `json:"items"`
	Total      int `json:"total"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
