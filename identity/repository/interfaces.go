package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/evidra/evidra/identity/domain"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	List(ctx context.Context) ([]*domain.Organization, error)
	Update(ctx context.Context, org *domain.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (*domain.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type APIKeyRepository interface {
	Create(ctx context.Context, key *domain.APIKey) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error)
	GetByKeyHash(ctx context.Context, hash string) (*domain.APIKey, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.APIKey, error)
	Update(ctx context.Context, key *domain.APIKey) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OIDCStateRepository interface {
	Create(ctx context.Context, state *domain.OIDCState) error
	GetByState(ctx context.Context, state string) (*domain.OIDCState, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type LinkedAccountRepository interface {
	Create(ctx context.Context, account *domain.LinkedAccount) error
	GetByProviderSubject(ctx context.Context, provider, subject string) (*domain.LinkedAccount, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.LinkedAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
