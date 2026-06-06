-- +goose Up
CREATE TABLE approvals (
    id          UUID        PRIMARY KEY,
    evidence_id UUID        NOT NULL REFERENCES evidence_items(id) ON DELETE CASCADE,
    reviewer_id UUID        NOT NULL,
    status      VARCHAR(30) NOT NULL,
    comment     TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approvals_evidence ON approvals(evidence_id);

-- +goose Down
DROP TABLE IF EXISTS approvals;
