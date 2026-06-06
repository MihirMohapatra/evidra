-- +goose Up
CREATE TABLE linked_accounts (
    id         UUID        PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider   VARCHAR(50) NOT NULL,
    subject    TEXT        NOT NULL,
    email      VARCHAR(255) NOT NULL DEFAULT '',
    name       VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, subject)
);

CREATE INDEX idx_linked_accounts_user ON linked_accounts(user_id);
CREATE INDEX idx_linked_accounts_provider_subject ON linked_accounts(provider, subject);

-- +goose Down
DROP TABLE IF EXISTS linked_accounts;
