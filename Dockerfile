# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS base
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# --- Identity Service ---
FROM base AS identity-build
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./identity/cmd/server

FROM alpine:3.20 AS identity
RUN apk add --no-cache ca-certificates tzdata
COPY --from=identity-build /server /server
COPY identity/identity-dev.yaml /configs/identity.yaml
EXPOSE 8081
ENTRYPOINT ["/server", "/configs/identity.yaml"]

# --- Questionnaire Service ---
FROM base AS questionnaire-build
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./questionnaire/cmd/server

FROM alpine:3.20 AS questionnaire
RUN apk add --no-cache ca-certificates tzdata
COPY --from=questionnaire-build /server /server
COPY questionnaire/questionnaire-dev.yaml /configs/questionnaire.yaml
EXPOSE 8082
ENTRYPOINT ["/server", "/configs/questionnaire.yaml"]

# --- Questionnaire Worker ---
FROM base AS questionnaire-worker-build
COPY . .
RUN CGO_ENABLED=0 go build -o /worker ./questionnaire/cmd/worker

FROM alpine:3.20 AS questionnaire-worker
RUN apk add --no-cache ca-certificates tzdata
COPY --from=questionnaire-worker-build /worker /worker
COPY questionnaire/questionnaire-dev.yaml /configs/questionnaire.yaml
ENTRYPOINT ["/worker", "/configs/questionnaire.yaml"]

# --- Evidence Service ---
FROM base AS evidence-build
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./evidence/cmd/server

FROM alpine:3.20 AS evidence
RUN apk add --no-cache ca-certificates tzdata
COPY --from=evidence-build /server /server
COPY evidence/evidence-dev.yaml /configs/evidence.yaml
EXPOSE 8083
ENTRYPOINT ["/server", "/configs/evidence.yaml"]

# --- Orchestrator Service ---
FROM base AS orchestrator-build
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./orchestrator/cmd/server

FROM alpine:3.20 AS orchestrator
RUN apk add --no-cache ca-certificates tzdata
COPY --from=orchestrator-build /server /server
COPY orchestrator/orchestrator-dev.yaml /configs/orchestrator.yaml
EXPOSE 8084
ENTRYPOINT ["/server", "/configs/orchestrator.yaml"]

# --- Audit Service ---
FROM base AS audit-build
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./audit/cmd/server

FROM alpine:3.20 AS audit
RUN apk add --no-cache ca-certificates tzdata
COPY --from=audit-build /server /server
COPY audit/audit-dev.yaml /configs/audit.yaml
EXPOSE 8085
ENTRYPOINT ["/server", "/configs/audit.yaml"]

# --- Orchestrator Worker ---
FROM base AS orchestrator-worker-build
COPY . .
RUN CGO_ENABLED=0 go build -o /worker ./orchestrator/cmd/worker

FROM alpine:3.20 AS orchestrator-worker
RUN apk add --no-cache ca-certificates tzdata
COPY --from=orchestrator-worker-build /worker /worker
COPY orchestrator/orchestrator-dev.yaml /configs/orchestrator.yaml
ENTRYPOINT ["/worker", "/configs/orchestrator.yaml"]
