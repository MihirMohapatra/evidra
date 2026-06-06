package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/audit/domain"
	"github.com/evidra/evidra/audit/repository"
)

type AuditRepo struct {
	pool *pgxpool.Pool
}

func NewAuditRepo(pool *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{pool: pool}
}

var _ repository.AuditRepository = (*AuditRepo)(nil)

func (r *AuditRepo) Append(ctx context.Context, event *domain.AuditEvent) error {
	meta, err := json.Marshal(event.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `INSERT INTO audit_events (id, tenant_id, actor_id, action, target_id, timestamp, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = r.pool.Exec(ctx, query,
		event.ID, event.TenantID, event.ActorID, event.Action, event.TargetID, event.Timestamp, meta)
	return err
}

func (r *AuditRepo) List(ctx context.Context, filter repository.AuditFilter) ([]*domain.AuditEvent, error) {
	query := `SELECT id, tenant_id, actor_id, action, target_id, timestamp, metadata FROM audit_events WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.TenantID != uuid.Nil {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, filter.TenantID)
		argIdx++
	}
	if filter.ActorID != uuid.Nil {
		query += fmt.Sprintf(" AND actor_id = $%d", argIdx)
		args = append(args, filter.ActorID)
		argIdx++
	}
	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, filter.Action)
		argIdx++
	}
	if filter.TargetID != "" {
		query += fmt.Sprintf(" AND target_id = $%d", argIdx)
		args = append(args, filter.TargetID)
		argIdx++
	}
	if !filter.Since.IsZero() {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, filter.Since)
		argIdx++
	}
	if !filter.Until.IsZero() {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, filter.Until)
		argIdx++
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	} else {
		query += " LIMIT 50"
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	defer rows.Close()

	var events []*domain.AuditEvent
	for rows.Next() {
		e := &domain.AuditEvent{}
		var metaJSON []byte
		if err := rows.Scan(&e.ID, &e.TenantID, &e.ActorID, &e.Action, &e.TargetID, &e.Timestamp, &metaJSON); err != nil {
			return nil, fmt.Errorf("scan audit: %w", err)
		}
		_ = json.Unmarshal(metaJSON, &e.Metadata)
		if e.Metadata == nil {
			e.Metadata = map[string]any{}
		}
		events = append(events, e)
	}
	if events == nil {
		events = []*domain.AuditEvent{}
	}
	return events, nil
}

func (r *AuditRepo) Count(ctx context.Context, filter repository.AuditFilter) (int, error) {
	query := `SELECT COUNT(*) FROM audit_events WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.TenantID != uuid.Nil {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, filter.TenantID)
		argIdx++
	}
	if filter.ActorID != uuid.Nil {
		query += fmt.Sprintf(" AND actor_id = $%d", argIdx)
		args = append(args, filter.ActorID)
		argIdx++
	}
	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, filter.Action)
		argIdx++
	}
	if !filter.Since.IsZero() {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, filter.Since)
		argIdx++
	}
	if !filter.Until.IsZero() {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, filter.Until)
		argIdx++
	}

	_ = argIdx

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}
