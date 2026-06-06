package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrLLMError          = errors.New("llm service error")
	ErrValidationFailed  = errors.New("response validation failed")
	ErrEmbeddingFailed   = errors.New("embedding generation failed")
	ErrEvidenceNotFound  = errors.New("no relevant evidence found")
)
