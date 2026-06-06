package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type APIKeyRepo struct {
	pool *pgxpool.Pool
}

func NewAPIKeyRepo(pool *pgxpool.Pool) *APIKeyRepo {
	return &APIKeyRepo{pool: pool}
}

func (r *APIKeyRepo) Create(ctx context.Context, key *domain.APIKey) error {
	query := `INSERT INTO api_keys (id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, key.ID, key.OrganizationID, key.Name, key.KeyHash, key.KeyPrefix, key.IsActive, key.CreatedAt, key.UpdatedAt)
	return err
}

func (r *APIKeyRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error) {
	query := `SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at FROM api_keys WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	key := &domain.APIKey{}
	err := row.Scan(&key.ID, &key.OrganizationID, &key.Name, &key.KeyHash, &key.KeyPrefix, &key.IsActive, &key.CreatedAt, &key.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return key, err
}

func (r *APIKeyRepo) GetByKeyHash(ctx context.Context, hash string) (*domain.APIKey, error) {
	query := `SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at FROM api_keys WHERE key_hash = $1`
	row := r.pool.QueryRow(ctx, query, hash)

	key := &domain.APIKey{}
	err := row.Scan(&key.ID, &key.OrganizationID, &key.Name, &key.KeyHash, &key.KeyPrefix, &key.IsActive, &key.CreatedAt, &key.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return key, err
}

func (r *APIKeyRepo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.APIKey, error) {
	query := `SELECT id, organization_id, name, key_hash, key_prefix, is_active, created_at, updated_at FROM api_keys WHERE organization_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*domain.APIKey
	for rows.Next() {
		key := &domain.APIKey{}
		if err := rows.Scan(&key.ID, &key.OrganizationID, &key.Name, &key.KeyHash, &key.KeyPrefix, &key.IsActive, &key.CreatedAt, &key.UpdatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (r *APIKeyRepo) Update(ctx context.Context, key *domain.APIKey) error {
	query := `UPDATE api_keys SET name = $1, key_hash = $2, key_prefix = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	tag, err := r.pool.Exec(ctx, query, key.Name, key.KeyHash, key.KeyPrefix, key.IsActive, key.UpdatedAt, key.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *APIKeyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
