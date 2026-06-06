package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/evidence/events"
	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/internal/config"
	"github.com/evidra/evidra/orchestrator/repository/postgres"
	"github.com/evidra/evidra/orchestrator/service"
	"github.com/evidra/evidra/pkg/queue"
)

func main() {
	cfgPath := "orchestrator-dev.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	natsBus, err := queue.NewNATS(queue.NATSConfig{URL: cfg.NATS.URL})
	if err != nil {
		slog.Error("failed to connect to nats", "error", err)
		os.Exit(1)
	}
	defer func() { _ = natsBus.Close() }()

	var embedder service.Embedder
	switch cfg.Embedder.Provider {
	case "openai":
		key := cfg.Embedder.OpenAIKey
		if key == "" {
			key = cfg.LLM.OpenAIKey
		}
		embedder = service.NewOpenAIEmbedder(key, cfg.Embedder.Model)
	default:
		slog.Error("unsupported embedder provider", "provider", cfg.Embedder.Provider)
		os.Exit(1)
	}

	embeddingRepo := postgres.NewEmbeddingRepo(pool)

	slog.Info("embedding worker starting", "provider", cfg.Embedder.Provider)

	if err := natsBus.Subscribe(events.SubjectEvidenceCreated, func(ctx context.Context, data []byte) error {
		var event events.EvidenceCreated
		if err := json.Unmarshal(data, &event); err != nil {
			slog.Error("failed to unmarshal event", "error", err)
			return err
		}

		go processEvidence(ctx, embedder, embeddingRepo, event)
		return nil
	}); err != nil {
		slog.Error("failed to subscribe", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("worker shutting down")
}

func processEvidence(ctx context.Context, embedder service.Embedder, repo *postgres.EmbeddingRepo, event events.EvidenceCreated) {
	log := slog.With("evidence_id", event.ID, "tenant_id", event.TenantID)

	text := event.Title
	if event.Content != "" {
		text = event.Title + "\n\n" + event.Content
	}

	vec, err := embedder.GenerateEmbedding(ctx, text)
	if err != nil {
		log.Error("failed to generate embedding", "error", err)
		return
	}

	chunk := &domain.EvidenceChunk{
		ID:         uuid.New(),
		TenantID:   event.TenantID,
		EvidenceID: event.ID,
		Content:    text,
		Embedding:  vec,
		Metadata: map[string]any{
			"title":    event.Title,
			"category": event.Category,
		},
	}

	if err := repo.Upsert(ctx, chunk); err != nil {
		log.Error("failed to upsert embedding", "error", err)
		return
	}

	log.Info("embedding created successfully")
}
