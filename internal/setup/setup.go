package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/diezfx/split-app-backend/internal/api"
	"github.com/diezfx/split-app-backend/internal/config"
	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/diezfx/split-app-backend/pkg/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupSplitService() (*http.Server, error) {
	cfg := config.Load()
	ctx := context.Background()
	if cfg.Environment == config.LocalEnv {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger.Info(ctx).String("config", fmt.Sprint(cfg)).Msg("Loaded config")

	psqlClient, err := postgres.New(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("create sqlite client: %w", err)
	}

	storageClient, err := storage.New(ctx, psqlClient)
	if err != nil {
		return nil, fmt.Errorf("create storage client: %w", err)
	}

	projectService := service.New(storageClient)

	router := api.InitAPI(cfg, projectService)

	srv := &http.Server{
		Handler: router.Handler,
		Addr:    cfg.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv, nil
}
