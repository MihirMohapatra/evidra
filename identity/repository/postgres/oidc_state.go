package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type OIDCStateRepo struct {
	pool *pgxpool.Pool
}

func NewOIDCStateRepo(pool *pgxpool.Pool) *OIDCStateRepo {
	return &OIDCStateRepo{pool: pool}
}

func (r *OIDCStateRepo) Create(ctx context.Context, state *domain.OIDCState) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO oidc_states (id, provider, state, nonce, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		state.ID, state.Provider, state.State, state.Nonce, state.ExpiresAt, state.CreatedAt)
	if err != nil {
		return fmt.Errorf("create oidc state: %w", err)
	}
	return nil
}

func (r *OIDCStateRepo) GetByState(ctx context.Context, state string) (*domain.OIDCState, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, provider, state, nonce, expires_at, created_at
		FROM oidc_states WHERE state = $1`, state)

	var s domain.OIDCState
	if err := row.Scan(&s.ID, &s.Provider, &s.State, &s.Nonce, &s.ExpiresAt, &s.CreatedAt); err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get oidc state: %w", err)
	}
	return &s, nil
}

func (r *OIDCStateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM oidc_states WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete oidc state: %w", err)
	}
	return nil
}

func (r *OIDCStateRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM oidc_states WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("delete expired oidc states: %w", err)
	}
	return nil
}
