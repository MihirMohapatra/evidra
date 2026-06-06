CREATE TABLE exports (
    id          UUID         PRIMARY KEY,
    tenant_id   UUID         NOT NULL,
    evidence_id UUID         NOT NULL,
    requester_id UUID        NOT NULL,
    format      VARCHAR(10)  NOT NULL,
    file_url    TEXT         NOT NULL DEFAULT '',
    file_size   BIGINT       NOT NULL DEFAULT 0,
    status      VARCHAR(20)  NOT NULL DEFAULT 'pending',
    error       TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exports_tenant ON exports(tenant_id);
CREATE INDEX idx_exports_evidence ON exports(evidence_id);
CREATE INDEX idx_exports_status ON exports(status);
