package config

import (
	"os"

	"github.com/diezfx/split-app-backend/pkg/postgres"
)

type Environment string

const (
	LocalEnv       Environment = "local"
	DevelopmentEnv Environment = "dev"
)

type Config struct {
	Addr        string
	Environment Environment
	LogLevel    string

	DB postgres.Config
}

func Load() Config {
	env := os.Getenv("ENVIRONMENT")
	if env == string(DevelopmentEnv) {
		return Config{
			Addr:        ":8080",
			Environment: DevelopmentEnv,
			LogLevel:    "debug",
			DB:          postgres.Config{Port: 5432, Host: "localhost", Database: "postgres", Username: "postgres", Password: "postgres", MigrationsDir: "db/migrations"},
		}
	}

	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		DB:          postgres.Config{Port: 5432, Host: "localhost", Database: "postgres", Username: "postgres", Password: "postgres", MigrationsDir: "db/migrations"},
	}
}
