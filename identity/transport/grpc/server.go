package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/identity/domain"
	"github.com/evidra/evidra/identity/service"
)

type IdentityServer struct {
	evidrav1.UnimplementedIdentityServiceServer
	svc *service.IdentityService
}

func NewServer(svc *service.IdentityService) *IdentityServer {
	return &IdentityServer{svc: svc}
}

func (s *IdentityServer) CreateOrganization(ctx context.Context, req *evidrav1.CreateOrganizationRequest) (*evidrav1.Organization, error) {
	org, err := s.svc.CreateOrganization(ctx, req.Name, req.Slug)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return organizationToProto(org), nil
}

func (s *IdentityServer) GetOrganization(ctx context.Context, req *evidrav1.GetOrganizationRequest) (*evidrav1.Organization, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	org, err := s.svc.GetOrganization(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return organizationToProto(org), nil
}

func (s *IdentityServer) ListOrganizations(ctx context.Context, _ *evidrav1.Empty) (*evidrav1.ListOrganizationsResponse, error) {
	orgs, err := s.svc.ListOrganizations(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.Organization, len(orgs))
	for i, o := range orgs {
		proto[i] = organizationToProto(o)
	}
	return &evidrav1.ListOrganizationsResponse{Organizations: proto}, nil
}

func (s *IdentityServer) UpdateOrganization(ctx context.Context, req *evidrav1.UpdateOrganizationRequest) (*evidrav1.Organization, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	org, err := s.svc.UpdateOrganization(ctx, id, req.Name)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return organizationToProto(org), nil
}

func (s *IdentityServer) DeleteOrganization(ctx context.Context, req *evidrav1.DeleteOrganizationRequest) (*evidrav1.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.DeleteOrganization(ctx, id); err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.Empty{}, nil
}

func (s *IdentityServer) CreateUser(ctx context.Context, req *evidrav1.CreateUserRequest) (*evidrav1.User, error) {
	orgID, err := uuid.Parse(req.OrganizationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization_id")
	}
	user, err := s.svc.CreateUser(ctx, service.CreateUserInput{
		OrganizationID: orgID,
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return userToProto(user), nil
}

func (s *IdentityServer) GetUser(ctx context.Context, req *evidrav1.GetUserRequest) (*evidrav1.User, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	user, err := s.svc.GetUser(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return userToProto(user), nil
}

func (s *IdentityServer) ListUsers(ctx context.Context, req *evidrav1.ListUsersRequest) (*evidrav1.ListUsersResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization_id")
	}
	users, err := s.svc.ListUsers(ctx, orgID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.User, len(users))
	for i, u := range users {
		proto[i] = userToProto(u)
	}
	return &evidrav1.ListUsersResponse{Users: proto}, nil
}

func (s *IdentityServer) UpdateUser(ctx context.Context, req *evidrav1.UpdateUserRequest) (*evidrav1.User, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	user, err := s.svc.UpdateUser(ctx, id, service.UpdateUserInput{Email: req.Email, Role: req.Role})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return userToProto(user), nil
}

func (s *IdentityServer) DeleteUser(ctx context.Context, req *evidrav1.DeleteUserRequest) (*evidrav1.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.DeleteUser(ctx, id); err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.Empty{}, nil
}

func (s *IdentityServer) Login(ctx context.Context, req *evidrav1.LoginRequest) (*evidrav1.Session, error) {
	session, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, toGRPCError(err)
	}
	user, _ := s.svc.GetUser(ctx, session.UserID)
	return sessionToProto(session, user), nil
}

func (s *IdentityServer) RefreshToken(ctx context.Context, req *evidrav1.RefreshTokenRequest) (*evidrav1.Session, error) {
	session, err := s.svc.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		return nil, toGRPCError(err)
	}
	user, _ := s.svc.GetUser(ctx, session.UserID)
	return sessionToProto(session, user), nil
}

func (s *IdentityServer) Logout(ctx context.Context, req *evidrav1.LogoutRequest) (*evidrav1.Empty, error) {
	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session_id")
	}
	if err := s.svc.Logout(ctx, sessionID); err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.Empty{}, nil
}

