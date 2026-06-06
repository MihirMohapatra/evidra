package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatus_Valid(t *testing.T) {
	tests := []struct {
		status Status
		valid  bool
	}{
		{StatusDraft, true},
		{StatusReview, true},
		{StatusApproved, true},
		{StatusExported, true},
		{"PENDING", false},
		{"REJECTED", false},
		{"EXPIRED", false},
		{"", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.valid, tt.status.Valid(), "status=%q", tt.status)
	}
}

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		current Status
		next    Status
		allowed bool
	}{
		{StatusDraft, StatusReview, true},
		{StatusDraft, StatusApproved, false},
		{StatusDraft, StatusExported, false},
		{StatusReview, StatusApproved, true},
		{StatusReview, StatusDraft, true},
		{StatusReview, StatusExported, false},
		{StatusApproved, StatusExported, true},
		{StatusApproved, StatusDraft, false},
		{StatusApproved, StatusReview, false},
		{StatusExported, StatusDraft, false},
		{StatusExported, StatusReview, false},
		{StatusExported, StatusApproved, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.allowed, tt.current.CanTransitionTo(tt.next), "%s -> %s", tt.current, tt.next)
	}
}

func TestNewStatusError(t *testing.T) {
	err := NewStatusError(StatusDraft, StatusApproved)
	assert.Equal(t, "cannot transition from DRAFT to APPROVED", err.Error())
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
}
