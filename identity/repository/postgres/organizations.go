package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type OrganizationRepo struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepo(pool *pgxpool.Pool) *OrganizationRepo {
	return &OrganizationRepo{pool: pool}
}

func (r *OrganizationRepo) Create(ctx context.Context, org *domain.Organization) error {
	query := `INSERT INTO organizations (id, name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, org.ID, org.Name, org.Slug, org.CreatedAt, org.UpdatedAt)
	return err
}

func (r *OrganizationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM organizations WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	org := &domain.Organization{}
	err := row.Scan(&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return org, err
}

func (r *OrganizationRepo) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM organizations WHERE slug = $1`
	row := r.pool.QueryRow(ctx, query, slug)

	org := &domain.Organization{}
	err := row.Scan(&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return org, err
}

func (r *OrganizationRepo) List(ctx context.Context) ([]*domain.Organization, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM organizations ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		org := &domain.Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, rows.Err()
}

func (r *OrganizationRepo) Update(ctx context.Context, org *domain.Organization) error {
	query := `UPDATE organizations SET name = $1, slug = $2, updated_at = $3 WHERE id = $4`
	tag, err := r.pool.Exec(ctx, query, org.Name, org.Slug, org.UpdatedAt, org.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *OrganizationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
