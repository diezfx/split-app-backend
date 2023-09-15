package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Service struct {
	projStorage ProjectStorage
}

// AddProjectUser implements api.ProjectService.
func (s *Service) AddProjectUser(ctx context.Context, projID uuid.UUID, userID string) error {
	if _, err := s.projStorage.GetUser(ctx, userID); err != nil {
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("get user: %w", err)
		}
		if err != nil && errors.Is(err, storage.ErrNotFound) {
			err := s.projStorage.AddUser(ctx, storage.User{ID: userID})
			if err != nil {
				return fmt.Errorf("add user:%w", err)
			}
		}
	}

	err := s.projStorage.AddProjectUser(ctx, projID, userID)
	if err != nil {
		return fmt.Errorf("add project user: %w", err)
	}
	return err
}

// GetProjectUsers implements api.ProjectService.
func (s *Service) GetProjectUsers(ctx context.Context, projectID uuid.UUID) ([]User, error) {
	sUsers, err := s.projStorage.GetProjectUsers(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("getProjectUsers: %w", err)
	}

	users := make([]User, 0, len(sUsers))

	for _, u := range sUsers {
		users = append(users, User{ID: u.ID})
	}
	return users, nil
}

// AddTransaction implements api.ProjectService.
func (s *Service) AddTransaction(ctx context.Context, projID uuid.UUID, transaction Transaction) error {
	_, err := s.projStorage.GetProjectByID(ctx, projID)
	if errors.Is(err, storage.ErrNotFound) {
		return ErrProjectNotFound
	}
	if err != nil {
		return fmt.Errorf("get project:%w", err)
	}
	err = s.projStorage.AddTransaction(ctx, projID, ToStorageTransaction(transaction))
	if err != nil {
		return fmt.Errorf("add transaction: %w", err)
	}
	return nil
}

func New(projStorage ProjectStorage) *Service {
	return &Service{projStorage: projStorage}
}

func (s *Service) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	proj, err := s.projStorage.GetProjectByID(ctx, id)
	if errors.Is(err, storage.ErrNotFound) {
		return Project{}, ErrProjectNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("get project:%w", err)
	}
	return FromStorageProject(proj), nil
}

func (s *Service) GetProjects(ctx context.Context) ([]Project, error) {
	projs, err := s.projStorage.GetProjects(ctx)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project:%w", err)
	}

	projectList := make([]Project, 0, len(projs))
	for _, p := range projs {
		projectList = append(projectList, FromStorageProject(p))
	}

	return projectList, nil
}

func (s *Service) AddProject(ctx context.Context, project Project) (Project, error) {
	_, err := s.projStorage.GetProjectByID(ctx, project.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return Project{}, fmt.Errorf("add project: %w", err)
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return Project{}, fmt.Errorf("add project: %w", err)
	}

	users, err := s.projStorage.GetUsers(ctx)
	if err != nil {
		return Project{}, fmt.Errorf("get users: %w", err)
	}

	for _, member := range project.Members {
		if slices.IndexFunc(users, func(u storage.User) bool { return member == u.ID }) == -1 {
			err = s.projStorage.AddUser(ctx, storage.User{ID: member})
			if err != nil {
				return Project{}, fmt.Errorf("add new user for project: %w", err)
			}
		}
	}

	proj, err := s.projStorage.AddProject(ctx, ToStorageProject(project))
	if err != nil {
		return Project{}, fmt.Errorf("add project: %w", err)
	}

	return FromStorageProject(proj), nil
}

// GetCostsByUser implements api.ProjectService.
// TODO improve implementation use pointers
func (s *Service) GetCostsByUser(ctx context.Context, userID string) (UserCosts, error) {
	incomeSt, err := s.projStorage.GetAllIncomingTransactionsByUserID(ctx, userID)
	if err != nil {
		return UserCosts{}, fmt.Errorf("get incoming transactions: %w", err)
	}
	income := make([]Transaction, 0, len(incomeSt))
	for _, tx := range incomeSt {
		income = append(income, FromStorageTransaction(tx))
	}

	outgoingSt, err := s.projStorage.GetAllOutgoingTransactionsByUserID(ctx, userID)
	if err != nil {
		return UserCosts{}, fmt.Errorf("get outgoing transactions: %w", err)
	}
	outgoing := make([]Transaction, 0, len(outgoingSt))
	for _, tx := range outgoingSt {
		fmt.Printf("\n%+v\n", tx)
		outgoing = append(outgoing, FromStorageTransaction(tx))
	}

	totalCost := Cost{
		Expenses: money.New(0, money.EUR),
		Income:   money.New(0, money.EUR),
		Balance:  money.New(0, money.EUR),
	}
	projectCosts := map[uuid.UUID]Cost{}

	for _, tx := range income {
		if _, ok := projectCosts[tx.ProjectID]; !ok {
			projectCosts[tx.ProjectID] = Cost{
				Expenses: money.New(0, money.EUR),
				Income:   money.New(0, money.EUR),
				Balance:  money.New(0, money.EUR),
			}
		}
		totalIncome, _ := totalCost.Income.Add(tx.Amount)
		totalCost.Income = totalIncome

		projectCost := projectCosts[tx.ProjectID]
		newProjectIncome, _ := projectCost.Income.Add(tx.Amount)
		projectCost.Income = newProjectIncome
		projectCosts[tx.ProjectID] = projectCost
	}

	for _, tx := range outgoing {
		if _, ok := projectCosts[tx.ProjectID]; !ok {
			projectCosts[tx.ProjectID] = Cost{
				Expenses: money.New(0, money.EUR),
				Income:   money.New(0, money.EUR),
				Balance:  money.New(0, money.EUR),
			}
		}
		totalExpenses, _ := totalCost.Expenses.Add(tx.Amount)
		totalCost.Expenses = totalExpenses

		projectCost := projectCosts[tx.ProjectID]
		newProjectExpenses, _ := projectCost.Expenses.Add(tx.Amount)
		projectCost.Expenses = newProjectExpenses
		projectCosts[tx.ProjectID] = projectCost

		if err != nil {
			return UserCosts{}, fmt.Errorf("add to totalExpeses: %w", err)
		}
	}

	balance, err := totalCost.Expenses.Subtract(totalCost.Income)
	if err != nil {
		return UserCosts{}, fmt.Errorf("calc balance: %w", err)
	}

	for pID, c := range projectCosts {
		balance, _ = c.Expenses.Subtract(c.Income)
		c.Balance = balance
		projectCosts[pID] = c
	}
	totalCost.Balance = balance
	userCosts := UserCosts{
		TotalCost:    totalCost,
		ProjectCosts: projectCosts,
	}
	return userCosts, nil
}
