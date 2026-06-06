package events

import "github.com/google/uuid"

const (
	SubjectQuestionnaireUploaded = "questionnaire.uploaded"
	SubjectQuestionnaireParsed   = "questionnaire.parsed"
	SubjectQuestionnaireFailed   = "questionnaire.failed"
)

type QuestionnaireUploaded struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
}

func (e QuestionnaireUploaded) Subject() string {
	return SubjectQuestionnaireUploaded
}

type QuestionnaireParsed struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	QuestionCount   int       `json:"question_count"`
	QuestionnaireID uuid.UUID `json:"questionnaire_id"`
}

func (e QuestionnaireParsed) Subject() string {
	return SubjectQuestionnaireParsed
}

type QuestionnaireFailed struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Error    string    `json:"error"`
}

func (e QuestionnaireFailed) Subject() string {
	return SubjectQuestionnaireFailed
}
