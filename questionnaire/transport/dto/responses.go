package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/questionnaire/domain"
)

type QuestionnaireResponse struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Title     string    `json:"title"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	Status    string    `json:"status"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToQuestionnaireResponse(q *domain.Questionnaire) QuestionnaireResponse {
	return QuestionnaireResponse{
		ID:        q.ID,
		TenantID:  q.TenantID,
		Title:     q.Title,
		FileName:  q.FileName,
		FileURL:   q.FileURL,
		FileType:  q.FileType,
		FileSize:  q.FileSize,
		Status:    string(q.Status),
		Version:   q.Version,
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
	}
}

type QuestionResponse struct {
	ID              uuid.UUID `json:"id"`
	QuestionnaireID uuid.UUID `json:"questionnaire_id"`
	Text            string    `json:"text"`
	Type            string    `json:"type"`
	Order           int       `json:"order"`
	Options         []string  `json:"options"`
	CreatedAt       time.Time `json:"created_at"`
}

func ToQuestionResponse(q *domain.Question) QuestionResponse {
	return QuestionResponse{
		ID:              q.ID,
		QuestionnaireID: q.QuestionnaireID,
		Text:            q.Text,
		Type:            string(q.Type),
		Order:           q.Order,
		Options:         q.Options,
		CreatedAt:       q.CreatedAt,
	}
}

type UploadResponse struct {
	Questionnaire QuestionnaireResponse `json:"questionnaire"`
	UploadURL     string                `json:"upload_url,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
