package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditEvent(t *testing.T) {
	tID := uuid.New()
	aID := uuid.New()
	meta := map[string]any{"key": "value"}
	e := NewAuditEvent(tID, aID, ActionEvidenceCreated, "ev-123", meta)
	require.NotNil(t, e)
	assert.Equal(t, tID, e.TenantID)
	assert.Equal(t, aID, e.ActorID)
	assert.Equal(t, ActionEvidenceCreated, e.Action)
	assert.Equal(t, "ev-123", e.TargetID)
	assert.Equal(t, meta, e.Metadata)
	assert.NotEqual(t, uuid.Nil, e.ID)
	assert.WithinDuration(t, time.Now(), e.Timestamp, time.Second)
}

func TestNewAuditEvent_NilMetadata(t *testing.T) {
	e := NewAuditEvent(uuid.New(), uuid.New(), ActionUserLogin, "", nil)
	assert.NotNil(t, e.Metadata)
	assert.Empty(t, e.Metadata)
}

func TestAuditActions(t *testing.T) {
	actions := []Action{
		ActionQuestionUploaded,
		ActionAIGenerated,
		ActionAnswerApproved,
		ActionDocumentExported,
		ActionRoleChanged,
		ActionEvidenceCreated,
		ActionEvidenceDeleted,
		ActionEvidenceExpired,
		ActionUserLogin,
		ActionAPIKeyCreated,
	}
	for _, a := range actions {
		assert.NotEmpty(t, string(a), "action %s should not be empty", a)
	}
}
