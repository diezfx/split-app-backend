package config

import "github.com/diezfx/split-app-backend/pkg/sqlite"

type Environment string

const (
	LocalEnv       Environment = "local"
	DevelopmentEnv Environment = "dev"
)

type Config struct {
	Addr        string
	Environment Environment
	LogLevel    string

	DB sqlite.Config
}

func Load() Config {
	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		DB:          sqlite.Config{Path: "ent.db", InMemory: false},
	}
}
