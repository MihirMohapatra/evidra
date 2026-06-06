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

	"github.com/evidra/evidra/identity/internal/config"
	"github.com/evidra/evidra/identity/repository/postgres"
	"github.com/evidra/evidra/identity/service"
	"github.com/evidra/evidra/identity/transport"
)

func main() {
	cfgPath := "configs/dev.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	pool.Config().MaxConns = int32(cfg.Database.MaxOpenConns)
	pool.Config().MinConns = int32(cfg.Database.MaxIdleConns)

	// Repositories
	orgRepo := postgres.NewOrganizationRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	sessRepo := postgres.NewSessionRepo(pool)
	keyRepo := postgres.NewAPIKeyRepo(pool)

	// Service
	svcCfg := service.Config{
		JWTSecret:         cfg.JWT.Secret,
		JWTIssuer:         cfg.JWT.Issuer,
		SessionTTL:        cfg.JWT.SessionTTL,
		APIKeyLength:      cfg.JWT.APIKeyLength,
		PasswordMinLength: 8,
	}
	svc := service.New(orgRepo, userRepo, sessRepo, keyRepo, svcCfg)

	// Router
	router := transport.NewRouter(svc)

	// Server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("identity service starting", "addr", addr)
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
