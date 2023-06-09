package storage

import (
	"context"
	"fmt"

	"github.com/diezfx/split-app-backend/gen/ent"
	"github.com/diezfx/split-app-backend/gen/ent/project"
	"github.com/google/uuid"
)

const seedUuid = "902b0687-f61c-41c4-86dc-f7d62db6ed7d"

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

func (c *Client) GetProjectByID(ctx context.Context, id uuid.UUID) (*ent.Project, error) {
	proj, err := c.entClient.Project.Query().WithTransactions().Where(project.ID(id)).First(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return proj, nil
}

func (c *Client) Seed() error {
	ctx := context.Background()
	id := uuid.MustParse(seedUuid)

	memberList := []string{"user1", "user2"}

	// Check if the user "rotemtam" already exists.
	r, err := c.entClient.Project.Query().
		Where(project.ID(id)).
		WithTransactions().
		Only(ctx)
	// If not, create the user.
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("query project: %w", err)
	}
	if ent.IsNotFound(err) {
		r, err = c.entClient.Project.Create().
			SetID(id).
			SetName("testProj1").
			SetMembers(memberList).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create project: %w", err)
		}
	}

	if len(r.Edges.Transactions) >= 2 {
		return nil
	}

	id1 := uuid.New()
	id2 := uuid.New()

	transactions := []*ent.TransactionCreate{
		c.entClient.Transaction.Create().SetID(id1).SetName("transaction1").SetAmount(25).SetSourceID("user1").SetTargetIds([]string{"user2"}),
		c.entClient.Transaction.Create().SetID(id2).SetName("transaction2").SetAmount(100).SetSourceID("user2").SetTargetIds([]string{"user3"})}
	err = c.entClient.Transaction.CreateBulk(transactions...).Exec(ctx)
	if err != nil {
		return fmt.Errorf("store transactions: %w", err)
	}

	err = r.Update().AddTransactionIDs(id1, id2).Exec(ctx)
	if err != nil {
		return fmt.Errorf("connect transactions to project")
	}

	return nil
}
