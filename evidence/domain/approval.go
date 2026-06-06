package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/pkg/workflow"
)

type Status string

const (
	StatusDraft    Status = "DRAFT"
	StatusReview   Status = "NEEDS_REVIEW"
	StatusApproved Status = "APPROVED"
	StatusExported Status = "EXPORTED"
)

var ApprovalMachine = workflow.New([]workflow.TransitionDef{
	{From: string(StatusDraft), To: string(StatusReview), Name: "submit"},
	{From: string(StatusReview), To: string(StatusApproved), Name: "approve"},
	{From: string(StatusReview), To: string(StatusDraft), Name: "reject"},
	{From: string(StatusApproved), To: string(StatusExported), Name: "export"},
})

func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusReview, StatusApproved, StatusExported:
		return true
	}
	return false
}

func (s Status) CanTransitionTo(next Status) bool {
	return ApprovalMachine.CanTransition(string(s), string(next))
}

func NewStatusError(current, target Status) error {
	return fmt.Errorf("cannot transition from %s to %s", current, target)
}

type Approval struct {
	ID         uuid.UUID
	EvidenceID uuid.UUID
	ReviewerID uuid.UUID
	Status     Status
	Comment    string
	CreatedAt  time.Time
}

func NewApproval(evidenceID, reviewerID uuid.UUID, status Status, comment string) *Approval {
	return &Approval{
		ID:         uuid.New(),
		EvidenceID: evidenceID,
		ReviewerID: reviewerID,
		Status:     status,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
}
