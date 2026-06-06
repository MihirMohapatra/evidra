package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/pkg/queue"
	"github.com/evidra/evidra/pkg/storage"
	"github.com/evidra/evidra/questionnaire/events"
	"github.com/evidra/evidra/questionnaire/internal/config"
	"github.com/evidra/evidra/questionnaire/repository/postgres"
	"github.com/evidra/evidra/questionnaire/service"
)

func main() {
	cfgPath := "questionnaire-dev.yaml"
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

	slog.Info("questionnaire worker starting", "concurrency", cfg.App.WorkerConcurrency)

	if err := natsBus.Subscribe(events.SubjectQuestionnaireUploaded, func(ctx context.Context, data []byte) error {
		var event events.QuestionnaireUploaded
		if err := json.Unmarshal(data, &event); err != nil {
			slog.Error("failed to unmarshal event", "error", err)
			return err
		}

		slog.Info("processing document", "questionnaire_id", event.ID)
		return svc.ProcessDocument(ctx, event)
	}); err != nil {
		slog.Error("failed to subscribe", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("worker shutting down")
}
