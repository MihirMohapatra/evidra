package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/repository"
)

type DraftRepo struct {
	pool *pgxpool.Pool
}

func NewDraftRepo(pool *pgxpool.Pool) *DraftRepo {
	return &DraftRepo{pool: pool}
}

var _ repository.DraftRepository = (*DraftRepo)(nil)

func (r *DraftRepo) Create(ctx context.Context, draft *domain.Draft) error {
	query := `INSERT INTO drafts (id, question_id, question_text, answer, confidence, model_used, evidence_ids, reasoning, status, feedback, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.pool.Exec(ctx, query,
		draft.ID, draft.QuestionID, draft.QuestionText, draft.Answer,
		draft.Confidence, draft.ModelUsed, draft.EvidenceIDs,
		draft.Reasoning, draft.Status, draft.Feedback,
		draft.CreatedAt, draft.UpdatedAt)
	return err
}

func (r *DraftRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Draft, error) {
	query := `SELECT id, question_id, question_text, answer, confidence, model_used, evidence_ids, reasoning, status, feedback, created_at, updated_at FROM drafts WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	d := &domain.Draft{}
	err := row.Scan(&d.ID, &d.QuestionID, &d.QuestionText, &d.Answer, &d.Confidence, &d.ModelUsed, &d.EvidenceIDs, &d.Reasoning, &d.Status, &d.Feedback, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return d, err
}

func (r *DraftRepo) List(ctx context.Context, limit, offset int) ([]*domain.Draft, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, question_id, question_text, answer, confidence, model_used, evidence_ids, reasoning, status, feedback, created_at, updated_at FROM drafts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list drafts: %w", err)
	}
	defer rows.Close()

	var drafts []*domain.Draft
	for rows.Next() {
		d := &domain.Draft{}
		if err := rows.Scan(&d.ID, &d.QuestionID, &d.QuestionText, &d.Answer, &d.Confidence, &d.ModelUsed, &d.EvidenceIDs, &d.Reasoning, &d.Status, &d.Feedback, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan draft: %w", err)
		}
		drafts = append(drafts, d)
	}
	if drafts == nil {
		drafts = []*domain.Draft{}
	}
	return drafts, nil
}

func (r *DraftRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DraftStatus, feedback string) error {
	query := `UPDATE drafts SET status = $1, feedback = $2, updated_at = $3 WHERE id = $4`
	tag, err := r.pool.Exec(ctx, query, status, feedback, time.Now(), id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
