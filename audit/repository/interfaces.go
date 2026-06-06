package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/audit/domain"
)

type AuditFilter struct {
	TenantID  uuid.UUID
	ActorID   uuid.UUID
	Action    domain.Action
	TargetID  string
	Since     time.Time
	Until     time.Time
	Limit     int
	Offset    int
}

type AuditRepository interface {
	Append(ctx context.Context, event *domain.AuditEvent) error
	List(ctx context.Context, filter AuditFilter) ([]*domain.AuditEvent, error)
	Count(ctx context.Context, filter AuditFilter) (int, error)
}
