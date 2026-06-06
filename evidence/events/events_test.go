package events

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/evidra/evidra/evidence/domain"
)

func TestEvidenceCreated_Subject(t *testing.T) {
	e := EvidenceCreated{ID: uuid.New(), TenantID: uuid.New(), Title: "test"}
	assert.Equal(t, SubjectEvidenceCreated, e.Subject())
}

func TestEvidenceUpdated_Subject(t *testing.T) {
	e := EvidenceUpdated{ID: uuid.New(), TenantID: uuid.New()}
	assert.Equal(t, SubjectEvidenceUpdated, e.Subject())
}

func TestEvidenceDeleted_Subject(t *testing.T) {
	e := EvidenceDeleted{ID: uuid.New(), TenantID: uuid.New()}
	assert.Equal(t, SubjectEvidenceDeleted, e.Subject())
}

func TestEvidenceExpired_Subject(t *testing.T) {
	e := EvidenceExpired{ID: uuid.New(), TenantID: uuid.New(), Title: "expired"}
	assert.Equal(t, SubjectEvidenceExpired, e.Subject())
}

func TestEvidenceStatusChanged_Subject(t *testing.T) {
	tests := []struct {
		status  domain.ApprovalStatus
		subject string
	}{
		{domain.StatusApproved, SubjectEvidenceApproved},
		{domain.StatusRejected, SubjectEvidenceRejected},
		{domain.StatusExported, SubjectEvidenceExported},
		{domain.StatusDraft, SubjectEvidenceUpdated},
		{domain.StatusPending, SubjectEvidenceUpdated},
	}
	for _, tt := range tests {
		e := EvidenceStatusChanged{
			ID:         uuid.New(),
			TenantID:   uuid.New(),
			Status:     tt.status,
			ReviewerID: uuid.New(),
		}
		assert.Equal(t, tt.subject, e.Subject(), "status=%s", tt.status)
	}
}

func TestEvidenceExported_Subject(t *testing.T) {
	e := EvidenceExported{ID: uuid.New(), TenantID: uuid.New(), Title: "exported"}
	assert.Equal(t, SubjectEvidenceExported, e.Subject())
}
