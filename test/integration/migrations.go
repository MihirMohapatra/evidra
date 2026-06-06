package integration

var identityMigrations = map[string]string{
	"001_create_organizations": `CREATE TABLE organizations (
		id         UUID PRIMARY KEY,
		name       VARCHAR(255) NOT NULL,
		slug       VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);
	CREATE INDEX idx_organizations_slug ON organizations(slug);`,

	"002_create_users": `CREATE TABLE users (
		id              UUID PRIMARY KEY,
		organization_id UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		email           VARCHAR(255) NOT NULL,
		password_hash   TEXT         NOT NULL,
		role            VARCHAR(20)  NOT NULL DEFAULT 'reviewer',
		is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
		created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		UNIQUE(organization_id, email)
	);
	CREATE INDEX idx_users_email ON users(email);
	CREATE INDEX idx_users_organization ON users(organization_id);`,

	"003_create_sessions": `CREATE TABLE sessions (
		id            UUID        PRIMARY KEY,
		user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token         TEXT        NOT NULL UNIQUE,
		refresh_token TEXT        NOT NULL UNIQUE,
		expires_at    TIMESTAMPTZ NOT NULL,
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX idx_sessions_token ON sessions(token);
	CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
	CREATE INDEX idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);`,

	"004_create_api_keys": `CREATE TABLE api_keys (
		id              UUID        PRIMARY KEY,
		organization_id UUID        NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		name            VARCHAR(255) NOT NULL,
		key_hash        TEXT        NOT NULL UNIQUE,
		key_prefix      VARCHAR(20) NOT NULL,
		is_active       BOOLEAN     NOT NULL DEFAULT TRUE,
		created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX idx_api_keys_organization ON api_keys(organization_id);
	CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);`,

	"005_create_oidc_states": `CREATE TABLE oidc_states (
		id         UUID        PRIMARY KEY,
		provider   VARCHAR(50) NOT NULL,
		state      TEXT        NOT NULL,
		nonce      TEXT        NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX idx_oidc_states_state ON oidc_states(state);
	CREATE INDEX idx_oidc_states_expires_at ON oidc_states(expires_at);`,

	"006_create_linked_accounts": `CREATE TABLE linked_accounts (
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
	CREATE INDEX idx_linked_accounts_provider_subject ON linked_accounts(provider, subject);`,
}

var evidenceMigrations = map[string]string{
	"001_create_evidence_items": `CREATE TABLE evidence_items (
		id         UUID         PRIMARY KEY,
		tenant_id  UUID         NOT NULL,
		title      VARCHAR(500) NOT NULL,
		content    TEXT         NOT NULL DEFAULT '',
		category   VARCHAR(30)  NOT NULL,
		status     VARCHAR(30)  NOT NULL DEFAULT 'DRAFT',
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
	CREATE INDEX idx_evidence_expires ON evidence_items(expires_at);`,

	"002_create_approvals": `CREATE TABLE approvals (
		id          UUID        PRIMARY KEY,
		evidence_id UUID        NOT NULL REFERENCES evidence_items(id) ON DELETE CASCADE,
		reviewer_id UUID        NOT NULL,
		status      VARCHAR(30) NOT NULL,
		comment     TEXT        NOT NULL DEFAULT '',
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX idx_approvals_evidence ON approvals(evidence_id);`,
}
