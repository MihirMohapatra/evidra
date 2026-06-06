package domain

import (
	"time"

	"github.com/google/uuid"
)

type Questionnaire struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Title     string
	FileName  string
	FileURL   string
	FileType  string
	FileSize  int64
	Status    Status
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQuestionnaire(tenantID uuid.UUID, title, fileName, fileURL, fileType string, fileSize int64) *Questionnaire {
	now := time.Now()
	return &Questionnaire{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Title:     title,
		FileName:  fileName,
		FileURL:   fileURL,
		FileType:  fileType,
		FileSize:  fileSize,
		Status:    StatusUploaded,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (q *Questionnaire) TransitionStatus(next Status) error {
	if !q.Status.CanTransitionTo(next) {
		return NewStatusError(q.Status, next)
	}
	q.Status = next
	q.UpdatedAt = time.Now()
	return nil
}
