# Database Design

Each service owns an independent PostgreSQL 16 database with isolated schemas. Migrations are managed per-service using [goose](https://github.com/pressly/goose).

## Identity Service (`evidra_identity`)

### Tables

```mermaid
erDiagram
    organizations {
        uuid id PK
        string name
        string slug UK
        text description
        jsonb metadata
        timestamptz created_at
        timestamptz updated_at
    }

    users {
        uuid id PK
        uuid organization_id FK
        string email UK
        string password_hash
        string display_name
        string role "admin | reviewer"
        timestamptz created_at
        timestamptz updated_at
    }

    sessions {
        uuid id PK
        uuid user_id FK
        string token UK
        timestamptz expires_at
        timestamptz created_at
    }

    api_keys {
        uuid id PK
        uuid organization_id FK
        string name
        string key_hash
        string key_prefix
        string[] permissions
        timestamptz expires_at
        boolean active
        timestamptz created_at
    }

    oidc_states {
        uuid id PK
        string state UK
        string nonce
        string provider
        string redirect_url
        timestamptz expires_at
        timestamptz created_at
    }

    linked_accounts {
        uuid id PK
        uuid user_id FK
        string provider
        string subject
        string email
        jsonb raw_attributes
        timestamptz created_at
    }

    organizations ||--o{ users : has
    users ||--o{ sessions : has
    organizations ||--o{ api_keys : has
    users ||--o{ linked_accounts : links
```

### Indexes
- `users(email)` UNIQUE
- `sessions(token)` UNIQUE
- `api_keys(key_hash)` UNIQUE
- `oidc_states(state)` UNIQUE
- `linked_accounts(provider, subject)` UNIQUE
- `linked_accounts(user_id)` FK

## Questionnaire Service (`evidra_questionnaire`)

### Tables

```mermaid
erDiagram
    questionnaires {
        uuid id PK
        uuid tenant_id
        string title
        string file_name
        string file_url
        string file_type
        bigint file_size
        string status "pending | parsing | parsed | failed"
        int version
        timestamptz created_at
        timestamptz updated_at
    }

    questions {
        uuid id PK
        uuid questionnaire_id FK
        string text
        string type "multiple_choice | true_false | open_ended | fill_blank | matching"
        int order
        text[] options
        timestamptz created_at
        timestamptz updated_at
    }

    questionnaires ||--o{ questions : contains
```

### Indexes
- `questions(questionnaire_id)` FK

## Evidence Service (`evidra_evidence`)

### Tables

```mermaid
erDiagram
    evidence {
        uuid id PK
        uuid tenant_id
        string title
        string description
        string category "policy | answer | claim | certification | architecture"
        string status "draft | needs_review | approved | exported"
        string file_url
        string file_type
        bigint file_size
        jsonb metadata
        uuid created_by
        timestamptz created_at
        timestamptz updated_at
    }

    approval_log {
        uuid id PK
        uuid evidence_id FK
        string from_status
        string to_status
        uuid actor_id
        text comment
        timestamptz created_at
    }

    evidence ||--o{ approval_log : tracks
```

### Indexes
- `evidence(tenant_id, status)` composite
- `evidence(created_by)` FK
- `approval_log(evidence_id)` FK

## Orchestrator Service (`evidra_orchestrator`, requires pgvector)

### Tables

```mermaid
erDiagram
    evidence_embeddings {
        uuid id PK
        uuid tenant_id
        uuid evidence_id UK
        vector(1536) embedding
        text content_snippet
        timestamptz created_at
        timestamptz updated_at
    }

    drafts {
        uuid id PK
        uuid question_id
        text question_text
        text answer
        float confidence
        string model_used
        uuid[] evidence_ids
        text reasoning
        string feedback
        string status "pending | approved | rejected"
        timestamptz created_at
        timestamptz updated_at
    }

    evidence_embeddings ||--o| evidence : represents
```

### Indexes
- `evidence_embeddings(tenant_id, evidence_id)` UNIQUE composite
- `evidence_embeddings` pgvector IVFFlat index on `embedding` for cosine similarity search

## Audit Service (`evidra_audit`)

### Tables

```mermaid
erDiagram
    audit_events {
        uuid id PK
        uuid tenant_id
        uuid actor_id
        string action "evidence_created | evidence_updated | evidence_deleted | evidence_exported | evidence_status_changed | questionnaire_uploaded | questionnaire_parsed | questionnaire_failed | draft_created | draft_approved | draft_rejected"
        string target_id
        timestamptz timestamp
        jsonb metadata
    }
```

### Indexes
- `audit_events(tenant_id, timestamp)` composite (primary query pattern)
- `audit_events(actor_id)` FK filter
- `audit_events(action)` filter
- `audit_events(target_id)` filter

## Export Service (`evidra_export`)

### Tables

```mermaid
erDiagram
    exports {
        uuid id PK
        uuid tenant_id
        uuid evidence_id
        string format "pdf | xlsx | docx"
        string status "pending | processing | completed | failed"
        string file_url
        bigint file_size
        string error_message
        timestamptz created_at
        timestamptz updated_at
    }
```

### Indexes
- `exports(tenant_id, status)` composite
- `exports(evidence_id)` FK

## Compliance Service (`evidra_compliance`)

### Tables

```mermaid
erDiagram
    frameworks {
        uuid id PK
        uuid tenant_id
        string name UK
        string slug UK
        string description
        string version
        jsonb metadata
        timestamptz created_at
        timestamptz updated_at
    }

    controls {
        uuid id PK
        uuid framework_id FK
        string control_id
        string title
        text description
        string category "administrative | technical | physical"
        jsonb metadata
        timestamptz created_at
        timestamptz updated_at
    }

    evidence_mappings {
        uuid id PK
        uuid control_id FK
        uuid evidence_id
        text notes
        timestamptz created_at
        timestamptz updated_at
    }

    question_mappings {
        uuid id PK
        uuid control_id FK
        uuid question_id
        text notes
        timestamptz created_at
    }

    frameworks ||--o{ controls : contains
    controls ||--o{ evidence_mappings : maps
    controls ||--o{ question_mappings : maps
```

### Indexes
- `frameworks(tenant_id, slug)` UNIQUE composite
- `controls(framework_id, control_id)` UNIQUE composite
- `evidence_mappings(control_id, evidence_id)` UNIQUE composite
- `question_mappings(control_id, question_id)` UNIQUE composite

## Migration Strategy

Each service has its own migration directory:

```
<service>/repository/migrations/
├── 001_create_organizations.sql
├── 002_create_users.sql
├── 003_create_sessions.sql
├── 004_create_api_keys.sql
├── 005_create_oidc_states.sql
└── 006_create_linked_accounts.sql
```

Run via Makefile:

```bash
make migrate-identity   IDENTITY_DB_URL="..."
make migrate-questionnaire QUESTIONNAIRE_DB_URL="..."
make migrate-evidence   EVIDENCE_DB_URL="..."
make migrate-orchestrator ORCHESTRATOR_DB_URL="..."
make migrate-audit      AUDIT_DB_URL="..."
make migrate-export     EXPORT_DB_URL="..."
make migrate-compliance COMPLIANCE_DB_URL="..."
```

## pgvector Setup

The orchestrator database requires the pgvector extension:

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE INDEX idx_evidence_embeddings_ivfflat
ON evidence_embeddings
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
```

Vector search query pattern (tenant pre-filter):

```sql
SELECT * FROM evidence_embeddings
WHERE tenant_id = $1
ORDER BY embedding <=> $2
LIMIT $3;
```
