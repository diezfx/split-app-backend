package storage

import (
	"context"
	"fmt"

	"github.com/diezfx/split-app-backend/gen/ent"
	"github.com/diezfx/split-app-backend/gen/ent/project"
	"github.com/diezfx/split-app-backend/gen/ent/transaction"
	"github.com/google/uuid"
)

type Client struct {
	entClient *ent.Client
}

func New(entClient *ent.Client) (*Client, error) {
	if err := entClient.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("create schema resources: %w", err)
	}
	client := Client{entClient: entClient}

	err := client.Seed()
	if err != nil {
		return nil, fmt.Errorf("seed db: %w", err)
	}
	return &client, nil
}

func (c *Client) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	proj, err := c.entClient.Project.Query().WithTransactions().Where(project.ID(id)).First(ctx)
	if ent.IsNotFound(err) {
		return Project{}, ErrNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("get project by id: %w", err)
	}

	return FromEntProject(proj), nil
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	projs, err := c.entClient.Project.Query().WithTransactions().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	var projectList []Project
	for _, p := range projs {
		projectList = append(projectList, FromEntProject(p))
	}
	return projectList, nil
}

func (c *Client) AddProject(ctx context.Context, proj Project) (Project, error) {
	result, err := c.entClient.Project.Create().
		SetID(proj.ID).
		SetName(proj.Name).
		SetMembers(proj.Members).Save(ctx)
	if err != nil {
		return Project{}, fmt.Errorf("add project to db: %w", err)
	}

	return FromEntProject(result), nil
}

func (c *Client) AddTransaction(ctx context.Context, projectID uuid.UUID, tx Transaction) error {
	_, err := c.entClient.Transaction.Create().
		SetID(tx.ID).
		SetName(tx.Name).
		SetAmount(tx.Amount.Amount()).
		SetSourceID(tx.SourceID).SetTransactionType(tx.TransactionType).
		SetTargetIds(tx.TargetIDs).SetProjectID(projectID).Save(ctx)
	if err != nil {
		return fmt.Errorf("add transaction to db: %w", err)
	}
	return nil
}

func (c *Client) GetAllOutgoingTransactionsByUserID(ctx context.Context, userID string) ([]*ent.Transaction, error) {
	txs, err := c.entClient.Transaction.Query().
		Where(transaction.And(transaction.SourceID(userID))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all outgoing edges: %w", err)
	}
	return txs, nil
}

func (c *Client) GetAllIncomingTransactionsByUserID(ctx context.Context, projectID uuid.UUID, userID string) ([]*ent.Transaction, error) {
	txs, err := c.entClient.Transaction.Query().
		Where(transaction.And(transaction.SourceID(userID), transaction.HasProjectWith(project.ID(projectID)))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all outgoing edges: %w", err)
	}
	return txs, nil
}
