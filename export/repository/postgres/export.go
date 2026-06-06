package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/export/domain"
)

type ExportRepo struct {
	pool *pgxpool.Pool
}

func NewExportRepo(pool *pgxpool.Pool) *ExportRepo {
	return &ExportRepo{pool: pool}
}

func (r *ExportRepo) Create(ctx context.Context, exp *domain.Export) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exports (id, tenant_id, evidence_id, requester_id, format, file_url, file_size, status, error, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		exp.ID, exp.TenantID, exp.EvidenceID, exp.RequesterID, exp.Format, exp.FileURL, exp.FileSize, exp.Status, exp.Error, exp.CreatedAt, exp.UpdatedAt)
	return err
}

func (r *ExportRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Export, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, evidence_id, requester_id, format, file_url, file_size, status, error, created_at, updated_at
		 FROM exports WHERE id = $1`, id)
	return scanExport(row)
}

func (r *ExportRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Export, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, evidence_id, requester_id, format, file_url, file_size, status, error, created_at, updated_at
		 FROM exports WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExports(rows)
}

func (r *ExportRepo) ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Export, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, evidence_id, requester_id, format, file_url, file_size, status, error, created_at, updated_at
		 FROM exports WHERE evidence_id = $1 ORDER BY created_at DESC`, evidenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExports(rows)
}

func (r *ExportRepo) Update(ctx context.Context, exp *domain.Export) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE exports SET status = $1, file_url = $2, file_size = $3, error = $4, updated_at = $5 WHERE id = $6`,
		exp.Status, exp.FileURL, exp.FileSize, exp.Error, exp.UpdatedAt, exp.ID)
	return err
}

func scanExport(row pgx.Row) (*domain.Export, error) {
	var e domain.Export
	err := row.Scan(&e.ID, &e.TenantID, &e.EvidenceID, &e.RequesterID, &e.Format, &e.FileURL, &e.FileSize, &e.Status, &e.Error, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func scanExports(rows pgx.Rows) ([]*domain.Export, error) {
	var exports []*domain.Export
	for rows.Next() {
		var e domain.Export
		if err := rows.Scan(&e.ID, &e.TenantID, &e.EvidenceID, &e.RequesterID, &e.Format, &e.FileURL, &e.FileSize, &e.Status, &e.Error, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		exports = append(exports, &e)
	}
	if exports == nil {
		return []*domain.Export{}, nil
	}
	return exports, nil
}
