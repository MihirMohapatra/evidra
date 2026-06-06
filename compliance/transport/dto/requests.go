package dto

type CreateFrameworkRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type CreateControlRequest struct {
	ControlID       string `json:"control_id" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	Category        string `json:"category" validate:"required"`
	RiskDescription string `json:"risk_description"`
}

type MapEvidenceRequest struct {
	EvidenceID string `json:"evidence_id" validate:"required"`
	ControlID  string `json:"control_id" validate:"required"`
	Notes      string `json:"notes"`
}
