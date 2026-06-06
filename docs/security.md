# Security

## Authentication

### JWT Sessions

Identity service issues JWTs for session-based authentication:

| Field | Description |
|-------|-------------|
| Algorithm | HS256 (HMAC-SHA256) |
| Secret | Configurable via `jwt.secret` in identity config |
| Claims | `sub` (user ID), `org` (organization ID), `role`, `exp`, `iat`, `iss` |
| TTL | Configurable per session (default: 24h) |

### API Keys

For programmatic access by other services:

- Generated server-side, returned once on creation
- Stored as bcrypt hash in `api_keys.key_hash`
- Prefix stored for identification (`api_keys.key_prefix`)
- Revocable via `DELETE /api/v1/api-keys/{id}`

### Auth Middleware

Supports two auth schemes:

1. **Bearer** — Validates JWT session token
2. **ApiKey** — Looks up `X-API-Key` header, hashes and compares against stored hash

Applied via middleware wrapping chi routes:

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth(svc))
    r.Get("/api/v1/users", h.ListUsers)
})
```

## OIDC/OAuth2 Login

Configurable providers in identity config:

```yaml
oidc:
  providers:
    - name: "google"
      issuer_url: "https://accounts.google.com"
      client_id: "..."
      client_secret: "..."
      redirect_url: "http://localhost:8081/api/v1/auth/oidc/google/callback"
      scopes:
        - "openid"
        - "email"
        - "profile"
```

### Flow

1. `GET /auth/oidc/providers` — Returns available provider names
2. `GET /auth/oidc/{provider}/login` — Generates state+nonce (stored in `oidc_states` table), redirects to provider
3. `GET /auth/oidc/{provider}/callback` — Exchanges code for token, fetches userinfo, auto-provisions user

### Auto-Provisioning

- First-time OIDC login creates a new user in the first organization with `reviewer` role
- Existing accounts are linked by email match (`linked_accounts` table)

### CSRF Protection

- State parameter with stored nonce in `oidc_states` table
- Auto-expired entries cleaned on verification

## Authorization (RBAC)

| Role | Permissions |
|------|-------------|
| `admin` | Full CRUD on organizations, users, API keys, evidence, questionnaires |
| `reviewer` | Read evidence, approve/reject evidence, answer questionnaires |

Fine-grained permission checks via `HasPermission` on the API key service layer.

## Password Storage

- BCrypt with cost factor 12
- Minimum password length: 8 characters
- Configurable via `PasswordMinLength` in service config

## Network Security

- Internal gRPC (port +1000) not exposed externally in production
- NATS can be configured with TLS for inter-service messaging
- File storage uses MinIO/S3 with configurable SSL (`storage.use_ssl`)
- CORS not currently configured (frontend not built yet)

## Secrets Management

- Dev secrets in `*-dev.yaml` files (not for production)
- All secrets overridable via `EVIDRA_*` environment variables
- CI/CD uses `secrets.GITHUB_TOKEN` for GHCR authentication

## Audit Trail

The audit service records all security-relevant events:

- User login/logout
- Evidence CRUD and status changes
- Questionnaire uploads and processing
- Draft creation, approval, rejection
- All events are immutable (append-only)

## Future Security Improvements

- [ ] TLS for all gRPC communication (mTLS)
- [ ] Rate limiting on auth endpoints
- [ ] Session refresh token rotation
- [ ] OAuth2 device code flow for CLI
- [ ] Row-level security (RLS) in PostgreSQL
- [ ] Audit event publishing from all services (currently only the audit service stores events)
- [ ] OpenID Connect discovery (currently uses direct HTTP calls instead of full discovery library)
