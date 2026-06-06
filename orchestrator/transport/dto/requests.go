package dto

type AnswerRequest struct {
	Question string `json:"question" validate:"required"`
	Context  string `json:"context"`
	TenantID string `json:"tenant_id"`
}

type ApproveDraftRequest struct {
	Feedback string `json:"feedback"`
}

type RejectDraftRequest struct {
	Feedback string `json:"feedback" validate:"required"`
}
