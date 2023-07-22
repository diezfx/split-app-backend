package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	projStorage ProjectStorage
}

func New(projStorage ProjectStorage) *Service {
	return &Service{projStorage: projStorage}
}

func (s *Service) GetProject(ctx context.Context, id uuid.UUID) (Project, error) {
	proj, err := s.projStorage.GetProjectByID(ctx, id)
	if errors.Is(err, storage.ErrNotFound) {
		return Project{}, ErrProjectNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("get project:%w", err)
	}
	return FromEntProject(proj), nil
}
