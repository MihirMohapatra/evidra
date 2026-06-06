package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/orchestrator/internal/config"
	"github.com/evidra/evidra/orchestrator/repository/postgres"
	"github.com/evidra/evidra/orchestrator/service"
	"github.com/evidra/evidra/orchestrator/transport"
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
		slog.Warn("nats not available, running without events", "error", err)
		natsBus = nil
	}

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

	var engine service.AIEngine
	model := cfg.LLM.LocalModel
	switch cfg.LLM.Provider {
	case "openai":
		if model == "" {
			model = "gpt-4o"
		}
		engine = service.NewOpenAIEngine(cfg.LLM.OpenAIKey, model, cfg.LLM.Temperature, cfg.LLM.MaxTokens)
	case "claude":
		if model == "" {
			model = "claude-sonnet-4-20250514"
		}
		engine = service.NewClaudeEngine(cfg.LLM.ClaudeKey, model, cfg.LLM.Temperature, cfg.LLM.MaxTokens)
	case "local":
		if model == "" {
			model = "llama3"
		}
		engine = service.NewLocalLLMEngine(cfg.LLM.LocalURL, model, cfg.LLM.Temperature, cfg.LLM.MaxTokens)
	default:
		slog.Error("unsupported LLM provider", "provider", cfg.LLM.Provider)
		os.Exit(1)
	}

	embeddingRepo := postgres.NewEmbeddingRepo(pool)
	draftRepo := postgres.NewDraftRepo(pool)

	svc := service.New(embedder, engine, embeddingRepo, draftRepo, natsBus, 5)

	router := transport.NewRouter(svc)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("orchestrator service starting", "addr", addr, "llm", cfg.LLM.Provider, "embedder", cfg.Embedder.Provider)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
