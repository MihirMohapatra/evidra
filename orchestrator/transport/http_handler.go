package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/service"
	"github.com/evidra/evidra/orchestrator/transport/dto"
)

type Handler struct {
	svc *service.OrchestratorService
}

func NewHandler(svc *service.OrchestratorService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Answer(w http.ResponseWriter, r *http.Request) {
	var req dto.AnswerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id is required")
		return
	}
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tenant_id")
		return
	}

	svcReq := service.AnswerRequest{
		Question: req.Question,
		Context:  req.Context,
		TenantID: tenantID,
	}

	result, err := h.svc.Answer(r.Context(), svcReq)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	evResp := make([]dto.EvidenceResponse, len(result.Evidence))
	for i, ev := range result.Evidence {
		evResp[i] = dto.ToEvidenceResponse(ev)
	}

	writeJSON(w, http.StatusOK, dto.AnswerResponse{
		Draft:    dto.ToDraftResponse(result.Draft),
		Evidence: evResp,
	})
}

func (h *Handler) GetDraft(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	draft, err := h.svc.GetDraft(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToDraftResponse(draft))
}

func (h *Handler) ListDrafts(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	drafts, err := h.svc.ListDrafts(r.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.DraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = dto.ToDraftResponse(d)
	}

	writeJSON(w, http.StatusOK, dto.ListDraftsResponse{
		Items:  resp,
		Total:  len(resp),
		Limit:  limit,
		Offset: offset,
	})
}

func (h *Handler) ApproveDraft(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if err := h.svc.ApproveDraft(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RejectDraft(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	var req dto.RejectDraftRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.RejectDraft(r.Context(), id, req.Feedback); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

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
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmbeddingFailed):
		writeError(w, http.StatusInternalServerError, "embedding generation failed")
	case errors.Is(err, domain.ErrLLMError):
		writeError(w, http.StatusInternalServerError, "llm service error")
	case errors.Is(err, domain.ErrValidationFailed):
		writeError(w, http.StatusUnprocessableEntity, "response validation failed")
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
