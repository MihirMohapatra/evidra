package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/audit/domain"
)

type AuditEventResponse struct {
	ID        uuid.UUID              `json:"id"`
	TenantID  uuid.UUID              `json:"tenant_id"`
	ActorID   uuid.UUID              `json:"actor_id"`
	Action    string                 `json:"action"`
	TargetID  string                 `json:"target_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]any         `json:"metadata,omitempty"`
}

func ToAuditEventResponse(e *domain.AuditEvent) AuditEventResponse {
	return AuditEventResponse{
		ID:        e.ID,
		TenantID:  e.TenantID,
		ActorID:   e.ActorID,
		Action:    string(e.Action),
		TargetID:  e.TargetID,
		Timestamp: e.Timestamp,
		Metadata:  e.Metadata,
	}
}

type PaginatedResponse struct {
	Items  any   `json:"items"`
	Total  int   `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
