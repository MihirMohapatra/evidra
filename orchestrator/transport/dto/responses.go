package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/orchestrator/domain"
)

type DraftResponse struct {
	ID           uuid.UUID   `json:"id"`
	QuestionID   uuid.UUID   `json:"question_id"`
	QuestionText string      `json:"question_text"`
	Answer       string      `json:"answer"`
	Confidence   float64     `json:"confidence"`
	ModelUsed    string      `json:"model_used"`
	EvidenceIDs  []uuid.UUID `json:"evidence_ids"`
	Reasoning    string      `json:"reasoning"`
	Status       string      `json:"status"`
	Feedback     string      `json:"feedback,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func ToDraftResponse(d *domain.Draft) DraftResponse {
	return DraftResponse{
		ID:           d.ID,
		QuestionID:   d.QuestionID,
		QuestionText: d.QuestionText,
		Answer:       d.Answer,
		Confidence:   d.Confidence,
		ModelUsed:    d.ModelUsed,
		EvidenceIDs:  d.EvidenceIDs,
		Reasoning:    d.Reasoning,
		Status:       string(d.Status),
		Feedback:     d.Feedback,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type EvidenceResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Score     float64   `json:"score"`
	SourceURL string    `json:"source_url,omitempty"`
}

func ToEvidenceResponse(e domain.Evidence) EvidenceResponse {
	return EvidenceResponse{
		ID:        e.ID,
		Title:     e.Title,
		Content:   e.Content,
		Category:  e.Category,
		Score:     e.Score,
		SourceURL: e.SourceURL,
	}
}

type AnswerResponse struct {
	Draft    DraftResponse      `json:"draft"`
	Evidence []EvidenceResponse `json:"evidence"`
}

type ListDraftsResponse struct {
	Items []DraftResponse `json:"items"`
	Total int             `json:"total"`
	Limit int             `json:"limit"`
	Offset int            `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
