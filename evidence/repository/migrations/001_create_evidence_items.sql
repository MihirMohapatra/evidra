-- +goose Up
CREATE TABLE evidence_items (
    id         UUID         PRIMARY KEY,
    tenant_id  UUID         NOT NULL,
    title      VARCHAR(500) NOT NULL,
    content    TEXT         NOT NULL DEFAULT '',
    category   VARCHAR(30)  NOT NULL,
    status     VARCHAR(30)  NOT NULL DEFAULT 'draft',
    owner_id   UUID         NOT NULL,
    source_url TEXT         NOT NULL DEFAULT '',
    tags       TEXT[]       NOT NULL DEFAULT '{}',
    version    INT          NOT NULL DEFAULT 1,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_evidence_tenant ON evidence_items(tenant_id);
CREATE INDEX idx_evidence_category ON evidence_items(category);
CREATE INDEX idx_evidence_status ON evidence_items(status);
CREATE INDEX idx_evidence_owner ON evidence_items(owner_id);
CREATE INDEX idx_evidence_expires ON evidence_items(expires_at);

-- +goose Down
DROP TABLE IF EXISTS evidence_items;
