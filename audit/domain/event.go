package domain

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	ActionQuestionUploaded   Action = "QUESTION_UPLOADED"
	ActionAIGenerated        Action = "AI_GENERATED"
	ActionAnswerApproved     Action = "ANSWER_APPROVED"
	ActionDocumentExported   Action = "DOCUMENT_EXPORTED"
	ActionRoleChanged        Action = "ROLE_CHANGED"
	ActionEvidenceCreated    Action = "EVIDENCE_CREATED"
	ActionEvidenceUpdated    Action = "EVIDENCE_UPDATED"
	ActionEvidenceDeleted    Action = "EVIDENCE_DELETED"
	ActionEvidenceSubmitted  Action = "EVIDENCE_SUBMITTED"
	ActionEvidenceExpired    Action = "EVIDENCE_EXPIRED"
	ActionUserLogin          Action = "USER_LOGIN"
	ActionUserLogout         Action = "USER_LOGOUT"
	ActionAPIKeyCreated      Action = "API_KEY_CREATED"
	ActionAPIKeyRevoked      Action = "API_KEY_REVOKED"
	ActionOrganizationCreated Action = "ORGANIZATION_CREATED"
	ActionOrganizationUpdated Action = "ORGANIZATION_UPDATED"
	ActionOrganizationDeleted Action = "ORGANIZATION_DELETED"
	ActionUserCreated        Action = "USER_CREATED"
	ActionUserUpdated        Action = "USER_UPDATED"
	ActionUserDeleted        Action = "USER_DELETED"
	ActionDraftCreated               Action = "DRAFT_CREATED"
	ActionDraftApproved              Action = "DRAFT_APPROVED"
	ActionDraftRejected              Action = "DRAFT_REJECTED"

	ActionComplianceFrameworkCreated Action = "COMPLIANCE_FRAMEWORK_CREATED"
	ActionComplianceControlCreated   Action = "COMPLIANCE_CONTROL_CREATED"
	ActionComplianceEvidenceMapped   Action = "COMPLIANCE_EVIDENCE_MAPPED"
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
