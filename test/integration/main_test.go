package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	identityPool *pgxpool.Pool
	evidencePool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	identityPool = startPostgres(ctx, "evidra_identity", "5433")
	evidencePool = startPostgres(ctx, "evidra_evidence", "5434")

	if identityPool == nil || evidencePool == nil {
		log.Fatal("failed to start postgres containers")
	}

	runMigrations(identityPool, "identity")
	runMigrations(evidencePool, "evidence")

	code := m.Run()

	if identityPool != nil {
		identityPool.Close()
	}
	if evidencePool != nil {
		evidencePool.Close()
	}

	os.Exit(code)
}

func startPostgres(ctx context.Context, dbName, port string) *pgxpool.Pool {
	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername("evidra"),
		postgres.WithPassword("evidra"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Printf("failed to start postgres container for %s: %v", dbName, err)
		return nil
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("failed to get connection string for %s: %v", dbName, err)
		return nil
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Printf("failed to create connection pool for %s: %v", dbName, err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("failed to ping database %s: %v", dbName, err)
		return nil
	}

	return pool
}

func runMigrations(pool *pgxpool.Pool, service string) {
	ctx := context.Background()
	var migrations []string

	switch service {
	case "identity":
		migrations = []string{
			identityMigrations["001_create_organizations"],
			identityMigrations["002_create_users"],
			identityMigrations["003_create_sessions"],
			identityMigrations["004_create_api_keys"],
			identityMigrations["005_create_oidc_states"],
			identityMigrations["006_create_linked_accounts"],
		}
	case "evidence":
		migrations = []string{
			evidenceMigrations["001_create_evidence_items"],
			evidenceMigrations["002_create_approvals"],
		}
	}

	for i, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			log.Fatalf("migration %d for %s failed: %v", i+1, service, err)
		}
	}
	fmt.Printf("migrations for %s applied successfully\n", service)
}
