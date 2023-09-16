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
	GetProjectUsers(ctx context.Context, projectID uuid.UUID) ([]service.User, error)
	AddProjectUser(ctx context.Context, projID uuid.UUID, userID string) error
	GetCostsByUser(ctx context.Context, userID string) (service.UserCosts, error)
	GetCostsByProject(ctx context.Context, projID uuid.UUID) (service.ProjectCosts, error)
	AddTransaction(ctx context.Context, projID uuid.UUID, transaction service.Transaction) error
}

type UserService interface {
	GetUsers(ctx context.Context) ([]service.User, error)
}
