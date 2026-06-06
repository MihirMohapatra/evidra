package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/evidra/evidra/evidence/domain"
	"github.com/evidra/evidra/evidence/service"
	"github.com/evidra/evidra/evidence/transport/dto"
)

type Handler struct {
	svc *service.EvidenceService
}

func NewHandler(svc *service.EvidenceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEvidenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid owner_id")
		return
	}

	cat, ok := domain.ParseCategory(req.Category)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category")
		return
	}

	var expiresAt time.Time
	if req.ExpiresAt != "" {
		expiresAt, err = time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid expires_at format (use RFC3339)")
			return
		}
	}

	input := domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  cat,
		OwnerID:   tenantID,
		SourceURL: req.SourceURL,
		Tags:      req.Tags,
		ExpiresAt: expiresAt,
	}

	item, err := h.svc.Create(r.Context(), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToEvidenceResponse(item))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	item, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := service.ListFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Tag:      r.URL.Query().Get("tag"),
		Search:   r.URL.Query().Get("search"),
		Expiring: r.URL.Query().Get("expiring") == "true",
		Expired:  r.URL.Query().Get("expired") == "true",
		Limit:    50,
		Offset:   0,
	}

	if tenantID := r.URL.Query().Get("tenant_id"); tenantID != "" {
		filter.TenantID, _ = uuid.Parse(tenantID)
	}
	if ownerID := r.URL.Query().Get("owner_id"); ownerID != "" {
		filter.OwnerID, _ = uuid.Parse(ownerID)
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := parseInt(limit); err == nil && l > 0 && l <= 200 {
			filter.Limit = l
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := parseInt(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	items, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.EvidenceResponse, len(items))
	for i, item := range items {
		resp[i] = dto.ToEvidenceResponse(item)
	}

	writeJSON(w, http.StatusOK, dto.PaginatedResponse{
		Items:  resp,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	var req dto.UpdateEvidenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := domain.CreateEvidenceInput{
		Title:     req.Title,
		Content:   req.Content,
		SourceURL: req.SourceURL,
		Tags:      req.Tags,
	}

	if req.Category != "" {
		cat, ok := domain.ParseCategory(req.Category)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid category")
			return
		}
		input.Category = cat
	}
	if req.OwnerID != "" {
		uid, err := uuid.Parse(req.OwnerID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid owner_id")
			return
		}
		input.OwnerID = uid
	}
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid expires_at format")
			return
		}
		input.ExpiresAt = t
	}

	item, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	item, err := h.svc.Submit(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	var req dto.ApproveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reviewerID, _ := uuid.Parse(r.Header.Get("X-User-ID"))

	item, err := h.svc.Approve(r.Context(), id, reviewerID, req.Comment)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	var req dto.RejectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reviewerID, _ := uuid.Parse(r.Header.Get("X-User-ID"))

	item, err := h.svc.Reject(r.Context(), id, reviewerID, req.Comment)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	item, err := h.svc.Export(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToEvidenceResponse(item))
}

func (h *Handler) GetApprovalHistory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}

	approvals, err := h.svc.GetApprovalHistory(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.ApprovalResponse, len(approvals))
	for i, a := range approvals {
		resp[i] = dto.ToApprovalResponse(a)
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- helpers ---

func decodeJSON(r *http.Request, v any) error {
	defer func() { _ = r.Body.Close() }()
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
	case errors.Is(err, domain.ErrInvalidTransition):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidCategory):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEvidenceExpired):
		writeError(w, http.StatusGone, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
