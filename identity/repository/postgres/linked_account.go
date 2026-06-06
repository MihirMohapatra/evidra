package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/identity/domain"
)

type LinkedAccountRepo struct {
	pool *pgxpool.Pool
}

func NewLinkedAccountRepo(pool *pgxpool.Pool) *LinkedAccountRepo {
	return &LinkedAccountRepo{pool: pool}
}

func (r *LinkedAccountRepo) Create(ctx context.Context, account *domain.LinkedAccount) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO linked_accounts (id, user_id, provider, subject, email, name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, account.UserID, account.Provider, account.Subject, account.Email, account.Name, account.CreatedAt)
	if err != nil {
		return fmt.Errorf("create linked account: %w", err)
	}
	return nil
}

func (r *LinkedAccountRepo) GetByProviderSubject(ctx context.Context, provider, subject string) (*domain.LinkedAccount, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, subject, email, name, created_at
		FROM linked_accounts WHERE provider = $1 AND subject = $2`, provider, subject)

	var a domain.LinkedAccount
	if err := row.Scan(&a.ID, &a.UserID, &a.Provider, &a.Subject, &a.Email, &a.Name, &a.CreatedAt); err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get linked account: %w", err)
	}
	return &a, nil
}

func (r *LinkedAccountRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.LinkedAccount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, provider, subject, email, name, created_at
		FROM linked_accounts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("list linked accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*domain.LinkedAccount
	for rows.Next() {
		var a domain.LinkedAccount
		if err := rows.Scan(&a.ID, &a.UserID, &a.Provider, &a.Subject, &a.Email, &a.Name, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan linked account: %w", err)
		}
		accounts = append(accounts, &a)
	}
	if accounts == nil {
		accounts = []*domain.LinkedAccount{}
	}
	return accounts, nil
}

func (r *LinkedAccountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM linked_accounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete linked account: %w", err)
	}
	return nil
}
