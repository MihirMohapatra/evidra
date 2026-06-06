package events

import (
	"github.com/google/uuid"
	"github.com/evidra/evidra/evidence/domain"
)

const (
	SubjectEvidenceCreated  = "evidence.created"
	SubjectEvidenceUpdated  = "evidence.updated"
	SubjectEvidenceDeleted  = "evidence.deleted"
	SubjectEvidenceApproved = "evidence.approved"
	SubjectEvidenceRejected = "evidence.rejected"
	SubjectEvidenceExported = "evidence.exported"
)

type EvidenceCreated struct {
	ID        uuid.UUID       `json:"id"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	Title     string          `json:"title"`
	Category  domain.Category `json:"category"`
	OwnerID   uuid.UUID       `json:"owner_id"`
	ExpiresAt string          `json:"expires_at"`
}

func (e EvidenceCreated) Subject() string { return SubjectEvidenceCreated }

type EvidenceUpdated struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
}

func (e EvidenceUpdated) Subject() string { return SubjectEvidenceUpdated }

type EvidenceDeleted struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
}

func (e EvidenceDeleted) Subject() string { return SubjectEvidenceDeleted }

type EvidenceStatusChanged struct {
	ID         uuid.UUID    `json:"id"`
	TenantID   uuid.UUID    `json:"tenant_id"`
	Status     domain.Status `json:"status"`
	ReviewerID uuid.UUID    `json:"reviewer_id"`
	Comment    string       `json:"comment"`
}

func (e EvidenceStatusChanged) Subject() string {
	switch e.Status {
	case domain.StatusApproved:
		return SubjectEvidenceApproved
	case domain.StatusDraft:
		return SubjectEvidenceRejected
	case domain.StatusExported:
		return SubjectEvidenceExported
	default:
		return SubjectEvidenceUpdated
	}
}

type EvidenceExported struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Title    string    `json:"title"`
}

func (e EvidenceExported) Subject() string { return SubjectEvidenceExported }
