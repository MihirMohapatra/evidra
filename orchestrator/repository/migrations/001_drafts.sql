-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE evidence_embeddings (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    evidence_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    source_url TEXT NOT NULL DEFAULT '',
    metadata JSONB DEFAULT '{}',
    embedding vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_embeddings_tenant ON evidence_embeddings(tenant_id);
CREATE INDEX idx_embeddings_vector ON evidence_embeddings
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE TABLE drafts (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL,
    question_text TEXT NOT NULL,
    answer TEXT NOT NULL,
    confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
    model_used TEXT NOT NULL DEFAULT '',
    evidence_ids UUID[] NOT NULL DEFAULT '{}',
    reasoning TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    feedback TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_drafts_question_id ON drafts(question_id);
CREATE INDEX idx_drafts_status ON drafts(status);

-- +goose Down
DROP TABLE IF EXISTS drafts;
DROP TABLE IF EXISTS evidence_embeddings;
