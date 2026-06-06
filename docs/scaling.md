# Scaling

## Service-Level Scaling

Each service is stateless and horizontally scalable. Add instances behind a load balancer.

| Service | State | Scalability |
|---------|-------|-------------|
| Identity | Stateless (JWT sessions in DB) | Horizontal |
| Questionnaire | Stateless (files in S3) | Horizontal |
| Questionnaire Worker | Stateless (NATS consumer) | Horizontal (partitioned) |
| Evidence | Stateless | Horizontal |
| Orchestrator | Stateless (LLM calls) | Horizontal |
| Orchestrator Worker | Stateless (NATS consumer) | Horizontal (partitioned) |
| Audit | Append-only writer | Horizontal (partitioned) |

## Database Scaling

Each service owns an independent PostgreSQL database:

- **Read replicas**: Add read replicas for query-heavy services (identity, audit)
- **Connection pooling**: pgx pool configured with `MaxConns: 25`, `MinConns: 5` per instance
- **pgvector**: IVFFlat index on `evidence_embeddings.embedding`; increase `lists` for larger datasets

## NATS Scaling

NATS JetStream is used for async event processing:

- **Queue groups**: Workers in the same queue group compete for messages (competing consumer pattern)
- **Subjects**: `evidence.*`, `questionnaire.*`, `orchestrator.*`

## gRPC Internal Communication

Services communicate via gRPC on port +1000:

- Low-latency binary protocol (protobuf)
- Reflection enabled for tooling (grpcurl, grpcui)
- Multiplexed HTTP/2 connections reduce connection overhead

## Cache Strategy

Not implemented yet. Future considerations:

- Redis for session caching (identity)
- Redis for pgvector query result caching (orchestrator)
- CDN for file downloads (questionnaire, evidence)

## Throughput Estimates

| Service | Bottleneck | Estimated Capacity (per instance) |
|---------|-----------|----------------------------------|
| Identity | DB writes (session create) | ~5000 req/s |
| Questionnaire | File parsing (CPU) | ~100 req/s |
| Evidence | DB reads | ~3000 req/s |
| Orchestrator | LLM API latency | ~10 req/s |
| Audit | DB append writes | ~10000 events/s |
| Workers | Embedding API latency | ~50 docs/s |

## Infrastructure Scaling (Future)

- **Kubernetes**: Horizontal Pod Autoscaler (HPA) based on CPU/memory
- **Terraform**: Infrastructure-as-Code for cloud deployments
- **Service mesh**: Istio or Linkerd for mTLS and traffic management
