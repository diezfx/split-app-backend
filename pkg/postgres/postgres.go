package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Port     int
	Host     string
	Database string
	Username string
	Password string

	MigrationsDir string
}

const connectionString = "postgres://%s:%s@%s:%d/%s"

type DB struct {
	cfg Config
	*sql.DB
}

func New(cfg Config) (*DB, error) {
	connString := buildConnectionString(cfg)
	logger.Info(context.Background()).String("connectionString", connString).Msg("connectionString")
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	return &DB{cfg: cfg, DB: conn}, nil
}

func buildConnectionString(cfg Config) string {
	return fmt.Sprintf(connectionString, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func (db *DB) Up(ctx context.Context) error {

	driver, err := migratepgx.WithInstance(db.DB, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", db.cfg.MigrationsDir), db.cfg.Database, driver)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil

}
