package domain

import "github.com/google/uuid"

type Evidence struct {
	ID        uuid.UUID
	Title     string
	Content   string
	Category  string
	Score     float64
	SourceURL string
}

type EvidenceChunk struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	EvidenceID uuid.UUID
	Content    string
	Embedding  []float32
	Metadata   map[string]any
	CreatedAt  string
}
