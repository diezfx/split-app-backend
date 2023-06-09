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

	Db sqlite.Config
}

func Load() Config {
	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		Db:          sqlite.Config{Path: "ent.db", InMemory: false},
	}
}
