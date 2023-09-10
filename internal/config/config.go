package config

import (
	"os"

	postgrescfg "github.com/diezfx/split-app-backend/internal/config/postgres"
	supabasecfg "github.com/diezfx/split-app-backend/internal/config/supabase"
	"github.com/diezfx/split-app-backend/pkg/auth"
	"github.com/diezfx/split-app-backend/pkg/configloader"
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
	Auth        auth.Config
	DB          postgres.Config
}

func Load() (Config, error) {
	env := os.Getenv("ENVIRONMENT")

	// read from stuff json
	loader := configloader.NewFileLoader("/etc/config", "/etc/secrets")

	// read postgres secrets

	if env == string(DevelopmentEnv) {
		pgDB, err := postgrescfg.LoadPostgresConfig(loader)
		if err != nil {
			return Config{}, err
		}

		authCfg, err := supabasecfg.LoadSupabaseConfig(loader)
		if err != nil {
			return Config{}, err
		}

		return Config{
			Addr:        ":8080",
			Environment: DevelopmentEnv,
			LogLevel:    "debug",
			DB:          pgDB,
			Auth:        authCfg,
		}, nil
	}

	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		DB: postgres.Config{
			Port: 5432, Host: "localhost", Database: "postgres",
			Username: "postgres", Password: "postgres",
			MigrationsDir: "db/migrations",
		},
	}, nil
}

func (cfg *Config) IsLocal() bool {
	return cfg.Environment == LocalEnv
}
