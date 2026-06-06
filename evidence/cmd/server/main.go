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
	"github.com/evidra/evidra/evidence/internal/config"
	"github.com/evidra/evidra/evidence/repository/postgres"
	"github.com/evidra/evidra/evidence/service"
	"github.com/evidra/evidra/evidence/transport"
	grpctransport "github.com/evidra/evidra/evidence/transport/grpc"
	"github.com/evidra/evidra/pkg/queue"
)

func main() {
	cfgPath := "evidence-dev.yaml"
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

	evidenceRepo := postgres.NewEvidenceRepo(pool)
	approvalRepo := postgres.NewApprovalRepo(pool)

	svc := service.New(evidenceRepo, approvalRepo, natsBus)

	router := transport.NewRouter(svc)

	httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpSrv := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	grpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port+1000)
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		slog.Error("failed to listen for grpc", "error", err)
		os.Exit(1)
	}

	grpcSrv := grpc.NewServer()
	evidrav1.RegisterEvidenceServiceServer(grpcSrv, grpctransport.NewServer(svc))
	reflection.Register(grpcSrv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("evidence http server starting", "addr", httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		slog.Info("evidence grpc server starting", "addr", grpcAddr)
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
