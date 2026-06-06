package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuestionnaire(t *testing.T) {
	tID := uuid.New()
	q := NewQuestionnaire(tID, "SOC 2 Questionnaire", "soc2.pdf", "https://minio/soc2.pdf", "application/pdf", 1024)
	require.NotNil(t, q)
	assert.Equal(t, tID, q.TenantID)
	assert.Equal(t, "SOC 2 Questionnaire", q.Title)
	assert.Equal(t, "soc2.pdf", q.FileName)
	assert.Equal(t, "https://minio/soc2.pdf", q.FileURL)
	assert.Equal(t, "application/pdf", q.FileType)
	assert.Equal(t, int64(1024), q.FileSize)
	assert.Equal(t, StatusUploaded, q.Status)
	assert.Equal(t, 1, q.Version)
	assert.NotEqual(t, uuid.Nil, q.ID)
	assert.WithinDuration(t, time.Now(), q.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), q.UpdatedAt, time.Second)
}

func TestQuestionnaireTransitionStatus_Valid(t *testing.T) {
	q := NewQuestionnaire(uuid.New(), "test", "f.pdf", "url", "pdf", 100)
	err := q.TransitionStatus(StatusQueued)
	assert.NoError(t, err)
	assert.Equal(t, StatusQueued, q.Status)
}

func TestQuestionnaireTransitionStatus_Invalid(t *testing.T) {
	q := NewQuestionnaire(uuid.New(), "test", "f.pdf", "url", "pdf", 100)
	err := q.TransitionStatus(StatusParsed)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot transition")
	assert.Equal(t, StatusUploaded, q.Status)
}

func TestQuestionnaireTransitionStatus_UpdatesTimestamp(t *testing.T) {
	q := NewQuestionnaire(uuid.New(), "test", "f.pdf", "url", "pdf", 100)
	original := q.UpdatedAt
	time.Sleep(time.Millisecond)
	_ = q.TransitionStatus(StatusQueued)
	assert.True(t, q.UpdatedAt.After(original))
}
