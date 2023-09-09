package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/diezfx/split-app-backend/pkg/configloader"
	"github.com/diezfx/split-app-backend/pkg/postgres"
)

const defaultNamespace = "postgres"

func LoadPostgresConfig(loader configloader.Loader) (postgres.Config, error) {
	cfg := postgres.Config{}

	content, err := loader.LoadConfig(defaultNamespace)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshal postgres: %w", err)
	}

	username, err := loader.LoadSecret(defaultNamespace, "username")
	if err != nil {
		return cfg, err
	}
	password, err := loader.LoadSecret(defaultNamespace, "password")
	if err != nil {
		return cfg, err
	}
	cfg.Username = username
	cfg.Password = password
	return cfg, nil
}
