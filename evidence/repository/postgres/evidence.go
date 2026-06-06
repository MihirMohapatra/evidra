package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/evidence/domain"
	"github.com/evidra/evidra/evidence/repository"
)

type EvidenceRepo struct {
	pool *pgxpool.Pool
}

func NewEvidenceRepo(pool *pgxpool.Pool) *EvidenceRepo {
	return &EvidenceRepo{pool: pool}
}

func (r *EvidenceRepo) Create(ctx context.Context, item *domain.EvidenceItem) error {
	query := `INSERT INTO evidence_items (id, tenant_id, title, content, category, status, owner_id, source_url, tags, version, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.pool.Exec(ctx, query, item.ID, item.TenantID, item.Title, item.Content, item.Category, item.Status, item.OwnerID, item.SourceURL, item.Tags, item.Version, item.ExpiresAt, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *EvidenceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.EvidenceItem, error) {
	query := `SELECT id, tenant_id, title, content, category, status, owner_id, source_url, tags, version, expires_at, created_at, updated_at FROM evidence_items WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	item := &domain.EvidenceItem{}
	err := row.Scan(&item.ID, &item.TenantID, &item.Title, &item.Content, &item.Category, &item.Status, &item.OwnerID, &item.SourceURL, &item.Tags, &item.Version, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return item, err
}

func (r *EvidenceRepo) List(ctx context.Context, filter repository.EvidenceFilter) ([]*domain.EvidenceItem, error) {
	query := `SELECT id, tenant_id, title, content, category, status, owner_id, source_url, tags, version, expires_at, created_at, updated_at FROM evidence_items WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.TenantID != uuid.Nil {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, filter.TenantID)
		argIdx++
	}
	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.OwnerID != uuid.Nil {
		query += fmt.Sprintf(" AND owner_id = $%d", argIdx)
		args = append(args, filter.OwnerID)
		argIdx++
	}
	if filter.Expiring {
		query += fmt.Sprintf(" AND expires_at IS NOT NULL AND expires_at > NOW() AND expires_at <= NOW() + $%d::interval", argIdx)
		args = append(args, fmt.Sprintf("%d seconds", int(filter.ExpireWin.Seconds())))
		argIdx++
	}
	if filter.Expired {
		query += " AND expires_at IS NOT NULL AND expires_at <= NOW()"
	}
	if filter.Tag != "" {
		query += fmt.Sprintf(" AND $%d = ANY(tags)", argIdx)
		args = append(args, filter.Tag)
		argIdx++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argIdx, argIdx+1)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
		argIdx += 2
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.EvidenceItem
	for rows.Next() {
		item := &domain.EvidenceItem{}
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Title, &item.Content, &item.Category, &item.Status, &item.OwnerID, &item.SourceURL, &item.Tags, &item.Version, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *EvidenceRepo) Count(ctx context.Context, filter repository.EvidenceFilter) (int, error) {
	query := `SELECT COUNT(*) FROM evidence_items WHERE 1=1`
	args := []any{}
	argIdx := 1

	buildWhere(&query, &args, &argIdx, filter)

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *EvidenceRepo) Update(ctx context.Context, item *domain.EvidenceItem) error {
	query := `UPDATE evidence_items SET title = $1, content = $2, category = $3, status = $4, owner_id = $5, source_url = $6, tags = $7, version = $8, expires_at = $9, updated_at = $10 WHERE id = $11`
	tag, err := r.pool.Exec(ctx, query, item.Title, item.Content, item.Category, item.Status, item.OwnerID, item.SourceURL, item.Tags, item.Version, item.ExpiresAt, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EvidenceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM evidence_items WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func buildWhere(query *string, args *[]any, argIdx *int, filter repository.EvidenceFilter) {
	if filter.TenantID != uuid.Nil {
		*query += fmt.Sprintf(" AND tenant_id = $%d", *argIdx)
		*args = append(*args, filter.TenantID)
		*argIdx++
	}
	if filter.Category != "" {
		*query += fmt.Sprintf(" AND category = $%d", *argIdx)
		*args = append(*args, filter.Category)
		*argIdx++
	}
	if filter.Status != "" {
		*query += fmt.Sprintf(" AND status = $%d", *argIdx)
		*args = append(*args, filter.Status)
		*argIdx++
	}
	if filter.OwnerID != uuid.Nil {
		*query += fmt.Sprintf(" AND owner_id = $%d", *argIdx)
		*args = append(*args, filter.OwnerID)
		*argIdx++
	}
	if filter.Expiring {
		*query += fmt.Sprintf(" AND expires_at IS NOT NULL AND expires_at > NOW() AND expires_at <= NOW() + $%d::interval", *argIdx)
		*args = append(*args, fmt.Sprintf("%d seconds", int(filter.ExpireWin.Seconds())))
		*argIdx++
	}
	if filter.Expired {
		*query += " AND expires_at IS NOT NULL AND expires_at <= NOW()"
	}
	if filter.Tag != "" {
		*query += fmt.Sprintf(" AND $%d = ANY(tags)", *argIdx)
		*args = append(*args, filter.Tag)
		*argIdx++
	}
	if filter.Search != "" {
		*query += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", *argIdx, *argIdx+1)
		searchPattern := "%" + filter.Search + "%"
		*args = append(*args, searchPattern, searchPattern)
		*argIdx += 2
	}
	if filter.Limit > 0 {
		*query += fmt.Sprintf(" LIMIT $%d", *argIdx)
		*args = append(*args, filter.Limit)
		*argIdx++
	}
	if filter.Offset > 0 {
		*query += fmt.Sprintf(" OFFSET $%d", *argIdx)
		*args = append(*args, filter.Offset)
		*argIdx++
	}
}

var _ repository.EvidenceRepository = (*EvidenceRepo)(nil)
