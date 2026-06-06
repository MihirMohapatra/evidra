package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/identity/domain"
)

type OrganizationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToOrganizationResponse(org *domain.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}

type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		Role:           string(user.Role),
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

type SessionResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserResponse `json:"user"`
}

func ToSessionResponse(session *domain.Session, user *domain.User) SessionResponse {
	return SessionResponse{
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		User:         ToUserResponse(user),
	}
}

type APIKeyResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	KeyPrefix      string    `json:"key_prefix"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToAPIKeyResponse(key *domain.APIKey) APIKeyResponse {
	return APIKeyResponse{
		ID:        key.ID,
		Name:      key.Name,
		KeyPrefix: key.KeyPrefix,
		IsActive:  key.IsActive,
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}
}

type CreateAPIKeyFullResponse struct {
	APIKeyResponse
	RawKey string `json:"raw_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
