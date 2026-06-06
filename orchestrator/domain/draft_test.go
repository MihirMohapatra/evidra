package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDraft(t *testing.T) {
	qID := uuid.New()
	evIDs := []uuid.UUID{uuid.New(), uuid.New()}
	d := NewDraft(qID, "What is SOC 2?", "SOC 2 is a compliance framework.", 0.95, "gpt-4o", evIDs, "Based on evidence...")
	require.NotNil(t, d)
	assert.Equal(t, qID, d.QuestionID)
	assert.Equal(t, "What is SOC 2?", d.QuestionText)
	assert.Equal(t, "SOC 2 is a compliance framework.", d.Answer)
	assert.Equal(t, 0.95, d.Confidence)
	assert.Equal(t, "gpt-4o", d.ModelUsed)
	assert.Equal(t, evIDs, d.EvidenceIDs)
	assert.Equal(t, "Based on evidence...", d.Reasoning)
	assert.Equal(t, DraftPending, d.Status)
	assert.NotEqual(t, uuid.Nil, d.ID)
	assert.WithinDuration(t, time.Now(), d.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), d.UpdatedAt, time.Second)
}

func TestNewDraft_NilEvidenceIDs(t *testing.T) {
	d := NewDraft(uuid.New(), "question", "answer", 0.5, "model", nil, "")
	assert.NotNil(t, d.EvidenceIDs)
	assert.Empty(t, d.EvidenceIDs)
}
