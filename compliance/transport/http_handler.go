package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/evidra/evidra/compliance/domain"
	"github.com/evidra/evidra/compliance/service"
	"github.com/evidra/evidra/compliance/transport/dto"
)

type Handler struct {
	svc *service.ComplianceService
}

func NewHandler(svc *service.ComplianceService) *Handler {
	return &Handler{svc: svc}
}

// --- Frameworks ---

func (h *Handler) CreateFramework(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFrameworkRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	f, err := h.svc.CreateFramework(r.Context(), req.Name, req.Slug, req.Description, req.Version)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToFrameworkResponse(f))
}

func (h *Handler) GetFramework(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid framework id")
		return
	}

	f, err := h.svc.GetFramework(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToFrameworkResponse(f))
}

func (h *Handler) ListFrameworks(w http.ResponseWriter, r *http.Request) {
	frameworks, err := h.svc.ListFrameworks(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.FrameworkResponse, len(frameworks))
	for i, f := range frameworks {
		resp[i] = dto.ToFrameworkResponse(f)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteFramework(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid framework id")
		return
	}

	if err := h.svc.DeleteFramework(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Controls ---

func (h *Handler) CreateControl(w http.ResponseWriter, r *http.Request) {
	frameworkID, err := uuid.Parse(chi.URLParam(r, "frameworkId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid framework id")
		return
	}

	var req dto.CreateControlRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.svc.CreateControl(r.Context(), service.CreateControlInput{
		FrameworkID:     frameworkID,
		ControlID:       req.ControlID,
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		RiskDescription: req.RiskDescription,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToControlResponse(c))
}

func (h *Handler) ListControls(w http.ResponseWriter, r *http.Request) {
	frameworkID, err := uuid.Parse(chi.URLParam(r, "frameworkId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid framework id")
		return
	}

	controls, err := h.svc.ListControls(r.Context(), frameworkID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.ControlResponse, len(controls))
	for i, c := range controls {
		resp[i] = dto.ToControlResponse(c)
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Mappings ---

func (h *Handler) MapEvidence(w http.ResponseWriter, r *http.Request) {
	var req dto.MapEvidenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	evidenceID, err := uuid.Parse(req.EvidenceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence_id")
		return
	}
	controlID, err := uuid.Parse(req.ControlID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid control_id")
		return
	}

	m, err := h.svc.MapEvidence(r.Context(), service.MapEvidenceInput{
		TenantID:   uuid.Nil,
		EvidenceID: evidenceID,
		ControlID:  controlID,
		Notes:      req.Notes,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToEvidenceMappingResponse(m))
}

func (h *Handler) UnmapEvidence(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid mapping id")
		return
	}

	if err := h.svc.UnmapEvidence(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListMappingsByControl(w http.ResponseWriter, r *http.Request) {
	controlID, err := uuid.Parse(chi.URLParam(r, "controlId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid control_id")
		return
	}

	mappings, err := h.svc.ListMappingsByControl(r.Context(), controlID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	resp := make([]dto.EvidenceMappingResponse, len(mappings))
	for i, m := range mappings {
		resp[i] = dto.ToEvidenceMappingResponse(m)
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Coverage ---

func (h *Handler) GetFrameworkCoverage(w http.ResponseWriter, r *http.Request) {
	frameworkID, err := uuid.Parse(chi.URLParam(r, "frameworkId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid framework id")
		return
	}

	coverage, err := h.svc.GetFrameworkCoverage(r.Context(), frameworkID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToFrameworkCoverageResponse(coverage))
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
	case errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrMappingExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrFrameworkInUse):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
