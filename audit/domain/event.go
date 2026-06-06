package domain

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	ActionQuestionUploaded Action = "QUESTION_UPLOADED"
	ActionAIGenerated      Action = "AI_GENERATED"
	ActionAnswerApproved   Action = "ANSWER_APPROVED"
	ActionDocumentExported Action = "DOCUMENT_EXPORTED"
	ActionRoleChanged      Action = "ROLE_CHANGED"
	ActionEvidenceCreated  Action = "EVIDENCE_CREATED"
	ActionEvidenceDeleted  Action = "EVIDENCE_DELETED"
	ActionEvidenceExpired  Action = "EVIDENCE_EXPIRED"
	ActionUserLogin        Action = "USER_LOGIN"
	ActionAPIKeyCreated    Action = "API_KEY_CREATED"
)

type AuditEvent struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	ActorID   uuid.UUID
	Action    Action
	TargetID  string
	Timestamp time.Time
	Metadata  map[string]any
}

func NewAuditEvent(tenantID, actorID uuid.UUID, action Action, targetID string, metadata map[string]any) *AuditEvent {
	if metadata == nil {
		metadata = map[string]any{}
	}
	return &AuditEvent{
		ID:        uuid.New(),
		TenantID:  tenantID,
		ActorID:   actorID,
		Action:    action,
		TargetID:  targetID,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}
}
