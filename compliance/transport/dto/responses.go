package dto

import (
	"time"

	"github.com/evidra/evidra/compliance/domain"
)

type FrameworkResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Version     string `json:"version"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ToFrameworkResponse(f *domain.Framework) FrameworkResponse {
	return FrameworkResponse{
		ID:          f.ID.String(),
		Name:        f.Name,
		Slug:        f.Slug,
		Description: f.Description,
		Version:     f.Version,
		CreatedAt:   f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   f.UpdatedAt.Format(time.RFC3339),
	}
}

type ControlResponse struct {
	ID              string `json:"id"`
	FrameworkID     string `json:"framework_id"`
	ControlID       string `json:"control_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	RiskDescription string `json:"risk_description"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func ToControlResponse(c *domain.Control) ControlResponse {
	return ControlResponse{
		ID:              c.ID.String(),
		FrameworkID:     c.FrameworkID.String(),
		ControlID:       c.ControlID,
		Name:            c.Name,
		Description:     c.Description,
		Category:        string(c.Category),
		RiskDescription: c.RiskDescription,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}
}

type EvidenceMappingResponse struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EvidenceID string `json:"evidence_id"`
	ControlID  string `json:"control_id"`
	Notes      string `json:"notes"`
	MappedBy   string `json:"mapped_by"`
	CreatedAt  string `json:"created_at"`
}

func ToEvidenceMappingResponse(m *domain.EvidenceMapping) EvidenceMappingResponse {
	return EvidenceMappingResponse{
		ID:         m.ID.String(),
		TenantID:   m.TenantID.String(),
		EvidenceID: m.EvidenceID.String(),
		ControlID:  m.ControlID.String(),
		Notes:      m.Notes,
		MappedBy:   m.MappedBy.String(),
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}

type ControlCoverageResponse struct {
	Control     ControlResponse `json:"control"`
	Status      string          `json:"status"`
	EvidenceIDs []string        `json:"evidence_ids,omitempty"`
}

type FrameworkCoverageResponse struct {
	Framework FrameworkResponse   `json:"framework"`
	Controls  []ControlCoverageResponse `json:"controls"`
	Total     int                 `json:"total"`
	Mapped    int                 `json:"mapped"`
}

func ToFrameworkCoverageResponse(fc *domain.FrameworkCoverage) FrameworkCoverageResponse {
	resp := FrameworkCoverageResponse{
		Framework: ToFrameworkResponse(&fc.Framework),
		Controls:  make([]ControlCoverageResponse, len(fc.Controls)),
		Total:     fc.Total,
		Mapped:    fc.Mapped,
	}
	for i, cc := range fc.Controls {
		evIDs := make([]string, len(cc.EvidenceIDs))
		for j, eid := range cc.EvidenceIDs {
			evIDs[j] = eid.String()
		}
		resp.Controls[i] = ControlCoverageResponse{
			Control:     ToControlResponse(&cc.Control),
			Status:      string(cc.Status),
			EvidenceIDs: evIDs,
		}
	}
	return resp
}

type ErrorResponse struct {
	Error string `json:"error"`
}
