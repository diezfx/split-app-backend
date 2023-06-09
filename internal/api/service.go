package api

import (
	"context"

	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/google/uuid"
)

type ProjectService interface {
	GetProject(ctx context.Context, id uuid.UUID) (service.Project, error)
}
