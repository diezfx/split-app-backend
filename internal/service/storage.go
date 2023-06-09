package service

import (
	"context"

	"github.com/diezfx/split-app-backend/gen/ent"
	"github.com/google/uuid"
)

type ProjectStorage interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (*ent.Project, error)
}
