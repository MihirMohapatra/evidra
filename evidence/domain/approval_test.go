package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovalStatus_Valid(t *testing.T) {
	tests := []struct {
		status ApprovalStatus
		valid  bool
	}{
		{StatusDraft, true},
		{StatusPending, true},
		{StatusApproved, true},
		{StatusRejected, true},
		{StatusExpired, true},
		{StatusArchived, true},
		{StatusExported, true},
		{"unknown", false},
		{"", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.valid, tt.status.Valid(), "status=%q", tt.status)
	}
}

func TestApprovalStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		current ApprovalStatus
		next    ApprovalStatus
		allowed bool
	}{
		{StatusDraft, StatusPending, true},
		{StatusDraft, StatusArchived, true},
		{StatusDraft, StatusApproved, false},
		{StatusDraft, StatusRejected, false},
		{StatusPending, StatusApproved, true},
		{StatusPending, StatusRejected, true},
		{StatusPending, StatusDraft, false},
		{StatusApproved, StatusExported, true},
		{StatusApproved, StatusArchived, true},
		{StatusApproved, StatusExpired, true},
		{StatusApproved, StatusDraft, false},
		{StatusExported, StatusArchived, true},
		{StatusExported, StatusDraft, false},
		{StatusRejected, StatusDraft, true},
		{StatusRejected, StatusArchived, true},
		{StatusRejected, StatusApproved, false},
		{StatusExpired, StatusDraft, true},
		{StatusExpired, StatusArchived, true},
		{StatusExpired, StatusApproved, false},
		{StatusArchived, StatusDraft, false},
		{StatusArchived, StatusPending, false},
		{StatusArchived, StatusApproved, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.allowed, tt.current.CanTransitionTo(tt.next), "%s -> %s", tt.current, tt.next)
	}
}

func TestNewStatusError(t *testing.T) {
	err := NewStatusError(StatusDraft, StatusApproved)
	assert.Equal(t, "cannot transition from draft to approved", err.Error())
}

func TestNewApproval(t *testing.T) {
	eID := uuid.New()
	rID := uuid.New()
	a := NewApproval(eID, rID, StatusApproved, "looks good")
	require.NotNil(t, a)
	assert.Equal(t, eID, a.EvidenceID)
	assert.Equal(t, rID, a.ReviewerID)
	assert.Equal(t, StatusApproved, a.Status)
	assert.Equal(t, "looks good", a.Comment)
	assert.NotEqual(t, uuid.Nil, a.ID)
	assert.WithinDuration(t, time.Now(), a.CreatedAt, time.Second)
}
