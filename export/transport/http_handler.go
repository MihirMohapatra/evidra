package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/evidra/evidra/export/domain"
	"github.com/evidra/evidra/export/service"
	"github.com/evidra/evidra/export/transport/dto"
)

type Handler struct {
	svc *service.ExportService
}

func NewHandler(svc *service.ExportService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	var req dto.ExportRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	evidenceID, err := uuid.Parse(req.EvidenceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence_id")
		return
	}

	exp, err := h.svc.Export(r.Context(), service.ExportInput{
		TenantID:    uuid.Nil,
		EvidenceID:  evidenceID,
		RequesterID: uuid.Nil,
		Format:      req.Format,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToExportResponse(exp))
}

func (h *Handler) GetExport(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid export id")
		return
	}

	exp, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToExportResponse(exp))
}

func (h *Handler) ListExports(w http.ResponseWriter, r *http.Request) {
	evidenceIDStr := r.URL.Query().Get("evidence_id")
	if evidenceIDStr != "" {
		evidenceID, err := uuid.Parse(evidenceIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid evidence_id")
			return
		}
		exports, err := h.svc.ListByEvidence(r.Context(), evidenceID)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		resp := make([]dto.ExportResponse, len(exports))
		for i, e := range exports {
			resp[i] = dto.ToExportResponse(e)
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	writeJSON(w, http.StatusOK, []dto.ExportResponse{})
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
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidFormat):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEvidenceNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrExportFailed):
		writeError(w, http.StatusInternalServerError, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
