package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/identity/internal/config"
	"github.com/evidra/evidra/identity/repository/postgres"
	"github.com/evidra/evidra/identity/service"
	"github.com/evidra/evidra/identity/transport"
	grpctransport "github.com/evidra/evidra/identity/transport/grpc"
	"github.com/evidra/evidra/pkg/queue"
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
	oidcStateRepo := postgres.NewOIDCStateRepo(pool)
	linkedAccountRepo := postgres.NewLinkedAccountRepo(pool)

	// OIDC providers
	oidcProviders := make([]service.ProviderConfig, 0, len(cfg.OIDC.Providers))
	for _, p := range cfg.OIDC.Providers {
		scopes := p.Scopes
		if len(scopes) == 0 {
			scopes = []string{"openid", "email", "profile"}
		}
		oidcProviders = append(oidcProviders, service.ProviderConfig{
			Name:         p.Name,
			IssuerURL:    p.IssuerURL,
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			RedirectURL:  p.RedirectURL,
			Scopes:       scopes,
		})
	}

	// NATS
	natsBus, err := queue.NewNATS(queue.NATSConfig{URL: cfg.NATS.URL})
	if err != nil {
		slog.Warn("nats not available, running without events", "error", err)
		natsBus = nil
	}

	// Service
	svcCfg := service.Config{
		JWTSecret:         cfg.JWT.Secret,
		JWTIssuer:         cfg.JWT.Issuer,
		SessionTTL:        cfg.JWT.SessionTTL,
		APIKeyLength:      cfg.JWT.APIKeyLength,
		PasswordMinLength: 8,
	}
	svc := service.New(orgRepo, userRepo, sessRepo, keyRepo, oidcStateRepo, linkedAccountRepo, svcCfg, oidcProviders, natsBus)

	// HTTP router
	router := transport.NewRouter(svc)

	httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpSrv := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// gRPC server
	grpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port+1000)
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		slog.Error("failed to listen for grpc", "error", err)
		os.Exit(1)
	}

	grpcSrv := grpc.NewServer()
	evidrav1.RegisterIdentityServiceServer(grpcSrv, grpctransport.NewServer(svc))
	reflection.Register(grpcSrv)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("identity http server starting", "addr", httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		slog.Info("identity grpc server starting", "addr", grpcAddr)
		if err := grpcSrv.Serve(grpcLis); err != nil {
			slog.Error("grpc server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	grpcSrv.GracefulStop()
	_ = httpSrv.Shutdown(ctx)

	slog.Info("server stopped")
}
