package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/questionnaire/domain"
)

type QuestionRepo struct {
	pool *pgxpool.Pool
}

func NewQuestionRepo(pool *pgxpool.Pool) *QuestionRepo {
	return &QuestionRepo{pool: pool}
}

func (r *QuestionRepo) Create(ctx context.Context, q *domain.Question) error {
	query := `INSERT INTO questions (id, questionnaire_id, text, type, "order", options, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, q.ID, q.QuestionnaireID, q.Text, q.Type, q.Order, q.Options, q.CreatedAt, q.UpdatedAt)
	return err
}

func (r *QuestionRepo) CreateBatch(ctx context.Context, questions []*domain.Question) error {
	if len(questions) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `INSERT INTO questions (id, questionnaire_id, text, type, "order", options, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	for _, q := range questions {
		_, err := tx.Exec(ctx, query, q.ID, q.QuestionnaireID, q.Text, q.Type, q.Order, q.Options, q.CreatedAt, q.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *QuestionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	query := `SELECT id, questionnaire_id, text, type, "order", options, created_at, updated_at FROM questions WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	q := &domain.Question{}
	err := row.Scan(&q.ID, &q.QuestionnaireID, &q.Text, &q.Type, &q.Order, &q.Options, &q.CreatedAt, &q.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return q, err
}

func (r *QuestionRepo) ListByQuestionnaire(ctx context.Context, questionnaireID uuid.UUID) ([]*domain.Question, error) {
	query := `SELECT id, questionnaire_id, text, type, "order", options, created_at, updated_at FROM questions WHERE questionnaire_id = $1 ORDER BY "order" ASC`
	rows, err := r.pool.Query(ctx, query, questionnaireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []*domain.Question
	for rows.Next() {
		q := &domain.Question{}
		if err := rows.Scan(&q.ID, &q.QuestionnaireID, &q.Text, &q.Type, &q.Order, &q.Options, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}
	return qs, rows.Err()
}

func (r *QuestionRepo) Update(ctx context.Context, q *domain.Question) error {
	query := `UPDATE questions SET text = $1, type = $2, "order" = $3, options = $4, updated_at = $5 WHERE id = $6`
	tag, err := r.pool.Exec(ctx, query, q.Text, q.Type, q.Order, q.Options, q.UpdatedAt, q.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *QuestionRepo) DeleteByQuestionnaire(ctx context.Context, questionnaireID uuid.UUID) error {
	query := `DELETE FROM questions WHERE questionnaire_id = $1`
	_, err := r.pool.Exec(ctx, query, questionnaireID)
	return err
}
