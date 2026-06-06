package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, organization_id, email, password_hash, role, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return user, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return user, err
}

func (r *UserRepo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.User, error) {
	query := `SELECT id, organization_id, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE organization_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, role = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	tag, err := r.pool.Exec(ctx, query, user.Email, user.PasswordHash, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
