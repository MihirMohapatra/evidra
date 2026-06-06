package domain

import (
	"time"

	"github.com/google/uuid"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeSingleChoice   QuestionType = "single_choice"
	QuestionTypeOpenEnded      QuestionType = "open_ended"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeFillBlank      QuestionType = "fill_blank"
)

type Question struct {
	ID              uuid.UUID
	QuestionnaireID uuid.UUID
	Text            string
	Type            QuestionType
	Order           int
	Options         []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewQuestion(questionnaireID uuid.UUID, text string, qtype QuestionType, order int, options []string) *Question {
	now := time.Now()
	if options == nil {
		options = []string{}
	}
	return &Question{
		ID:              uuid.New(),
		QuestionnaireID: questionnaireID,
		Text:            text,
		Type:            qtype,
		Order:           order,
		Options:         options,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
