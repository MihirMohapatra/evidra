-- +goose Up
CREATE TABLE oidc_states (
    id         UUID        PRIMARY KEY,
    provider   VARCHAR(50) NOT NULL,
    state      TEXT        NOT NULL,
    nonce      TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oidc_states_state ON oidc_states(state);
CREATE INDEX idx_oidc_states_expires_at ON oidc_states(expires_at);

-- +goose Down
DROP TABLE IF EXISTS oidc_states;
