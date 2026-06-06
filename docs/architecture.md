# Architecture

## Overview

Evidra is a modular microservice platform for managing evidence repositories, compliance questionnaires, and AI-assisted audit workflows. Each service is independently deployable, owns its own data store, and communicates via NATS messaging and synchronous gRPC/REST APIs.

## Service Map

```
                            ┌──────────────────────────────────────────────────────────────────┐
                            │                        External Clients                           │
                            │        (CLI, CI/CD pipelines, Web UI via Next.js)                │
                            └────────┬──────────────┬────────────────┬──────────────────────────┘
                                     │              │                │
                            ┌────────▼──────┐ ┌─────▼──────┐  ┌────▼──────────────┐
                            │  HTTP REST    │ │   gRPC     │  │  NATS JetStream   │
                            │  (chi/v5)     │ │ (internal) │  │  (async events)   │
                            │  Port :80xx   │ │ Port +1000 │  │  Port 4222        │
                            └────────┬──────┘ └─────┬──────┘  └─────────┬─────────┘
                                     │              │                  │
          ┌──────────────────────────┼──────────────┼──────────────────┼──────────────────────────┐
          │                          │              │                  │                          │
          │              ┌───────────▼──────────────▼──────────────────▼───────────────┐          │
          │              │                     Identity :8081                          │          │
          │              │  JWT auth, sessions, API keys, organizations, users, roles   │          │
          │              │  DB: evidra_identity (pg port 5433)                          │          │
          │              └───────────┬──────────────┬────────────────────────────────────┘          │
          │                          │              │                                             │
          │              ┌───────────▼──────────────▼────┐  ┌──────────────────────────────┐      │
          │              │   Questionnaire :8082         │  │  Questionnaire Worker        │      │
          │              │  Upload, parse PDF/XLSX/DOCX  │  │  (NATS sub) extracts text,   │      │
          │              │  DB: evidra_questionnaire     │  │  detects questions           │      │
          │              │  (pg port 5434)               │  │                              │      │
          │              └───────────┬────────────────────┘  └──────────────┬───────────────┘      │
          │                          │                                      │                      │
          │              ┌───────────▼──────────────────────────────────────▼──────┐              │
          │              │              Evidence :8083                              │              │
          │              │  CRUD, approval state machine (pkg/workflow)             │              │
          │              │  5 categories, 4 statuses, NATS events                  │              │
          │              │  DB: evidra_evidence (pg port 5435)                     │              │
          │              └───────────┬──────────────────────────────────────────────┘              │
          │                          │                                                             │
          │              ┌───────────▼──────────────────────────────────────────────┐              │
          │              │           Orchestrator :8084                              │              │
          │              │  RAG pipeline: question→embedding→pgvector search→       │              │
          │              │  LLM (OpenAI/Claude/Local) → draft                       │              │
          │              │  DB: evidra_orchestrator (pg port 5436, pgvector ext)    │              │
          │              └───────────┬──────────────────────────────────────────────┘              │
          │                          │                                                             │
          │              ┌───────────▼──────────────────────────────┐                              │
          │              │       Orchestrator Worker                │                              │
          │              │  (NATS sub) listens for evidence.created │                              │
          │              │  generates embeddings, upserts pgvector  │                              │
          │              └───────────┬──────────────────────────────┘                              │
          │                          │                                                             │
          │              ┌───────────▼──────────────────────────────┐                              │
          │              │          Audit :8085                      │                              │
          │              │  Event-sourced append-only audit log     │                              │
          │              │  Predefined actions, JSONB metadata      │                              │
          │              │  DB: evidra_audit (pg port 5437)         │                              │
          │              └───────────┬──────────────────────────────┘                              │
          │                          │                                                             │
          │              ┌───────────▼──────────────────────────────┐                              │
          │              │          Export :8086                     │                              │
          │              │  PDF (gofpdf), XLSX (excelize),          │                              │
          │              │  DOCX (custom OOXML) generation          │                              │
          │              │  uploads to MinIO, NATS events           │                              │
          │              │  DB: evidra_export (pg port 5438)        │                              │
          │              └───────────┬──────────────────────────────┘                              │
          │                          │                                                             │
          │              ┌───────────▼──────────────────────────────┐                              │
          │              │      Compliance Mapper :8087              │                              │
          │              │  6 framework seeds (SOC2, ISO27001, etc) │                              │
          │              │  Control library (admin/tech/physical)   │                              │
          │              │  Evidence mapping & coverage reports     │                              │
          │              │  DB: evidra_compliance (pg port 5439)    │                              │
          │              └──────────────────────────────────────────┘                              │
          │                                                                                        │
          └────────────────────────────────────────────────────────────────────────────────────────┘
```

## Communication Patterns

| Pattern | Protocol | Use Case |
|---------|----------|----------|
| Synchronous | HTTP REST (chi/v5) | Client-facing API, CRUD operations |
| Internal RPC | gRPC (protobuf) | Inter-service queries (port = HTTP + 1000) |
| Async events | NATS pub/sub | Document processing, evidence indexing, notifications |
| File storage | S3 API (MinIO) | Document uploads, evidence artifacts |

## Event Flow

```
questionnaire.uploaded  ──►  Questionnaire Worker  ──►  questions.saved
                                                                 │
evidence.created        ──►  Orchestrator Worker   ──►  embeddings.upserted

export.requested        ──►  Export Service        ──►  export.completed / export.failed
compliance.evidence.mapped ──►  Compliance Mapper   ──►  coverage.updated
```

## Domain-Driven Design per Service

```
<service>/
├── cmd/
│   ├── server/            # Service entry point (HTTP + gRPC)
│   └── worker/            # Async worker entry point (NATS subscriber)
├── domain/                # Business entities, value objects, enums
├── events/                # NATS event type definitions
├── repository/
│   ├── interfaces.go      # Repository interface
│   └── postgres/          # pgx implementation
├── service/               # Use cases, business logic orchestration
├── transport/
│   ├── http_handler.go    # chi router, handlers, middleware
│   ├── dto/               # Request/response DTOs
│   └── grpc/              # gRPC server implementation
└── internal/config/       # Viper-based config loading
```

## Shared Libraries (`pkg/`)

| Package | Purpose |
|---------|---------|
| `pkg/storage` | MinIO/S3 abstraction (upload, download, delete, signed URLs) |
| `pkg/queue` | NATS pub/sub EventBus interface |
| `pkg/workflow` | Generic state machine engine |

## API Specifications (`api/`)

```
api/
├── proto/
│   ├── buf.gen.yaml
│   ├── generate.ps1
│   └── evidra/v1/
│       ├── common.proto
│       ├── identity.proto
│       ├── evidence.proto
│       ├── questionnaire.proto
│       ├── orchestrator.proto
│       └── audit.proto
└── gen/evidra/v1/          # Generated Go code
    ├── *.pb.go
    └── *_grpc.pb.go
```

## Service Ports

| Service | HTTP | gRPC | Database |
|---------|------|------|----------|
| Identity | 8081 | 9081 | 5433 |
| Questionnaire | 8082 | 9082 | 5434 |
| Evidence | 8083 | 9083 | 5435 |
| Orchestrator | 8084 | 9084 | 5436 |
| Audit | 8085 | 9085 | 5437 |
| Export | 8086 | 9086 | 5438 |
| Compliance | 8087 | 9087 | 5439 |