func (s *IdentityServer) CreateAPIKey(ctx context.Context, req *evidrav1.CreateAPIKeyRequest) (*evidrav1.CreateAPIKeyFullResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization_id")
	}
	key, rawKey, err := s.svc.CreateAPIKey(ctx, orgID, req.Name)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.CreateAPIKeyFullResponse{
		Key:    apiKeyToProto(key),
		RawKey: rawKey,
	}, nil
}

func (s *IdentityServer) ListAPIKeys(ctx context.Context, req *evidrav1.ListAPIKeysRequest) (*evidrav1.ListAPIKeysResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization_id")
	}
	keys, err := s.svc.ListAPIKeys(ctx, orgID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.APIKeyResponse, len(keys))
	for i, k := range keys {
		proto[i] = apiKeyToProto(k)
	}
	return &evidrav1.ListAPIKeysResponse{Keys: proto}, nil
}

func (s *IdentityServer) RevokeAPIKey(ctx context.Context, req *evidrav1.RevokeAPIKeyRequest) (*evidrav1.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.RevokeAPIKey(ctx, id); err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.Empty{}, nil
}

func (s *IdentityServer) ListOIDCProviders(ctx context.Context, _ *evidrav1.Empty) (*evidrav1.ListOIDCProvidersResponse, error) {
	providers := s.svc.GetOIDCProviders()
	proto := make([]*evidrav1.OIDCProviderInfo, len(providers))
	for i, p := range providers {
		proto[i] = &evidrav1.OIDCProviderInfo{Name: p.Name, RedirectUrl: p.RedirectURL}
	}
	return &evidrav1.ListOIDCProvidersResponse{Providers: proto}, nil
}

func (s *IdentityServer) OIDCLogin(ctx context.Context, req *evidrav1.OIDCLoginRequest) (*evidrav1.OIDCLoginResponse, error) {
	authURL, err := s.svc.InitiateOIDCLogin(ctx, req.Provider)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.OIDCLoginResponse{AuthUrl: authURL}, nil
}

func (s *IdentityServer) OIDCCallback(ctx context.Context, req *evidrav1.OIDCCallbackRequest) (*evidrav1.Session, error) {
	session, err := s.svc.HandleOIDCCallback(ctx, req.Provider, req.Code, req.State)
	if err != nil {
		return nil, toGRPCError(err)
	}
	user, _ := s.svc.GetUser(ctx, session.UserID)
	return sessionToProto(session, user), nil
}

// --- converters ---

func organizationToProto(o *domain.Organization) *evidrav1.Organization {
	return &evidrav1.Organization{
		Id:        o.ID.String(),
		Name:      o.Name,
		Slug:      o.Slug,
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}

func userToProto(u *domain.User) *evidrav1.User {
	return &evidrav1.User{
		Id:             u.ID.String(),
		OrganizationId: u.OrganizationID.String(),
		Email:          u.Email,
		Role:           string(u.Role),
		IsActive:       u.IsActive,
		CreatedAt:      timestamppb.New(u.CreatedAt),
		UpdatedAt:      timestamppb.New(u.UpdatedAt),
	}
}

func sessionToProto(s *domain.Session, u *domain.User) *evidrav1.Session {
	return &evidrav1.Session{
		Token:        s.Token,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    timestamppb.New(s.ExpiresAt),
		User:         userToProto(u),
	}
}

func apiKeyToProto(k *domain.APIKey) *evidrav1.APIKeyResponse {
	return &evidrav1.APIKeyResponse{
		Id:        k.ID.String(),
		Name:      k.Name,
		KeyPrefix: k.KeyPrefix,
		IsActive:  k.IsActive,
		CreatedAt: timestamppb.New(k.CreatedAt),
		UpdatedAt: timestamppb.New(k.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidRole):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrForbidden), errors.Is(err, domain.ErrUserInactive):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
