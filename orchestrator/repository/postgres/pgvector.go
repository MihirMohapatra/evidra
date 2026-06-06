package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/repository"
)

type EmbeddingRepo struct {
	pool *pgxpool.Pool
}

func NewEmbeddingRepo(pool *pgxpool.Pool) *EmbeddingRepo {
	return &EmbeddingRepo{pool: pool}
}

var _ repository.EmbeddingRepository = (*EmbeddingRepo)(nil)

func (r *EmbeddingRepo) SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]domain.Evidence, error) {
	vec := formatVector(embedding)
	query := fmt.Sprintf(`
		SELECT id, evidence_id, title, content, category, source_url, 1 - (embedding <=> %s) AS score
		FROM evidence_embeddings
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> %s
		LIMIT $1`, vec, vec)

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}
	defer rows.Close()

	var results []domain.Evidence
	for rows.Next() {
		var e domain.Evidence
		var id, evID uuid.UUID
		if err := rows.Scan(&id, &evID, &e.Title, &e.Content, &e.Category, &e.SourceURL, &e.Score); err != nil {
			return nil, fmt.Errorf("scan evidence: %w", err)
		}
		e.ID = evID
		results = append(results, e)
	}
	if results == nil {
		results = []domain.Evidence{}
	}
	return results, nil
}

func (r *EmbeddingRepo) Upsert(ctx context.Context, chunk *domain.EvidenceChunk) error {
	vec := formatVector(chunk.Embedding)
	query := fmt.Sprintf(`
		INSERT INTO evidence_embeddings (id, evidence_id, title, content, category, source_url, metadata, embedding, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, %s, $8)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata`, vec)

	_, err := r.pool.Exec(ctx, query,
		chunk.ID, chunk.EvidenceID, "", chunk.Content, "", "", chunk.Metadata, chunk.CreatedAt)
	return err
}

func formatVector(v []float32) string {
	parts := make([]string, len(v))
	for i, val := range v {
		parts[i] = fmt.Sprintf("%f", val)
	}
	return fmt.Sprintf("'[%s]'", strings.Join(parts, ","))
}
