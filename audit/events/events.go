package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubjectAuditRecorded = "audit.recorded"
)

type AuditRecorded struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ActorID   uuid.UUID `json:"actor_id"`
	Action    string    `json:"action"`
	TargetID  string    `json:"target_id,omitempty"`
	Timestamp string    `json:"timestamp"`
}

func (e AuditRecorded) Subject() string { return SubjectAuditRecorded }

func NewAuditRecorded(tenantID, actorID uuid.UUID, action, targetID string) AuditRecorded {
	return AuditRecorded{
		ID:        uuid.New(),
		TenantID:  tenantID,
		ActorID:   actorID,
		Action:    action,
		TargetID:  targetID,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}
}
