package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/diezfx/split-app-backend/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Client struct {
	conn *postgres.DB
}

func New(ctx context.Context, sqlConn *postgres.DB) (*Client, error) {
	//TODO use migrate up

	client := Client{conn: sqlConn}

	err := client.conn.Up(ctx)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}

	return &client, nil
}

func (c *Client) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.amount,t.source_id,t.transaction_type,tt.user_id
	FROM projects as p
	LEFT JOIN transactions as t
	ON p.id=t.project_id 
	LEFT JOIN transaction_targets as tt
	ON t.id=tt.transaction_id
	where p.id=$1
	`
	rows, err := c.conn.QueryContext(ctx, sqlQuery, id)
	if err != nil {
		return Project{}, fmt.Errorf("query projects: %w", err)
	}

	var project = Project{ID: id, Transactions: []Transaction{}}
	rowCount := 0
	for rows.Next() {
		rowCount++
		var projectID, projectName, transactionID, transactionName, transactionSourceID, transactionType, targetUserID sql.NullString
		var transactionAmount sql.NullInt64
		err := rows.Scan(&projectID, &projectName, &transactionID, &transactionName, &transactionAmount, &transactionSourceID, &transactionType, &targetUserID)
		if err != nil {
			return project, fmt.Errorf("scan projects row: %w", err)
		}
		project.Name = projectName.String

		if !transactionID.Valid {
			return project, nil
		}

		var transaction Transaction
		transIndex := slices.IndexFunc(project.Transactions, func(t Transaction) bool {
			return transactionID.String == t.ID.String()
		})
		if transIndex == -1 {
			transaction = Transaction{
				ID:              uuid.MustParse(transactionID.String),
				Name:            transactionName.String,
				Amount:          int(transactionAmount.Int64),
				SourceID:        transactionSourceID.String,
				TargetIDs:       []string{},
				TransactionType: transactionType.String}
			project.Transactions = append(project.Transactions, transaction)
			transIndex = len(project.Transactions) - 1
		} else {
			transaction = project.Transactions[transIndex]
		}

		if targetUserID.Valid {
			transaction.TargetIDs = append(transaction.TargetIDs, targetUserID.String)
		}

		project.Transactions[transIndex] = transaction

	}
	if rowCount == 0 {
		return project, ErrNotFound
	}

	return project, nil
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	sqlQuery := `
	SELECT p.id, p.name, t.id, t.name,t.amount,t.source_id,t.transaction_type,tt.user_id
	FROM projects as p
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	LEFT JOIN transaction_targets as tt
	ON t.id=tt.transaction_id
	ORDER BY p.id
	`
	rows, err := c.conn.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}

	projects := []Project{}
	for rows.Next() {
		var projectID, projectName, transactionID, transactionName, transactionSourceID, transactionType, targetUserID sql.NullString
		var transactionAmount sql.NullInt64
		err := rows.Scan(&projectID, &projectName, &transactionID, &transactionName, &transactionAmount, &transactionSourceID, &transactionType, &targetUserID)
		if err != nil {
			return nil, fmt.Errorf("scan projects row: %w", err)
		}

		index := slices.IndexFunc(projects, func(p Project) bool {
			return projectID.String == p.ID.String()
		})
		var project Project
		if index == -1 {
			project = Project{
				ID:   uuid.MustParse(projectID.String),
				Name: projectName.String, Transactions: []Transaction{}}
			projects = append(projects, project)
			index = len(projects) - 1
		} else {
			project = projects[index]
		}
		if !transactionID.Valid {
			continue
		}

		var transaction Transaction
		transIndex := slices.IndexFunc(project.Transactions, func(t Transaction) bool {
			return transactionID.String == t.ID.String()
		})
		if transIndex == -1 {
			transaction = Transaction{
				ID:              uuid.MustParse(transactionID.String),
				Name:            transactionName.String,
				Amount:          int(transactionAmount.Int64),
				SourceID:        transactionSourceID.String,
				TargetIDs:       []string{},
				TransactionType: transactionType.String}
			project.Transactions = append(project.Transactions, transaction)
			transIndex = len(project.Transactions) - 1
		} else {
			transaction = project.Transactions[transIndex]
		}

		if targetUserID.Valid {
			transaction.TargetIDs = append(transaction.TargetIDs, targetUserID.String)
		}

		project.Transactions[transIndex] = transaction
		projects[index] = project

	}

	return projects, nil
}

func (c *Client) AddProject(ctx context.Context, proj Project) (Project, error) {

	addProjectFunc := func(ctx context.Context, tx *sql.Tx) error {

		sqlQuery := `
		insert into projects (id,name)
		values($1,$2)
		`
		_, err := tx.ExecContext(ctx, sqlQuery, proj.ID, proj.Name)
		if err != nil {
			return fmt.Errorf("insert project: %w", err)
		}

		sqlUserInsert := `
		insert into project_memberships (project_id,user_id)
		values($1,$2)
		`
		stmt, err := tx.Prepare(sqlUserInsert)
		if err != nil {
			return fmt.Errorf("prepare add project users: %w", err)
		}
		for _, user := range proj.Members {

			_, err := stmt.ExecContext(ctx, proj.ID, user)
			if err != nil {
				return fmt.Errorf("insert user: %w", err)
			}
		}
		return nil

	}

	logger.Info(ctx).Msg("this is reached")
	err := WithTransaction(ctx, c.conn.DB, addProjectFunc)
	if err != nil {
		return proj, fmt.Errorf("execute add project transaction: %w", err)
	}

	return proj, nil
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

func WithTransaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) error {
	// Begin a transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Execute the provided function within the transaction
	if err := fn(ctx, tx); err != nil {
		rbErr := tx.Rollback()
		return errors.Join(fmt.Errorf("execute function: %w", err), rbErr) // Return any error from the inner function
	}

	// Commit the transaction if everything was successful
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil // Return nil to indicate success
}
