package sqlite

import (
	"fmt"

	"github.com/diezfx/split-app-backend/gen/ent"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Path     string
	InMemory bool
}

const connectionString = "file:%s?mode=%s&cache=shared&_fk=1"

func NewSqliteDB(cfg Config) (*ent.Client, error) {

	client, err := ent.Open("sqlite3", buildConnectionString(cfg))
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	return client, nil
}

func buildConnectionString(cfg Config) string {
	if cfg.InMemory {
		return fmt.Sprintf(connectionString, uuid.New(), "memory")
	}
	return fmt.Sprintf(connectionString, cfg.Path, "rwc")
}
