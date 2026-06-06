.PHONY: dev build test lint clean docker-up docker-down docker-build migrate

# Development
dev: docker-up
	@echo "All services running. Use Ctrl+C to stop."

build:
	go build ./...

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clean:
	rm -rf dist/ vendor/

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f

# Services
run-identity:
	go run ./identity/cmd/server ./identity/identity-dev.yaml

run-questionnaire:
	go run ./questionnaire/cmd/server ./questionnaire/questionnaire-dev.yaml

run-questionnaire-worker:
	go run ./questionnaire/cmd/worker ./questionnaire/questionnaire-dev.yaml

run-evidence:
	go run ./evidence/cmd/server ./evidence/evidence-dev.yaml

run-orchestrator:
	go run ./orchestrator/cmd/server ./orchestrator/orchestrator-dev.yaml

run-orchestrator-worker:
	go run ./orchestrator/cmd/worker ./orchestrator/orchestrator-dev.yaml

run-audit:
	go run ./audit/cmd/server ./audit/audit-dev.yaml

# Database migrations
migrate-identity:
	goose -dir identity/repository/migrations postgres "$(IDENTITY_DB_URL)" up

migrate-questionnaire:
	goose -dir questionnaire/repository/migrations postgres "$(QUESTIONNAIRE_DB_URL)" up

migrate-evidence:
	goose -dir evidence/repository/migrations postgres "$(EVIDENCE_DB_URL)" up

migrate-orchestrator:
	goose -dir orchestrator/repository/migrations postgres "$(ORCHESTRATOR_DB_URL)" up

migrate-audit:
	goose -dir audit/repository/migrations postgres "$(AUDIT_DB_URL)" up

# Binaries
build-binaries:
	mkdir -p dist
	go build -o dist/identity-server ./identity/cmd/server
	go build -o dist/questionnaire-server ./questionnaire/cmd/server
	go build -o dist/questionnaire-worker ./questionnaire/cmd/worker
	go build -o dist/evidence-server ./evidence/cmd/server
	go build -o dist/orchestrator-server ./orchestrator/cmd/server
	go build -o dist/orchestrator-worker ./orchestrator/cmd/worker
	go build -o dist/audit-server ./audit/cmd/server

.PHONY: dev build test lint vet clean docker-up docker-down docker-build docker-logs
.PHONY: run-identity run-questionnaire run-questionnaire-worker run-evidence run-orchestrator run-orchestrator-worker run-audit
.PHONY: migrate-identity migrate-questionnaire migrate-evidence migrate-orchestrator migrate-audit
.PHONY: build-binaries
