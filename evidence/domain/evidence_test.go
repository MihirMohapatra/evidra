package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvidence(t *testing.T) {
	now := time.Now()
	input := CreateEvidenceInput{
		TenantID:  uuid.New(),
		Title:     "Test Evidence",
		Content:   "Some content",
		Category:  CategoryPolicy,
		OwnerID:   uuid.New(),
		SourceURL: "https://example.com/doc.pdf",
		Tags:      []string{"tag1", "tag2"},
		ExpiresAt: now.Add(24 * time.Hour),
	}
	item := NewEvidence(input)
	require.NotNil(t, item)
	assert.Equal(t, input.TenantID, item.TenantID)
	assert.Equal(t, input.Title, item.Title)
	assert.Equal(t, input.Content, item.Content)
	assert.Equal(t, input.Category, item.Category)
	assert.Equal(t, input.OwnerID, item.OwnerID)
	assert.Equal(t, input.SourceURL, item.SourceURL)
	assert.Equal(t, input.Tags, item.Tags)
	assert.Equal(t, input.ExpiresAt, item.ExpiresAt)
	assert.Equal(t, StatusDraft, item.Status)
	assert.Equal(t, 1, item.Version)
	assert.NotEqual(t, uuid.Nil, item.ID)
	assert.WithinDuration(t, time.Now(), item.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), item.UpdatedAt, time.Second)
}

func TestNewEvidence_NilTags(t *testing.T) {
	input := CreateEvidenceInput{
		TenantID: uuid.New(),
		Title:    "No tags",
		Category: CategoryAnswer,
		OwnerID:  uuid.New(),
	}
	item := NewEvidence(input)
	assert.NotNil(t, item.Tags)
	assert.Empty(t, item.Tags)
}

func TestTransitionStatus_Valid(t *testing.T) {
	item := NewEvidence(validInput())
	err := item.TransitionStatus(StatusPending)
	assert.NoError(t, err)
	assert.Equal(t, StatusPending, item.Status)
}

func TestTransitionStatus_Invalid(t *testing.T) {
	item := NewEvidence(validInput())
	err := item.TransitionStatus(StatusApproved)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot transition")
	assert.Equal(t, StatusDraft, item.Status)
}

func TestTransitionStatus_UpdatesTimestamp(t *testing.T) {
	item := NewEvidence(validInput())
	original := item.UpdatedAt
	time.Sleep(time.Millisecond)
	_ = item.TransitionStatus(StatusPending)
	assert.True(t, item.UpdatedAt.After(original))
}

func TestIsExpired(t *testing.T) {
	item := NewEvidence(validInput())
	assert.False(t, item.IsExpired())

	item.ExpiresAt = time.Now().Add(-time.Hour)
	assert.True(t, item.IsExpired())

	item.ExpiresAt = time.Time{}
	assert.False(t, item.IsExpired())
}

func TestRenew_NotExpired(t *testing.T) {
	item := NewEvidence(validInput())
	item.Status = StatusDraft
	newExpiry := time.Now().Add(72 * time.Hour)
	original := item.UpdatedAt
	time.Sleep(time.Millisecond)
	item.Renew(newExpiry)
	assert.Equal(t, newExpiry, item.ExpiresAt)
	assert.Equal(t, StatusDraft, item.Status)
	assert.True(t, item.UpdatedAt.After(original))
}

func TestRenew_Expired(t *testing.T) {
	item := NewEvidence(validInput())
	item.Status = StatusExpired
	newExpiry := time.Now().Add(72 * time.Hour)
	item.Renew(newExpiry)
	assert.Equal(t, StatusDraft, item.Status)
	assert.Equal(t, newExpiry, item.ExpiresAt)
}

func validInput() CreateEvidenceInput {
	return CreateEvidenceInput{
		TenantID: uuid.New(),
		Title:    "test",
		Category: CategoryPolicy,
		OwnerID:  uuid.New(),
	}
}
