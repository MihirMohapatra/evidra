package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/evidence/domain"
	"github.com/evidra/evidra/evidence/repository"
)

type ApprovalRepo struct {
	pool *pgxpool.Pool
}

func NewApprovalRepo(pool *pgxpool.Pool) *ApprovalRepo {
	return &ApprovalRepo{pool: pool}
}

func (r *ApprovalRepo) Create(ctx context.Context, a *domain.Approval) error {
	query := `INSERT INTO approvals (id, evidence_id, reviewer_id, status, comment, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, a.ID, a.EvidenceID, a.ReviewerID, a.Status, a.Comment, a.CreatedAt)
	return err
}

func (r *ApprovalRepo) ListByEvidence(ctx context.Context, evidenceID uuid.UUID) ([]*domain.Approval, error) {
	query := `SELECT id, evidence_id, reviewer_id, status, comment, created_at FROM approvals WHERE evidence_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, evidenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []*domain.Approval
	for rows.Next() {
		a := &domain.Approval{}
		if err := rows.Scan(&a.ID, &a.EvidenceID, &a.ReviewerID, &a.Status, &a.Comment, &a.CreatedAt); err != nil {
			return nil, err
		}
		approvals = append(approvals, a)
	}
	return approvals, rows.Err()
}

var _ repository.ApprovalRepository = (*ApprovalRepo)(nil)
