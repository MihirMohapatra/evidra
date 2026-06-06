package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"github.com/evidra/evidra/identity/domain"
	idpg "github.com/evidra/evidra/identity/repository/postgres"
	idservice "github.com/evidra/evidra/identity/service"
)

type IdentitySuite struct {
	suite.Suite
	ctx      context.Context
	service  *idservice.IdentityService
	orgRepo  *idpg.OrganizationRepo
	userRepo *idpg.UserRepo
	sessRepo *idpg.SessionRepo
	keyRepo  *idpg.APIKeyRepo
}

func TestIdentitySuite(t *testing.T) {
	if identityPool == nil {
		t.Skip("identity postgres container not available")
	}
	suite.Run(t, new(IdentitySuite))
}

func (s *IdentitySuite) SetupSuite() {
	s.ctx = context.Background()

	s.orgRepo = idpg.NewOrganizationRepo(identityPool)
	s.userRepo = idpg.NewUserRepo(identityPool)
	s.sessRepo = idpg.NewSessionRepo(identityPool)
	s.keyRepo = idpg.NewAPIKeyRepo(identityPool)
	oidcRepo := idpg.NewOIDCStateRepo(identityPool)
	linkRepo := idpg.NewLinkedAccountRepo(identityPool)

	s.service = idservice.New(
		s.orgRepo,
		s.userRepo,
		s.sessRepo,
		s.keyRepo,
		oidcRepo,
		linkRepo,
		idservice.Config{
			JWTSecret:         "test-secret",
			JWTIssuer:         "evidra-test",
			SessionTTL:        1 * time.Hour,
			APIKeyLength:      32,
			PasswordMinLength: 8,
		},
		nil,
	)
}

func (s *IdentitySuite) TearDownTest() {
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM sessions")
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM linked_accounts")
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM oidc_states")
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM api_keys")
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM users")
	_, _ = identityPool.Exec(s.ctx, "DELETE FROM organizations")
}

func (s *IdentitySuite) TestCreateAndGetOrganization() {
	org, err := s.service.CreateOrganization(s.ctx, "Test Org", "testorg")
	s.Require().NoError(err)
	s.Require().NotNil(org)
	s.Equal("Test Org", org.Name)
	s.Equal("testorg", org.Slug)
	s.NotEqual(uuid.Nil, org.ID)

	got, err := s.service.GetOrganization(s.ctx, org.ID)
	s.Require().NoError(err)
	s.Equal(org.Name, got.Name)
	s.Equal(org.Slug, got.Slug)

	_, err = s.service.CreateOrganization(s.ctx, "Test Org", "testorg")
	s.Error(err)
	s.ErrorContains(err, "already exists")

	_, err = s.service.GetOrganization(s.ctx, uuid.New())
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *IdentitySuite) TestCreateAndGetUser() {
	org, err := s.service.CreateOrganization(s.ctx, "Users Org", "usersorg")
	s.Require().NoError(err)

	user, err := s.service.CreateUser(s.ctx, idservice.CreateUserInput{
		OrganizationID: org.ID,
		Email:          "user@example.com",
		Password:       "password123",
		Role:           "reviewer",
	})
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Equal("user@example.com", user.Email)
	s.Equal(domain.RoleReviewer, user.Role)
	s.True(user.IsActive)

	got, err := s.service.GetUser(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Equal(user.Email, got.Email)

	_, err = s.service.GetUser(s.ctx, uuid.New())
	s.ErrorIs(err, domain.ErrNotFound)

	_, err = s.service.CreateUser(s.ctx, idservice.CreateUserInput{
		OrganizationID: org.ID,
		Email:          "user@example.com",
		Password:       "password123",
		Role:           "reviewer",
	})
	s.Error(err)
	s.ErrorContains(err, "already exists")
}

func (s *IdentitySuite) TestLogin() {
	org, err := s.service.CreateOrganization(s.ctx, "Login Org", "loginorg")
	s.Require().NoError(err)

	_, err = s.service.CreateUser(s.ctx, idservice.CreateUserInput{
		OrganizationID: org.ID,
		Email:          "login@example.com",
		Password:       "securepass",
		Role:           "admin",
	})
	s.Require().NoError(err)

	session, err := s.service.Login(s.ctx, "login@example.com", "securepass")
	s.Require().NoError(err)
	s.Require().NotNil(session)
	s.NotEmpty(session.Token)
	s.NotEmpty(session.RefreshToken)
	s.False(session.IsExpired())

	_, err = s.service.Login(s.ctx, "login@example.com", "wrongpassword")
	s.ErrorIs(err, domain.ErrInvalidCredentials)

	_, err = s.service.Login(s.ctx, "nonexistent@example.com", "securepass")
	s.ErrorIs(err, domain.ErrInvalidCredentials)
}

func (s *IdentitySuite) TestCreateAndValidateAPIKey() {
	org, err := s.service.CreateOrganization(s.ctx, "APIKey Org", "apikeyorg")
	s.Require().NoError(err)

	key, rawKey, err := s.service.CreateAPIKey(s.ctx, org.ID, "test-key")
	s.Require().NoError(err)
	s.Require().NotNil(key)
	s.Equal("test-key", key.Name)
	s.True(key.IsActive)
	s.NotEmpty(rawKey)
	s.Len(rawKey, 64)

	validated, err := s.service.ValidateAPIKey(s.ctx, rawKey)
	s.Require().NoError(err)
	s.Equal(key.ID, validated.ID)

	err = s.service.RevokeAPIKey(s.ctx, key.ID)
	s.Require().NoError(err)

	_, err = s.service.ValidateAPIKey(s.ctx, rawKey)
	s.ErrorIs(err, domain.ErrUnauthorized)

	keys, err := s.service.ListAPIKeys(s.ctx, org.ID)
	s.Require().NoError(err)
	s.Len(keys, 1)
	s.False(keys[0].IsActive)
}

func TestPasswordHash(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte("password123"))
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte("wrong"))
	if err == nil {
		t.Fatal("expected mismatch")
	}
}
