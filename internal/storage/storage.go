package storage

import (
	"context"
	"fmt"

	"github.com/diezfx/split-app-backend/pkg/postgres"
	"github.com/google/uuid"
)

type Client struct {
	conn *postgres.DB
}

func New(ctx context.Context, sqlConn *postgres.DB) (*Client, error) {
	//TODO use migrate up

	client := Client{conn: sqlConn}

	err := client.conn.Up(ctx)
	if err != nil {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}
	err = client.Seed()
	if err != nil {
		return nil, fmt.Errorf("seed db: %w", err)
	}
	return &client, nil
}

func (c *Client) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)

	return Project{}, nil
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)
	return nil, nil
}

func (c *Client) AddProject(ctx context.Context, proj Project) (Project, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)

	return Project{}, nil
}

func (c *Client) AddTransaction(ctx context.Context, projectID uuid.UUID, tx Transaction) error {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)
	return nil
}

func (c *Client) GetAllOutgoingTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)
	return nil, nil
}

func (c *Client) GetAllIncomingTransactionsByUserID(ctx context.Context, projectID uuid.UUID, userID string) ([]Transaction, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.transaction_type  
	FROM projects as p,
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	`
	c.conn.QueryRowContext(ctx, sqlQuery)
	return nil, nil
}
