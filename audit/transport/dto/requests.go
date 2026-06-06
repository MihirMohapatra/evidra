package dto

type RecordAuditRequest struct {
	TenantID string         `json:"tenant_id" validate:"required,uuid"`
	ActorID  string         `json:"actor_id" validate:"required,uuid"`
	Action   string         `json:"action" validate:"required"`
	TargetID string         `json:"target_id"`
	Metadata map[string]any `json:"metadata"`
}

type ListAuditRequest struct {
	TenantID string `json:"tenant_id"`
	ActorID  string `json:"actor_id"`
	Action   string `json:"action"`
	TargetID string `json:"target_id"`
	Since    string `json:"since"`
	Until    string `json:"until"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
