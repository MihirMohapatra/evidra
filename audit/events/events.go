package events

import "github.com/google/uuid"

const (
	SubjectAuditRecorded = "audit.recorded"
)

type AuditRecorded struct {
	ID        uuid.UUID    `json:"id"`
	TenantID  uuid.UUID    `json:"tenant_id"`
	ActorID   uuid.UUID    `json:"actor_id"`
	Action    string       `json:"action"`
	TargetID  string       `json:"target_id,omitempty"`
	Timestamp string       `json:"timestamp"`
}

func (e AuditRecorded) Subject() string { return SubjectAuditRecorded }
