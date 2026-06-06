package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ApprovalStatus string

const (
	StatusDraft     ApprovalStatus = "draft"
	StatusPending   ApprovalStatus = "pending_review"
	StatusApproved  ApprovalStatus = "approved"
	StatusRejected  ApprovalStatus = "rejected"
	StatusExpired   ApprovalStatus = "expired"
	StatusArchived  ApprovalStatus = "archived"
)

func (s ApprovalStatus) Valid() bool {
	switch s {
	case StatusDraft, StatusPending, StatusApproved, StatusRejected, StatusExpired, StatusArchived:
		return true
	}
	return false
}

func (s ApprovalStatus) CanTransitionTo(next ApprovalStatus) bool {
	transitions := map[ApprovalStatus][]ApprovalStatus{
		StatusDraft:    {StatusPending, StatusArchived},
		StatusPending:  {StatusApproved, StatusRejected},
		StatusApproved: {StatusExpired, StatusArchived},
		StatusRejected: {StatusDraft, StatusArchived},
		StatusExpired:  {StatusDraft, StatusArchived},
		StatusArchived: {},
	}
	allowed, ok := transitions[s]
	if !ok {
		return false
	}
	for _, st := range allowed {
		if st == next {
			return true
		}
	}
	return false
}

func NewStatusError(current, target ApprovalStatus) error {
	return fmt.Errorf("cannot transition from %s to %s", current, target)
}

type Approval struct {
	ID         uuid.UUID
	EvidenceID uuid.UUID
	ReviewerID uuid.UUID
	Status     ApprovalStatus
	Comment    string
	CreatedAt  time.Time
}

func NewApproval(evidenceID, reviewerID uuid.UUID, status ApprovalStatus, comment string) *Approval {
	return &Approval{
		ID:         uuid.New(),
		EvidenceID: evidenceID,
		ReviewerID: reviewerID,
		Status:     status,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
}
