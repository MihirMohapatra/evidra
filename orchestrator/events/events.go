package events

import (
	"github.com/google/uuid"
)

const (
	SubjectDraftCreated = "orchestrator.draft.created"
	SubjectDraftApproved = "orchestrator.draft.approved"
	SubjectDraftRejected = "orchestrator.draft.rejected"
)

type DraftCreated struct {
	ID           uuid.UUID `json:"id"`
	QuestionID   uuid.UUID `json:"question_id"`
	QuestionText string    `json:"question_text"`
	ModelUsed    string    `json:"model_used"`
	Confidence   float64   `json:"confidence"`
	EvidenceIDs  []uuid.UUID `json:"evidence_ids"`
}

func (e DraftCreated) Subject() string { return SubjectDraftCreated }

type DraftStatusChanged struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Feedback string   `json:"feedback"`
}

func (e DraftStatusChanged) Subject() string {
	switch e.Status {
	case "approved":
		return SubjectDraftApproved
	case "rejected":
		return SubjectDraftRejected
	default:
		return SubjectDraftCreated
	}
}
