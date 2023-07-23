package service

import (
	"context"

	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type ProjectStorage interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (storage.Project, error)
	AddProject(ctx context.Context, project storage.Project) error
}
