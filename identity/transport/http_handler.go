package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/evidra/evidra/identity/domain"
	"github.com/evidra/evidra/identity/service"
	"github.com/evidra/evidra/identity/transport/dto"
)

type Handler struct {
	svc *service.IdentityService
}

func NewHandler(svc *service.IdentityService) *Handler {
	return &Handler{svc: svc}
}

// --- Organizations ---

func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrganizationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	org, err := h.svc.CreateOrganization(r.Context(), req.Name, req.Slug)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToOrganizationResponse(org))
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	org, err := h.svc.GetOrganization(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToOrganizationResponse(org))
}

func (h *Handler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.svc.ListOrganizations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list organizations")
		return
	}

	resp := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		resp[i] = dto.ToOrganizationResponse(org)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	org, err := h.svc.UpdateOrganization(r.Context(), id, req.Name)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToOrganizationResponse(org))
}

func (h *Handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	if err := h.svc.DeleteOrganization(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Users ---

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	input := service.CreateUserInput{
		OrganizationID: orgID,
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
	}

	user, err := h.svc.CreateUser(r.Context(), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToUserResponse(user))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToUserResponse(user))
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		writeError(w, http.StatusBadRequest, "organization_id query parameter is required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization_id")
		return
	}

	users, err := h.svc.ListUsers(r.Context(), orgID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = dto.ToUserResponse(u)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := service.UpdateUserInput{
		Email: req.Email,
		Role:  req.Role,
	}

	user, err := h.svc.UpdateUser(r.Context(), id, input)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToUserResponse(user))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Auth ---

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	user, _ := h.svc.GetUser(r.Context(), session.UserID)
	writeJSON(w, http.StatusOK, dto.ToSessionResponse(session, user))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("session_id")
	if sessionIDStr == "" {
		writeError(w, http.StatusBadRequest, "session_id is required")
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session_id")
		return
	}

	if err := h.svc.Logout(r.Context(), sessionID); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.svc.RefreshSession(r.Context(), req.RefreshToken)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	user, _ := h.svc.GetUser(r.Context(), session.UserID)
	writeJSON(w, http.StatusOK, dto.ToSessionResponse(session, user))
}

// --- OIDC ---

func (h *Handler) ListOIDCProviders(w http.ResponseWriter, r *http.Request) {
	providers := h.svc.GetOIDCProviders()
	writeJSON(w, http.StatusOK, providers)
}

func (h *Handler) OIDCLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "provider is required")
		return
	}

	authURL, err := h.svc.InitiateOIDCLogin(r.Context(), provider)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "provider is required")
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		writeError(w, http.StatusBadRequest, "missing code or state")
		return
	}

	session, err := h.svc.HandleOIDCCallback(r.Context(), provider, code, state)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	user, _ := h.svc.GetUser(r.Context(), session.UserID)
	writeJSON(w, http.StatusOK, dto.ToSessionResponse(session, user))
}

// --- API Keys ---

func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		writeError(w, http.StatusBadRequest, "organization_id is required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization_id")
		return
	}

	var req dto.CreateAPIKeyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	key, rawKey, err := h.svc.CreateAPIKey(r.Context(), orgID, req.Name)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateAPIKeyFullResponse{
		APIKeyResponse: dto.ToAPIKeyResponse(key),
		RawKey:         rawKey,
	})
}

func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		writeError(w, http.StatusBadRequest, "organization_id is required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization_id")
		return
	}

	keys, err := h.svc.ListAPIKeys(r.Context(), orgID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.APIKeyResponse, len(keys))
	for i, k := range keys {
		resp[i] = dto.ToAPIKeyResponse(k)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid api key id")
		return
	}

	if err := h.svc.RevokeAPIKey(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, dto.ErrorResponse{Error: msg})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidRole):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrSessionExpired):
		writeError(w, http.StatusUnauthorized, "session expired")
	case errors.Is(err, domain.ErrUserInactive):
		writeError(w, http.StatusForbidden, "user is inactive")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
