package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/evidra/evidra/audit/domain"
	"github.com/evidra/evidra/audit/events"
	"github.com/evidra/evidra/audit/internal/config"
	"github.com/evidra/evidra/audit/repository/postgres"
	"github.com/evidra/evidra/audit/service"
	"github.com/evidra/evidra/audit/transport"
	"github.com/evidra/evidra/pkg/queue"
)

func main() {
	cfgPath := "audit-dev.yaml"
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

	repo := postgres.NewAuditRepo(pool)
	svc := service.New(repo)
	router := transport.NewRouter(svc)

	natsBus, err := queue.NewNATS(queue.NATSConfig{URL: cfg.NATS.URL})
	if err != nil {
		slog.Warn("nats not available, running without event subscription", "error", err)
	} else {
		if err := natsBus.Subscribe(events.SubjectAuditRecorded, func(ctx context.Context, data []byte) error {
			var evt events.AuditRecorded
			if err := json.Unmarshal(data, &evt); err != nil {
				slog.Error("failed to unmarshal audit event", "error", err)
				return nil
			}
			_, err := svc.Record(ctx, service.RecordInput{
				TenantID: evt.TenantID,
				ActorID:  evt.ActorID,
				Action:   domain.Action(evt.Action),
				TargetID: evt.TargetID,
			})
			if err != nil {
				slog.Error("failed to record audit event", "error", err)
			}
			return nil
		}); err != nil {
			slog.Error("failed to subscribe to audit events", "error", err)
		}
	}

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
		slog.Info("audit service starting", "addr", addr)
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
