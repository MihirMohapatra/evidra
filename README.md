# Evidra

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)
[![Build](https://img.shields.io/badge/Build-passing-green)](https://github.com/MihirMohapatra/evidra/actions)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)](https://github.com/MihirMohapatra/evidra/pulls)

**Evidra** is a modular, microservice-based platform for managing evidence repositories, compliance frameworks, questionnaires, and AI-assisted audit workflows. Built with Go, it follows domain-driven design with clean architecture boundaries.

## Architecture

```
                    ┌──────────────┐
                    │  Frontend    │
                    │  (Next.js)   │
                    └──────┬───────┘
                           │ HTTP / gRPC
           ┌───────────────┼─────────────────────────────────┐
           │               │                                 │
    ┌──────▼──────┐ ┌─────▼──────┐  ┌─────────▼──────────┐ ┌─▼──────────┐
    │   Identity  │ │Questionnaire│ │  Evidence          │ │  Export    │
    │   Service   │ │  Service    │ │  Repository        │ │  Service   │
    └──────┬──────┘ └─────┬──────┘  └─────────┬──────────┘ └─────┬──────┘
           │               │                   │                  │
           │         ┌─────▼──────┐            │                  │
           │         │   Worker   │            │                  │
           │         │ (NATS sub) │            │                  │
           │         └─────┬──────┘            │                  │
           │               │                   │                  │
    ┌──────▼──────┐ ┌─────▼──────┐  ┌─────────▼──────────┐ ┌─────▼──────┐
    │    Audit    │ │Orchestrator│  │  Worker            │ │ Compliance │
    │   Service   │ │  Service   │  │  (NATS sub)        │ │  Mapper    │
    └─────────────┘ └────────────┘  └─────────────────────┘ └────────────┘
           │               │                   │                  │
           └───────────────┼───────────────────┼──────────────────┘
                           │                   │
                    ┌──────▼───────────────────▼──────┐
                    │           NATS                   │
                    │        Message Bus               │
                    └──────────────────────────────────┘
```

### Services

| Service | Port | Status | Description |
|---------|------|--------|-------------|
| **Identity** | 8081 | ✅ Live | Organizations, users, roles, JWT auth, API keys, OIDC |
| **Questionnaire** | 8082 | ✅ Live | Upload, parse, extract questions from PDF/XLSX/DOCX |
| **Evidence** | 8083 | ✅ Live | Evidence repository with embeddings & approval workflow |
| **Orchestrator** | 8084 | ✅ Live | RAG, LLM integration (OpenAI/Claude/Local), draft generation |
| **Audit** | 8085 | ✅ Live | Event sourcing & audit trail with NATS ingestion |
| **Export** | 8086 | ✅ Live | PDF/XLSX/DOCX generation with MinIO storage |
| **Compliance** | 8087 | ✅ Live | Framework mapper (SOC2, ISO27001, NIST, PCI-DSS, HIPAA, FedRAMP) |
| **Frontend** | 3000 | ✅ Live | Next.js 15 dashboard with TypeScript |

### Supporting Infrastructure

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Database** | PostgreSQL 16 + pgvector | State persistence + vector embeddings |
| **Messaging** | NATS (JetStream) | Async event bus between services |
| **Storage** | MinIO (S3-compatible) | File/attachment storage |
| **Observability** | OpenTelemetry + Prometheus + Grafana | Distributed tracing, metrics, dashboards |
| **Container** | Docker / Docker Compose | Local development orchestration |
| **Orchestration** | Kubernetes + Kustomize | Production deployment |
| **IaC** | Terraform (AWS) | Cloud infrastructure provisioning |

## Tech Stack

| Category | Choice |
|----------|--------|
| **Language** | Go 1.25+ |
| **HTTP** | chi/v5 |
| **Database** | PostgreSQL 16 + pgx/v5 + pgvector |
| **Migrations** | goose |
| **Validation** | go-playground/validator |
| **Auth** | JWT (golang-jwt/v5) + bcrypt + OIDC/OAuth2 |
| **Messaging** | NATS |
| **Storage** | MinIO (S3-compatible) via minio-go |
| **Config** | Viper + YAML + env vars |
| **Logging** | slog |
| **Tracing** | OpenTelemetry (OTLP gRPC) |
| **Metrics** | Prometheus client_golang |
| **gRPC** | google.golang.org/grpc + protobuf |
| **Frontend** | Next.js 15 + TypeScript + Tailwind CSS |
| **Testing** | testify + testcontainers-go |
| **Docs** | OpenAPI 3.0, Protocol Buffers |
| **CI/CD** | GitHub Actions |

## Project Structure

```
evidra/
├── frontend/                 # Next.js 15 TypeScript frontend
│   └── src/
│       ├── app/              # App Router pages
│       ├── components/       # Reusable UI components
│       ├── contexts/         # React contexts (auth)
│       └── lib/              # API client & types
│
├── identity/                 # Identity & access management
│   ├── cmd/server/           # HTTP + gRPC server
│   ├── domain/               # Business entities & rules
│   ├── repository/           # Data access (postgres)
│   ├── service/              # Use cases & business logic
│   └── transport/            # HTTP handlers, middleware, gRPC
│
├── questionnaire/            # Questionnaire management
│   ├── cmd/server/           # API server
│   ├── cmd/worker/           # Document processor (NATS sub)
│   ├── domain/
│   ├── events/               # NATS event definitions
│   ├── parser/               # PDF/XLSX/DOCX extraction
│   ├── repository/
│   ├── service/
│   └── transport/
│
├── evidence/                 # Evidence repository
│   ├── cmd/server/
│   ├── domain/
│   ├── events/
│   ├── repository/
│   ├── service/
│   └── transport/
│
├── orchestrator/             # AI orchestrator (RAG + LLM)
│   ├── cmd/server/
│   ├── cmd/worker/
│   ├── domain/
│   ├── events/
│   ├── repository/
│   ├── service/
│   └── transport/
│
├── audit/                    # Audit trail service
│   ├── cmd/server/
│   ├── domain/
│   ├── events/
│   ├── repository/
│   ├── service/
│   └── transport/
│
├── export/                   # Document export service
│   ├── cmd/server/
│   ├── domain/
│   ├── events/
│   ├── repository/
│   ├── service/              # PDF/XLSX/DOCX generators
│   └── transport/
│
├── compliance/               # Compliance framework mapper
│   ├── cmd/server/
│   ├── domain/
│   ├── events/
│   ├── repository/
│   ├── service/
│   └── transport/
│
├── pkg/                      # Shared libraries
│   ├── queue/                # NATS message bus abstraction
│   ├── storage/              # MinIO/S3 abstraction
│   └── telemetry/            # OpenTelemetry + Prometheus
│
├── api/                      # API specifications
│   ├── openapi/              # OpenAPI 3.0 specs per service
│   ├── proto/                # Protocol Buffer definitions
│   └── gen/                  # Generated code
│
├── deployments/              # Production deployments
│   ├── kubernetes/           # K8s manifests (16 resources)
│   ├── terraform/            # AWS IaC with modules
│   └── grafana/              # Dashboard + Prometheus config
│
├── test/                     # Integration tests
│   └── integration/          # Testcontainers-based
│
├── configs/                  # Environment YAML configs
├── migrations/               # Database migrations
├── scripts/                  # Utility scripts
├── docs/                     # Architecture & design docs
├── docker-compose.yml        # Local dev orchestration
├── Dockerfile                # Multi-stage build
└── Makefile                  # Build & dev commands
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Node.js 20+ (for frontend)

### Local Development

```bash
# Clone
git clone https://github.com/MihirMohapatra/evidra.git
cd evidra

# Start dependencies (PostgreSQL ×5, NATS, MinIO)
docker compose up -d

# Run migrations for all services
go run ./migrations

# Start backend services (in separate terminals)
go run ./identity/cmd/server
go run ./questionnaire/cmd/server
go run ./questionnaire/cmd/worker
go run ./evidence/cmd/server
go run ./orchestrator/cmd/server
go run ./audit/cmd/server
go run ./export/cmd/server
go run ./compliance/cmd/server

# Start frontend (in another terminal)
cd frontend
npm install
npm run dev
```

### Configuration

Each service reads a YAML config file. Environment variables with the `EVIDRA_` prefix override config values .

```yaml
# configs/dev.yaml
server:
  host: "0.0.0.0"
  port: 8081

database:
  url: "postgres://evidra:evidra@localhost:5432/evidra_identity?sslmode=disable"

nats:
  url: "nats://localhost:4222"
```

## Development

### Build

```bash
go build ./...                    # Build all Go packages
cd frontend && npm run build      # Build frontend
```

### Test

```bash
go test ./...                     # Unit tests
go test ./test/integration/...    # Integration tests (requires Docker)
cd frontend && npm run lint       # Frontend lint
```

### Lint

```bash
golangci-lint run ./...
```

### Database Migrations

```bash
# Run all migrations
go run ./migrations

# Or run per-service
go run ./identity/migrations
```

## API Reference

### Identity Service (port 8081)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/login` | Authenticate user |
| POST | `/api/v1/auth/refresh` | Refresh session token |
| POST | `/api/v1/auth/logout` | Invalidate session |
| GET | `/api/v1/auth/oidc/providers` | List OIDC providers |
| GET/POST | `/api/v1/organizations` | List/create organizations |
| GET/PUT/DELETE | `/api/v1/organizations/{id}` | Organization CRUD |
| GET/POST | `/api/v1/users` | List/create users |
| GET/PUT/DELETE | `/api/v1/users/{id}` | User CRUD |
| GET/POST | `/api/v1/api-keys` | List/create API keys |
| DELETE | `/api/v1/api-keys/{id}` | Revoke API key |

### Evidence Service (port 8083)

| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/api/v1/evidence` | List/create evidence |
| GET/PUT/DELETE | `/api/v1/evidence/{id}` | Evidence CRUD |
| POST | `/api/v1/evidence/{id}/submit` | Submit for review |
| POST | `/api/v1/evidence/{id}/approve` | Approve evidence |
| POST | `/api/v1/evidence/{id}/reject` | Reject evidence |
| POST | `/api/v1/evidence/{id}/export` | Mark as exported |
| GET | `/api/v1/evidence/{id}/approvals` | Approval history |

### Orchestrator Service (port 8084)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/orchestrator/answer` | Generate answer (RAG + LLM) |
| GET | `/api/v1/orchestrator/drafts` | List drafts |
| GET | `/api/v1/orchestrator/drafts/{id}` | Get draft |
| POST | `/api/v1/orchestrator/drafts/{id}/approve` | Approve draft |
| POST | `/api/v1/orchestrator/drafts/{id}/reject` | Reject draft |

### Questionnaire Service (port 8082)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/questionnaires/upload` | Upload document |
| GET | `/api/v1/questionnaires` | List questionnaires |
| GET/DELETE | `/api/v1/questionnaires/{id}` | Questionnaire detail/delete |
| GET | `/api/v1/questionnaires/{id}/questions` | Extracted questions |

### Audit Service (port 8085)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/audit/events` | Record audit event |
| GET | `/api/v1/audit/events` | List audit events |

### Export Service (port 8086)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/exports` | Export evidence (PDF/XLSX/DOCX) |
| GET | `/api/v1/exports/{id}` | Get export status/details |
| GET | `/api/v1/exports` | List exports (filter by `evidence_id`) |

### Compliance Service (port 8087)

| Method | Path | Description |
|--------|------|-------------|
| POST/GET | `/api/v1/compliance/frameworks` | Create/list frameworks |
| GET/DELETE | `/api/v1/compliance/frameworks/{frameworkId}` | Framework detail/delete |
| POST/GET | `/api/v1/compliance/frameworks/{frameworkId}/controls` | Create/list controls |
| POST | `/api/v1/compliance/mappings` | Map evidence to control |
| DELETE | `/api/v1/compliance/mappings/{id}` | Remove mapping |
| GET | `/api/v1/compliance/mappings/by-control/{controlId}` | Mappings by control |
| GET | `/api/v1/compliance/frameworks/{frameworkId}/coverage` | Coverage report |

## Observability

- **Metrics**: Prometheus `/metrics` endpoint on every service
- **Tracing**: OpenTelemetry OTLP gRPC exporter (configurable endpoint)
- **Dashboards**: Grafana dashboard at `deployments/grafana/dashboard.json`
- **Scraping**: Prometheus config at `deployments/grafana/prometheus.yml`

## Deployment

### Docker

```bash
docker compose up -d --build
```

### Kubernetes

```bash
kubectl apply -k deployments/kubernetes/
```

### Terraform

```bash
cd deployments/terraform
terraform init
terraform workspace select dev
terraform apply
```

## Roadmap

- [x] Identity service (auth, orgs, users, API keys, OIDC)
- [x] Questionnaire service (upload, parse, question extraction)
- [x] Evidence repository service
- [x] AI orchestrator with RAG (OpenAI/Claude/Local)
- [x] Approval workflow engine
- [x] Audit service with NATS event ingestion
- [x] OpenAPI specs for all services
- [x] Frontend dashboard (Next.js + TypeScript)
- [x] Kubernetes deployment manifests
- [x] Terraform infrastructure-as-code
- [x] CI/CD pipeline
- [x] OpenTelemetry tracing + Prometheus metrics
- [x] Export service (PDF/XLSX/DOCX)
- [x] Compliance framework mapper
- [ ] WebSocket real-time updates
- [ ] Mobile app (React Native)

## Documentation

- [Architecture](docs/architecture.md) — Service map, communication patterns, DDD structure
- [Database Design](docs/database-design.md) — ER diagrams, indexes, migrations, pgvector
- [Scaling](docs/scaling.md) — Horizontal scalability, throughput estimates, caching
- [Security](docs/security.md) — JWT, API keys, OIDC/OAuth2, RBAC, audit trail
- [Deployment](docs/deployment.md) — Local dev, Docker, CI/CD, Kubernetes, Terraform

## License

MIT
