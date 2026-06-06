package dto

type CreateEvidenceRequest struct {
	Title     string   `json:"title" validate:"required,max=500"`
	Content   string   `json:"content"`
	Category  string   `json:"category" validate:"required,oneof=policy answer claim certification architecture"`
	OwnerID   string   `json:"owner_id" validate:"required,uuid"`
	SourceURL string   `json:"source_url"`
	Tags      []string `json:"tags"`
	ExpiresAt string   `json:"expires_at"` // RFC3339
}

type UpdateEvidenceRequest struct {
	Title     string   `json:"title" validate:"omitempty,max=500"`
	Content   string   `json:"content"`
	Category  string   `json:"category" validate:"omitempty,oneof=policy answer claim certification architecture"`
	OwnerID   string   `json:"owner_id" validate:"omitempty,uuid"`
	SourceURL string   `json:"source_url"`
	Tags      []string `json:"tags"`
	ExpiresAt string   `json:"expires_at"`
}

type ApproveRequest struct {
	Comment string `json:"comment"`
}

type RejectRequest struct {
	Comment string `json:"comment" validate:"required"`
}

type ExportRequest struct {
	Comment string `json:"comment"`
}

type ListEvidenceRequest struct {
	TenantID string `json:"tenant_id"`
	Category string `json:"category"`
	Status   string `json:"status"`
	OwnerID  string `json:"owner_id"`
	Expiring bool   `json:"expiring"`
	Expired  bool   `json:"expired"`
	Tag      string `json:"tag"`
	Search   string `json:"search"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
