package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/questionnaire/domain"
)

type QuestionnaireRepo struct {
	pool *pgxpool.Pool
}

func NewQuestionnaireRepo(pool *pgxpool.Pool) *QuestionnaireRepo {
	return &QuestionnaireRepo{pool: pool}
}

func (r *QuestionnaireRepo) Create(ctx context.Context, q *domain.Questionnaire) error {
	query := `INSERT INTO questionnaires (id, tenant_id, title, file_name, file_url, file_type, file_size, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, q.ID, q.TenantID, q.Title, q.FileName, q.FileURL, q.FileType, q.FileSize, q.Status, q.Version, q.CreatedAt, q.UpdatedAt)
	return err
}

func (r *QuestionnaireRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Questionnaire, error) {
	query := `SELECT id, tenant_id, title, file_name, file_url, file_type, file_size, status, version, created_at, updated_at FROM questionnaires WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	q := &domain.Questionnaire{}
	err := row.Scan(&q.ID, &q.TenantID, &q.Title, &q.FileName, &q.FileURL, &q.FileType, &q.FileSize, &q.Status, &q.Version, &q.CreatedAt, &q.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return q, err
}

func (r *QuestionnaireRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Questionnaire, error) {
	query := `SELECT id, tenant_id, title, file_name, file_url, file_type, file_size, status, version, created_at, updated_at FROM questionnaires WHERE tenant_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []*domain.Questionnaire
	for rows.Next() {
		q := &domain.Questionnaire{}
		if err := rows.Scan(&q.ID, &q.TenantID, &q.Title, &q.FileName, &q.FileURL, &q.FileType, &q.FileSize, &q.Status, &q.Version, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}
	return qs, rows.Err()
}

func (r *QuestionnaireRepo) Update(ctx context.Context, q *domain.Questionnaire) error {
	query := `UPDATE questionnaires SET title = $1, file_name = $2, file_url = $3, file_type = $4, file_size = $5, status = $6, version = $7, updated_at = $8 WHERE id = $9`
	tag, err := r.pool.Exec(ctx, query, q.Title, q.FileName, q.FileURL, q.FileType, q.FileSize, q.Status, q.Version, q.UpdatedAt, q.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *QuestionnaireRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM questionnaires WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
