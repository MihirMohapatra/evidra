package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/audit/domain"
	"github.com/evidra/evidra/audit/repository"
	"github.com/evidra/evidra/audit/service"
	"github.com/evidra/evidra/audit/transport/dto"
)

type Handler struct {
	svc *service.AuditService
}

func NewHandler(svc *service.AuditService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Record(w http.ResponseWriter, r *http.Request) {
	var req dto.RecordAuditRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tenant_id")
		return
	}
	actorID, err := uuid.Parse(req.ActorID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid actor_id")
		return
	}

	event, err := h.svc.Record(r.Context(), service.RecordInput{
		TenantID: tenantID,
		ActorID:  actorID,
		Action:   domain.Action(req.Action),
		TargetID: req.TargetID,
		Metadata: req.Metadata,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToAuditEventResponse(event))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := repository.AuditFilter{
		Limit:  50,
		Offset: 0,
	}

	if tenantID := r.URL.Query().Get("tenant_id"); tenantID != "" {
		filter.TenantID, _ = uuid.Parse(tenantID)
	}
	if actorID := r.URL.Query().Get("actor_id"); actorID != "" {
		filter.ActorID, _ = uuid.Parse(actorID)
	}
	if action := r.URL.Query().Get("action"); action != "" {
		filter.Action = domain.Action(action)
	}
	if targetID := r.URL.Query().Get("target_id"); targetID != "" {
		filter.TargetID = targetID
	}
	if since := r.URL.Query().Get("since"); since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			filter.Since = t
		}
	}
	if until := r.URL.Query().Get("until"); until != "" {
		if t, err := time.Parse(time.RFC3339, until); err == nil {
			filter.Until = t
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 200 {
			filter.Limit = l
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	events, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.AuditEventResponse, len(events))
	for i, e := range events {
		resp[i] = dto.ToAuditEventResponse(e)
	}

	writeJSON(w, http.StatusOK, dto.PaginatedResponse{
		Items:  resp,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

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
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
