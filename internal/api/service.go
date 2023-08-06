package api

import (
	"context"

	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/google/uuid"
)

type ProjectService interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (service.Project, error)
	GetProjects(ctx context.Context) ([]service.Project, error)
	AddProject(ctx context.Context, proj service.Project) (service.Project, error)

	AddTransaction(ctx context.Context, projID uuid.UUID, transaction service.Transaction) error
}
