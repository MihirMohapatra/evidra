package domain

import (
	"time"

	"github.com/google/uuid"
)

type EvidenceItem struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Title     string
	Content   string
	Category  Category
	Status    Status
	OwnerID   uuid.UUID
	SourceURL string
	Tags      []string
	Version   int
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateEvidenceInput struct {
	TenantID  uuid.UUID
	Title     string
	Content   string
	Category  Category
	OwnerID   uuid.UUID
	SourceURL string
	Tags      []string
	ExpiresAt time.Time
}

func NewEvidence(input CreateEvidenceInput) *EvidenceItem {
	now := time.Now()
	if input.Tags == nil {
		input.Tags = []string{}
	}
	return &EvidenceItem{
		ID:        uuid.New(),
		TenantID:  input.TenantID,
		Title:     input.Title,
		Content:   input.Content,
		Category:  input.Category,
		Status:    StatusDraft,
		OwnerID:   input.OwnerID,
		SourceURL: input.SourceURL,
		Tags:      input.Tags,
		Version:   1,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (e *EvidenceItem) Submit() error {
	if err := ApprovalMachine.Transition(string(e.Status), string(StatusReview)); err != nil {
		return err
	}
	e.Status = StatusReview
	e.UpdatedAt = time.Now()
	return nil
}

func (e *EvidenceItem) Approve() error {
	if err := ApprovalMachine.Transition(string(e.Status), string(StatusApproved)); err != nil {
		return err
	}
	e.Status = StatusApproved
	e.UpdatedAt = time.Now()
	return nil
}

func (e *EvidenceItem) Reject() error {
	if err := ApprovalMachine.Transition(string(e.Status), string(StatusDraft)); err != nil {
		return err
	}
	e.Status = StatusDraft
	e.UpdatedAt = time.Now()
	return nil
}

func (e *EvidenceItem) Export() error {
	if err := ApprovalMachine.Transition(string(e.Status), string(StatusExported)); err != nil {
		return err
	}
	e.Status = StatusExported
	e.UpdatedAt = time.Now()
	return nil
}

func (e *EvidenceItem) IsExpired() bool {
	return !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt)
}
