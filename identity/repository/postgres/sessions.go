package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

func (r *SessionRepo) Create(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, token, refresh_token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, session.ID, session.UserID, session.Token, session.RefreshToken, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	query := `SELECT id, user_id, token, refresh_token, expires_at, created_at FROM sessions WHERE token = $1`
	row := r.pool.QueryRow(ctx, query, token)

	session := &domain.Session{}
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return session, err
}

func (r *SessionRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	query := `SELECT id, user_id, token, refresh_token, expires_at, created_at FROM sessions WHERE refresh_token = $1`
	row := r.pool.QueryRow(ctx, query, refreshToken)

	session := &domain.Session{}
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return session, err
}

func (r *SessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SessionRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *SessionRepo) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.pool.Exec(ctx, query)
	return err
}
