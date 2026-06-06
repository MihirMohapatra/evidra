.PHONY: dev build run-api run-worker migrate test lint

dev:
	@echo "Starting dev..."

build:
	go build ./...

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

migrate:
	go run ./cmd/migration

test:
	go test ./...

lint:
	golangci-lint run ./...
