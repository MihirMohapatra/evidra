# Evidra

**Evidra** is a modular, microservice-based platform for managing evidence repositories, compliance frameworks, questionnaires, and AI-assisted audit workflows. Built with Go, it follows a domain-driven design with clean architecture boundaries.

## Architecture

Evidra is composed of independently deployable services that communicate via NATS messaging and REST/gRPC APIs.

```
                   ┌──────────────┐
                   │   API GW     │
                   └──────┬───────┘
                          │
          ┌───────────────┼───────────────────┐
          │               │                   │
   ┌──────▼──────┐ ┌─────▼──────┐  ┌─────────▼──────────┐
   │   Identity  │ │Questionnaire│  │  Evidence          │
   │   Service   │ │  Service    │  │  Repository        │
   └──────┬──────┘ └─────┬──────┘  └─────────┬──────────┘
          │               │                   │
          │         ┌─────▼──────┐            │
          │         │   Worker   │            │
          │         │ (NATS sub) │            │
          │         └─────┬──────┘            │
          │               │                   │
          └───────────────┼───────────────────┘
                          │
                   ┌──────▼──────┐
                   │    NATS     │
                   │   Message   │
                   │    Bus      │
                   └─────────────┘
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| **Identity** | 8081 | Organizations, users, roles, auth (JWT), API keys |
| **Questionnaire** | 8082 | Upload, parse, extract questions from PDF/XLSX/DOCX |
| **Evidence** | _planned_ | Evidence repository with embeddings & approval workflow |
| **AI Orchestrator** | _planned_ | RAG, LLM integration, draft generation |
| **Compliance** | _planned_ | Framework mapping & control management |
| **Audit** | _planned_ | Event sourcing & audit trail |
| **Export** | _planned_ | PDF/XLSX/DOCX report generation |

## Tech Stack

| Category | Choice |
|----------|--------|
| **Language** | Go 1.22+ |
| **HTTP** | chi/v5 |
| **Database** | PostgreSQL 16 + pgx/v5 |
| **Migrations** | goose |
| **Validation** | go-playground/validator |
| **Auth** | JWT (golang-jwt/v5) + bcrypt |
| **Messaging** | NATS |
| **Storage** | MinIO (S3-compatible) |
| **Config** | Viper |
| **Logging** | slog / zerolog |
| **Testing** | testify |
| **Docs** | OpenAPI 3.0, Protocol Buffers |

## Project Structure

```
evidra/
├── cmd/                    # Application entry points
│   ├── api/                # API gateway server
│   ├── worker/             # Async job processor
│   └── migration/          # Database migration runner
│
├── identity/               # Identity & access management service
│   ├── cmd/server/         # Service entry point
│   ├── domain/             # Business entities & rules
│   ├── repository/         # Data access (interfaces + postgres)
│   ├── service/            # Use cases & business logic
│   └── transport/          # HTTP handlers & DTOs
│
├── questionnaire/          # Questionnaire management service
│   ├── cmd/server/         # API server entry point
│   ├── cmd/worker/         # Document processor worker
│   ├── domain/             # Business entities & status machine
│   ├── parser/             # PDF/XLSX/DOCX text extraction
│   ├── repository/         # Data access layer
│   ├── service/            # Upload & process orchestration
│   ├── events/             # NATS event definitions
│   └── transport/          # HTTP handlers & DTOs
│
├── pkg/                    # Shared libraries
│   ├── storage/            # MinIO/S3 file storage abstraction
│   └── queue/              # NATS message bus abstraction
│
├── api/                    # API specifications
│   ├── openapi.yaml
│   └── proto/              # Protocol Buffer definitions
│
├── internal/               # Internal shared packages
│   └── shared/             # Errors, middleware, logger, security
│
├── migrations/             # Root-level database migrations
├── deployments/            # Docker, Kubernetes, Terraform
├── configs/                # Environment configuration
├── scripts/                # Utility scripts
└── tests/                  # Integration & E2E tests
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 16
- NATS server
- MinIO (or S3-compatible storage)

### Local Development

```bash
# Clone the repository
git clone https://github.com/MihirMohapatra/evidra.git
cd evidra

# Start dependencies (PostgreSQL, NATS, MinIO)
docker compose up -d

# Run database migrations
go run ./cmd/migration

# Start identity service
go run ./identity/cmd/server

# Start questionnaire service
go run ./questionnaire/cmd/server

# Start document worker
go run ./questionnaire/cmd/worker
```

### Configuration

Each service reads a YAML configuration file (default: `<service>/<service>-dev.yaml`). Environment variables with the `EVIDRA_` prefix override config values.

```yaml
# identity/identity-dev.yaml
server:
  host: "0.0.0.0"
  port: 8081

database:
  url: "postgres://evidra:evidra@localhost:5432/evidra_identity?sslmode=disable"

jwt:
  secret: "dev-secret-change-in-production"
  session_ttl: 24h
```

## Development

### Build

```bash
go build ./...               # Build all packages
go build ./identity/...      # Build identity service
go build ./questionnaire/... # Build questionnaire service
```

### Test

```bash
go test ./...                # Run all tests
go test ./identity/...       # Run identity tests
```

### Lint

```bash
golangci-lint run ./...
```

### Database Migrations

Migrations use [goose](https://github.com/pressly/goose):

```bash
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir identity/repository/migrations postgres "$IDENTITY_DB_URL" up
```

## Services Detail

### Identity Service

Handles organizations, user management, authentication, and API keys.

**Endpoints:**

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/login` | Authenticate user |
| POST | `/api/v1/auth/refresh` | Refresh session token |
| POST | `/api/v1/auth/logout` | Invalidate session |
| GET/POST | `/api/v1/organizations` | List/create organizations |
| GET/PUT/DELETE | `/api/v1/organizations/{id}` | Organization CRUD |
| GET/POST | `/api/v1/users` | List/create users |
| GET/PUT/DELETE | `/api/v1/users/{id}` | User CRUD |
| GET/POST | `/api/v1/api-keys` | List/create API keys |
| DELETE | `/api/v1/api-keys/{id}` | Revoke API key |

### Questionnaire Service

Handles document upload, parsing, and question extraction.

**Endpoints:**

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/questionnaires/upload` | Upload document (multipart) |
| GET | `/api/v1/questionnaires` | List questionnaires |
| GET | `/api/v1/questionnaires/{id}` | Get questionnaire details |
| DELETE | `/api/v1/questionnaires/{id}` | Delete questionnaire |
| GET | `/api/v1/questionnaires/{id}/questions` | Get extracted questions |

## Roadmap

- [x] Identity service (auth, orgs, users, API keys)
- [x] Questionnaire service (upload, parse, question extraction)
- [ ] Evidence repository service
- [ ] AI orchestrator with RAG
- [ ] Compliance framework mapper
- [ ] Approval workflow engine
- [ ] Audit service
- [ ] Export service
- [ ] gRPC internal communication
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline

## License

MIT
