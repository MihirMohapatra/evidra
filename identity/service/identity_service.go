package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/evidra/evidra/identity/domain"
	"github.com/evidra/evidra/identity/repository"
)

type Config struct {
	JWTSecret         string
	JWTIssuer         string
	SessionTTL        time.Duration
	APIKeyLength      int
	PasswordMinLength int
}

type IdentityService struct {
	orgs  repository.OrganizationRepository
	users repository.UserRepository
	sess  repository.SessionRepository
	keys  repository.APIKeyRepository
	v     *validator.Validate
	cfg   Config
}

func New(
	orgs repository.OrganizationRepository,
	users repository.UserRepository,
	sess repository.SessionRepository,
	keys repository.APIKeyRepository,
	cfg Config,
) *IdentityService {
	return &IdentityService{
		orgs:  orgs,
		users: users,
		sess:  sess,
		keys:  keys,
		v:     validator.New(),
		cfg:   cfg,
	}
}

// --- Organizations ---

func (s *IdentityService) CreateOrganization(ctx context.Context, name, slug string) (*domain.Organization, error) {
	if err := s.v.Var(name, "required,max=255"); err != nil {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if err := s.v.Var(slug, "required,alphanum,max=100"); err != nil {
		return nil, fmt.Errorf("%w: slug must be alphanumeric", domain.ErrInvalidInput)
	}

	existing, err := s.orgs.GetBySlug(ctx, slug)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: organization with slug %q already exists", domain.ErrAlreadyExists, slug)
	}

	org := domain.NewOrganization(name, slug)
	if err := s.orgs.Create(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *IdentityService) GetOrganization(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	return s.orgs.GetByID(ctx, id)
}

func (s *IdentityService) ListOrganizations(ctx context.Context) ([]*domain.Organization, error) {
	return s.orgs.List(ctx)
}

func (s *IdentityService) UpdateOrganization(ctx context.Context, id uuid.UUID, name string) (*domain.Organization, error) {
	org, err := s.orgs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	org.Name = name
	org.UpdatedAt = time.Now()
	if err := s.orgs.Update(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *IdentityService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return s.orgs.Delete(ctx, id)
}

// --- Users ---

type CreateUserInput struct {
	OrganizationID uuid.UUID `validate:"required"`
	Email          string    `validate:"required,email,max=255"`
	Password       string    `validate:"required,min=8"`
	Role           string    `validate:"required"`
}

func (s *IdentityService) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	if err := s.v.Struct(input); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}

	role, ok := domain.ParseRole(input.Role)
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidRole, input.Role)
	}

	if _, err := s.orgs.GetByID(ctx, input.OrganizationID); err != nil {
		return nil, fmt.Errorf("organization: %w", err)
	}

	existing, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: user with email %q already exists", domain.ErrAlreadyExists, input.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(input.OrganizationID, input.Email, string(hash), role)
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *IdentityService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *IdentityService) ListUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.User, error) {
	return s.users.ListByOrganization(ctx, orgID)
}

type UpdateUserInput struct {
	Email string `validate:"omitempty,email"`
	Role  string `validate:"omitempty"`
}

func (s *IdentityService) UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Email != "" {
		if err := s.v.Var(input.Email, "email"); err != nil {
			return nil, fmt.Errorf("%w: invalid email", domain.ErrInvalidInput)
		}
		user.Email = input.Email
	}
	if input.Role != "" {
		role, ok := domain.ParseRole(input.Role)
		if !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrInvalidRole, input.Role)
		}
		user.Role = role
	}

	user.UpdatedAt = time.Now()
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *IdentityService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.users.Delete(ctx, id)
}

// --- Auth ---

func (s *IdentityService) Login(ctx context.Context, email, password string) (*domain.Session, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	session := domain.NewSession(user.ID, token, refreshToken, s.cfg.SessionTTL)
	if err := s.sess.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *IdentityService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.sess.Delete(ctx, sessionID)
}

func (s *IdentityService) ValidateSession(ctx context.Context, token string) (*domain.User, error) {
	claims, err := s.parseJWT(token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	session, err := s.sess.GetByToken(ctx, token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if session.IsExpired() {
		return nil, domain.ErrSessionExpired
	}

	sub, _ := claims["sub"].(string)
	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	return user, nil
}

func (s *IdentityService) RefreshSession(ctx context.Context, refreshToken string) (*domain.Session, error) {
	session, err := s.sess.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if session.IsExpired() {
		return nil, domain.ErrSessionExpired
	}

	if err := s.sess.Delete(ctx, session.ID); err != nil {
		return nil, err
	}

	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	token, newRefreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	newSession := domain.NewSession(user.ID, token, newRefreshToken, s.cfg.SessionTTL)
	if err := s.sess.Create(ctx, newSession); err != nil {
		return nil, err
	}

	return newSession, nil
}

// --- API Keys ---

func (s *IdentityService) CreateAPIKey(ctx context.Context, orgID uuid.UUID, name string) (*domain.APIKey, string, error) {
	if err := s.v.Var(name, "required,max=255"); err != nil {
		return nil, "", fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}

	rawKey := generateRawKey(s.cfg.APIKeyLength)
	hash := sha256Hex(rawKey)
	prefix := rawKey[:8]

	key := domain.NewAPIKey(orgID, name, hash, prefix)
	if err := s.keys.Create(ctx, key); err != nil {
		return nil, "", err
	}

	return key, rawKey, nil
}

func (s *IdentityService) ValidateAPIKey(ctx context.Context, rawKey string) (*domain.APIKey, error) {
	hash := sha256Hex(rawKey)
	key, err := s.keys.GetByKeyHash(ctx, hash)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if !key.IsActive {
		return nil, domain.ErrUnauthorized
	}
	return key, nil
}

func (s *IdentityService) RevokeAPIKey(ctx context.Context, id uuid.UUID) error {
	key, err := s.keys.GetByID(ctx, id)
	if err != nil {
		return err
	}
	key.IsActive = false
	key.UpdatedAt = time.Now()
	return s.keys.Update(ctx, key)
}

func (s *IdentityService) ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]*domain.APIKey, error) {
	return s.keys.ListByOrganization(ctx, orgID)
}

// --- internal helpers ---

func (s *IdentityService) generateTokenPair(user *domain.User) (token string, refreshToken string, err error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"org": user.OrganizationID.String(),
		"role": string(user.Role),
		"iat": now.Unix(),
		"exp": now.Add(s.cfg.SessionTTL).Unix(),
		"iss": s.cfg.JWTIssuer,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", err
	}
	refreshToken = hex.EncodeToString(refreshBytes)

	return token, refreshToken, nil
}

func (s *IdentityService) parseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func generateRawKey(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
