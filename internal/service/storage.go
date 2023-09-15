package service

import (
	"context"

	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type ProjectStorage interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (storage.Project, error)
	GetProjects(ctx context.Context) ([]storage.Project, error)
	GetProjectUsers(ctx context.Context, projectID uuid.UUID) ([]storage.User, error)
	AddProject(ctx context.Context, project storage.Project) (storage.Project, error)
	AddTransaction(ctx context.Context, projectID uuid.UUID, transaction storage.Transaction) error
	GetUsers(ctx context.Context) ([]storage.User, error)
	GetUser(ctx context.Context, userID string) (storage.User, error)
	AddUser(ctx context.Context, user storage.User) error
	AddProjectUser(ctx context.Context, projectID uuid.UUID, userID string) error

	GetAllOutgoingTransactionsByUserID(ctx context.Context, userID string) ([]storage.Transaction, error)
	GetAllIncomingTransactionsByUserID(ctx context.Context, userID string) ([]storage.Transaction, error)
}
