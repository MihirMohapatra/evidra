package dto

type ExportRequest struct {
	EvidenceID string `json:"evidence_id" validate:"required"`
	Format     string `json:"format" validate:"required"`
}
