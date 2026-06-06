package domain

import (
	"time"

	"github.com/google/uuid"
)

type DraftStatus string

const (
	DraftPending  DraftStatus = "pending"
	DraftApproved DraftStatus = "approved"
	DraftRejected DraftStatus = "rejected"
)

type Draft struct {
	ID           uuid.UUID
	QuestionID   uuid.UUID
	QuestionText string
	Answer       string
	Confidence   float64
	ModelUsed    string
	EvidenceIDs  []uuid.UUID
	Reasoning    string
	Status       DraftStatus
	Feedback     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewDraft(questionID uuid.UUID, questionText, answer string, confidence float64, modelUsed string, evidenceIDs []uuid.UUID, reasoning string) *Draft {
	now := time.Now()
	if evidenceIDs == nil {
		evidenceIDs = []uuid.UUID{}
	}
	return &Draft{
		ID:           uuid.New(),
		QuestionID:   questionID,
		QuestionText: questionText,
		Answer:       answer,
		Confidence:   confidence,
		ModelUsed:    modelUsed,
		EvidenceIDs:  evidenceIDs,
		Reasoning:    reasoning,
		Status:       DraftPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
