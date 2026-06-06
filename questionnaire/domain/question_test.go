package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuestion(t *testing.T) {
	qID := uuid.New()
	q := NewQuestion(qID, "What is SOC 2?", QuestionTypeOpenEnded, 1, nil)
	require.NotNil(t, q)
	assert.Equal(t, qID, q.QuestionnaireID)
	assert.Equal(t, "What is SOC 2?", q.Text)
	assert.Equal(t, QuestionTypeOpenEnded, q.Type)
	assert.Equal(t, 1, q.Order)
	assert.NotNil(t, q.Options)
	assert.Empty(t, q.Options)
	assert.NotEqual(t, uuid.Nil, q.ID)
	assert.WithinDuration(t, time.Now(), q.CreatedAt, time.Second)
}

func TestNewQuestion_WithOptions(t *testing.T) {
	qID := uuid.New()
	opts := []string{"Option A", "Option B", "Option C"}
	q := NewQuestion(qID, "Choose?", QuestionTypeSingleChoice, 2, opts)
	require.NotNil(t, q)
	assert.Equal(t, opts, q.Options)
	assert.Equal(t, 2, q.Order)
}
