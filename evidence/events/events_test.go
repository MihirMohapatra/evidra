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

func TestEvidenceStatusChanged_Subject(t *testing.T) {
	tests := []struct {
		status  domain.Status
		subject string
	}{
		{domain.StatusApproved, SubjectEvidenceApproved},
		{domain.StatusDraft, SubjectEvidenceRejected},
		{domain.StatusExported, SubjectEvidenceExported},
		{domain.StatusReview, SubjectEvidenceUpdated},
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
