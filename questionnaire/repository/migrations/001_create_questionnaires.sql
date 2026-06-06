-- +goose Up
CREATE TABLE questionnaires (
    id         UUID         PRIMARY KEY,
    tenant_id  UUID         NOT NULL,
    title      VARCHAR(255) NOT NULL,
    file_name  VARCHAR(255) NOT NULL,
    file_url   TEXT         NOT NULL,
    file_type  VARCHAR(20)  NOT NULL,
    file_size  BIGINT       NOT NULL DEFAULT 0,
    status     VARCHAR(20)  NOT NULL DEFAULT 'uploaded',
    version    INT          NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questionnaires_tenant ON questionnaires(tenant_id);
CREATE INDEX idx_questionnaires_status ON questionnaires(status);

-- +goose Down
DROP TABLE IF EXISTS questionnaires;
