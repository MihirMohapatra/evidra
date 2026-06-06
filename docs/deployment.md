# Deployment

## Local Development

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 16 (optional, use Docker)
- NATS (optional, use Docker)
- MinIO (optional, use Docker)

### Quick Start

```bash
# Start infrastructure (PostgreSQL x5, NATS, MinIO)
make docker-up

# Run migrations
make migrate-identity   IDENTITY_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_identity?sslmode=disable"
make migrate-questionnaire QUESTIONNAIRE_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_questionnaire?sslmode=disable"
make migrate-evidence   EVIDENCE_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_evidence?sslmode=disable"
make migrate-orchestrator ORCHESTRATOR_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_orchestrator?sslmode=disable"
make migrate-audit      AUDIT_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_audit?sslmode=disable"
make migrate-export     EXPORT_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_export?sslmode=disable"
make migrate-compliance COMPLIANCE_DB_URL="postgres://evidra:evidra@localhost:5432/evidra_compliance?sslmode=disable"

# Start all services (separate terminals)
make run-identity
make run-questionnaire
make run-questionnaire-worker
make run-evidence
make run-orchestrator
make run-orchestrator-worker
make run-audit
make run-export
make run-compliance
```

### Docker Compose

```bash
# Build and start everything
docker compose up --build -d

# View logs
docker compose logs -f

# Stop
make docker-down
```

The compose file runs: 7 PostgreSQL databases, NATS JetStream, MinIO + bucket creation, all 9 service images.

## Docker

Multi-stage Dockerfile builds 7 targets:

```bash
# Build specific service
docker build --target identity -t evidra/identity .
docker build --target questionnaire -t evidra/questionnaire .
docker build --target questionnaire-worker -t evidra/questionnaire-worker .
docker build --target evidence -t evidra/evidence .
docker build --target orchestrator -t evidra/orchestrator .
docker build --target orchestrator-worker -t evidra/orchestrator-worker .
docker build --target audit -t evidra/audit .
docker build --target export -t evidra/export .
docker build --target compliance -t evidra/compliance .
```

Base image: `golang:1.25-alpine` → `alpine:3.20` (final stage).

## CI/CD Pipeline

### CI (`.github/workflows/ci.yml`)

Triggers on push to `master`, `feat/**`, `fix/**` and pull requests to `master`:

| Step | Tool | Description |
|------|------|-------------|
| Lint | golangci-lint v2 | Static analysis |
| Build | `go build ./...` | Compile all packages + binaries |
| Test | `go test -race -count=1 ./...` | Unit tests (30+, testify) |
| Docker | Docker Buildx | Build all 9 images (no push) |

### CD (`.github/workflows/cd.yml`)

Triggers on push to `master` or `v*` tags:

| Step | Description |
|------|-------------|
| Build & Push 9 images | Multi-architecture, pushed to `ghcr.io/<owner>/evidra/<service>:latest` and `:<sha>` |

Authentication: `secrets.GITHUB_TOKEN` for GHCR access.

## Infrastructure (Future)

### Kubernetes (planned)

Expected structure:

```
deployments/kubernetes/
├── namespace.yaml
├── configmaps/
├── secrets/
├── services/
│   ├── identity.yaml
│   ├── questionnaire.yaml
│   ├── evidence.yaml
│   ├── orchestrator.yaml
│   ├── audit.yaml
│   ├── export.yaml
│   └── compliance.yaml
├── deployments/
│   ├── identity.yaml
│   ├── questionnaire.yaml
│   ├── questionnaire-worker.yaml
│   ├── evidence.yaml
│   ├── orchestrator.yaml
│   ├── orchestrator-worker.yaml
│   ├── audit.yaml
│   ├── export.yaml
│   └── compliance.yaml
├── ingress.yaml
└── kustomization.yaml
```

### Terraform (planned)

Expected structure:

```
deployments/terraform/
├── main.tf
├── variables.tf
├── outputs.tf
├── modules/
│   ├── postgres/
│   ├── nats/
│   └── storage/
└── environments/
    ├── dev/
    ├── staging/
    └── prod/
```

### Monitoring (planned)

- **OpenTelemetry**: Trace propagation across services
- **Prometheus**: `/metrics` endpoints per service
- **Grafana**: Dashboards for latency, errors, throughput

## Configuration

Each service reads a YAML config file (default: `<service>/<service>-dev.yaml`).
All values overridable via `EVIDRA_*` environment variables (Viper binding).

Example environment variables:

```bash
export EVIDRA_DATABASE_URL="postgres://user:pass@host:5432/db"
export EVIDRA_NATS_URL="nats://nats:4222"
export EVIDRA_JWT_SECRET="production-secret"
export EVIDRA_LLM_OPENAI_KEY="sk-..."
export EVIDRA_EMBEDDER_OPENAI_KEY="sk-..."
```
