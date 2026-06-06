package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	NATS     NATSConfig
	LLM      LLMConfig
	Embedder EmbedderConfig
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

type NATSConfig struct {
	URL string
}

type LLMConfig struct {
	Provider    string  // "openai", "claude", "local"
	OpenAIKey   string  `mapstructure:"openai_key"`
	ClaudeKey   string  `mapstructure:"claude_key"`
	LocalURL    string  `mapstructure:"local_url"`
	LocalModel  string  `mapstructure:"local_model"`
	Temperature float64
	MaxTokens   int     `mapstructure:"max_tokens"`
}

type EmbedderConfig struct {
	Provider  string // "openai", "local"
	OpenAIKey string `mapstructure:"openai_key"`
	LocalURL  string `mapstructure:"local_url"`
	Model     string
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("EVIDRA")

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8084)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "60s")
	v.SetDefault("server.shutdown_timeout", "30s")

	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("nats.url", "nats://localhost:4222")

	v.SetDefault("llm.provider", "openai")
	v.SetDefault("llm.temperature", 0.3)
	v.SetDefault("llm.max_tokens", 2048)

	v.SetDefault("embedder.provider", "openai")
	v.SetDefault("embedder.model", "text-embedding-3-small")

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
