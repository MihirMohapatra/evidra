package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	NATS     NATSConfig
	App      AppConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type StorageConfig struct {
	Bucket       string
	Endpoint     string
	AccessKey    string
	SecretKey    string
	UseSSL       bool
	Region       string
	MaxFileSize  int64
}

type NATSConfig struct {
	URL string
}

type AppConfig struct {
	WorkerConcurrency int
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("EVIDRA")

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8082)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.shutdown_timeout", "30s")

	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("storage.bucket", "questionnaires")
	v.SetDefault("storage.endpoint", "localhost:9000")
	v.SetDefault("storage.access_key", "minioadmin")
	v.SetDefault("storage.secret_key", "minioadmin")
	v.SetDefault("storage.use_ssl", false)
	v.SetDefault("storage.region", "us-east-1")
	v.SetDefault("storage.max_file_size", 52428800)

	v.SetDefault("nats.url", "nats://localhost:4222")

	v.SetDefault("app.worker_concurrency", 5)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
