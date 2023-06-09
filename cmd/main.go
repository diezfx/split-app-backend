package main

import (
	"context"

	"github.com/diezfx/split-app-backend/internal/setup"
	"github.com/diezfx/split-app-backend/pkg/logger"
)

func main() {
	srv, err := setup.SetupSplitService()
	if err != nil {
		logger.Fatal(context.Background(), err).Msg("failed setup")
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(context.Background(), err).Msg("failed listening and serve")
	}
}
