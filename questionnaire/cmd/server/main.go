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

	"github.com/evidra/evidra/pkg/queue"
	"github.com/evidra/evidra/pkg/storage"
	"github.com/evidra/evidra/questionnaire/internal/config"
	"github.com/evidra/evidra/questionnaire/repository/postgres"
	"github.com/evidra/evidra/questionnaire/service"
	"github.com/evidra/evidra/questionnaire/transport"
)

func main() {
	cfgPath := "questionnaires-dev.yaml"
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

	store, err := storage.NewMinIO(storage.MinIOConfig{
		Endpoint:  cfg.Storage.Endpoint,
		AccessKey: cfg.Storage.AccessKey,
		SecretKey: cfg.Storage.SecretKey,
		UseSSL:    cfg.Storage.UseSSL,
		Region:    cfg.Storage.Region,
	})
	if err != nil {
		slog.Error("failed to create storage", "error", err)
		os.Exit(1)
	}

	natsBus, err := queue.NewNATS(queue.NATSConfig{URL: cfg.NATS.URL})
	if err != nil {
		slog.Error("failed to connect to nats", "error", err)
		os.Exit(1)
	}
	defer func() { _ = natsBus.Close() }()

	qRepo := postgres.NewQuestionnaireRepo(pool)
	quRepo := postgres.NewQuestionRepo(pool)

	svc := service.New(qRepo, quRepo, store, natsBus, cfg.Storage.Bucket, cfg.Storage.MaxFileSize)

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
		slog.Info("questionnaire service starting", "addr", addr)
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
