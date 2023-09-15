package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/diezfx/split-app-backend/pkg/postgres"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Client struct {
	conn *postgres.DB
}

func New(ctx context.Context, sqlConn *postgres.DB) (*Client, error) {
	client := Client{conn: sqlConn}

	err := client.conn.Up(ctx)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}

	return &client, nil
}

func (c *Client) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	sqlQuery := `
	SELECT p.id as project_id, p.name as project_name,
		t.id as transaction_id, t.name as transaction_name,t.amount,t.source_id,t.transaction_type,tt.user_id as target_id
	FROM projects as p
	LEFT JOIN transactions as t
	ON p.id=t.project_id 
	LEFT JOIN transaction_targets as tt
	ON t.id=tt.transaction_id
	where p.id=$1
	`
	var projectQueryElements []projectQueryElement

	err := sqlscan.Select(ctx, c.conn.DB, &projectQueryElements, sqlQuery, id)
	if err != nil {
		return Project{}, fmt.Errorf("select queryElements: %w", err)
	}
	projects := mergeProject(projectQueryElements)
	if len(projects) == 0 {
		return Project{}, ErrNotFound
	}
	return projects[0], nil
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	sqlQuery := `
	SELECT p.id as project_id, p.name as project_name,
		t.id as transaction_id, t.name as transaction_name,t.amount,t.source_id,t.transaction_type,tt.user_id as target_id
	FROM projects as p
	LEFT JOIN transactions as t
	ON p.id=t.project_id
	LEFT JOIN transaction_targets as tt
	ON t.id=tt.transaction_id
	ORDER BY p.id
	`
	var projectQueryElements []projectQueryElement

	err := sqlscan.Select(ctx, c.conn.DB, &projectQueryElements, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("select queryElements: %w", err)
	}
	projects := mergeProject(projectQueryElements)
	return projects, nil
}

func (c *Client) AddProjectUser(ctx context.Context, projectID uuid.UUID, userID string) error {
	err := withTransaction(ctx, c.conn.DB, func(ctx context.Context, tx *sql.Tx) error {
		return addUsers(ctx, tx, projectID, []string{userID})
	})
	if err != nil {
		return fmt.Errorf("addUser: %w", err)
	}
	return nil
}

func addUsers(ctx context.Context, tx *sql.Tx, projectID uuid.UUID, userIDs []string) error {
	sqlUserInsert := `insert into project_memberships (project_id,user_id)
	values($1,$2)
	`
	stmt, err := tx.PrepareContext(ctx, sqlUserInsert)
	if err != nil {
		return fmt.Errorf("prepare add project users: %w", err)
	}
	for _, user := range userIDs {
		_, err := stmt.ExecContext(ctx, projectID, user)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}
	}
	return nil
}

func (c *Client) AddProject(ctx context.Context, proj Project) (Project, error) {
	addProjectFunc := func(ctx context.Context, tx *sql.Tx) error {
		sqlQuery := `insert into projects (id,name)
		values($1,$2)
		`
		_, err := tx.ExecContext(ctx, sqlQuery, proj.ID, proj.Name)
		if err != nil {
			return fmt.Errorf("insert project: %w", err)
		}
		return addUsers(ctx, tx, proj.ID, proj.Members)
	}

	err := withTransaction(ctx, c.conn.DB, addProjectFunc)
	if err != nil {
		return proj, fmt.Errorf("execute add project transaction: %w", err)
	}

	return proj, nil
}

func (c *Client) AddTransaction(ctx context.Context, projectID uuid.UUID, transaction Transaction) error {
	addTransactionFunc := func(ctx context.Context, tx *sql.Tx) error {
		const sqlQuery = `
		INSERT INTO transactions (id,name,amount,source_id,transaction_type,project_id)
		VALUES($1,$2,$3,$4,$5,$6)
		`
		_, err := tx.ExecContext(ctx, sqlQuery,
			transaction.ID, transaction.Name, transaction.Amount, transaction.SourceID, transaction.TransactionType, projectID)
		if err != nil {
			return fmt.Errorf("insert project: %w", err)
		}

		const insertTransactionTargetsQuery = `
		INSERT INTO transaction_targets (transaction_id,user_id)
		VALUES($1,$2)`

		stmt, err := tx.Prepare(insertTransactionTargetsQuery)
		if err != nil {
			return fmt.Errorf("prepare add project users: %w", err)
		}
		for _, target := range transaction.TargetIDs {
			_, err := stmt.ExecContext(ctx, transaction.ID, target)
			if err != nil {
				return fmt.Errorf("insert target: %w", err)
			}
		}
		return nil
	}

	return withTransaction(ctx, c.conn.DB, addTransactionFunc)
}

func (c *Client) GetAllOutgoingTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
	sqlQuery := `
	SELECT t.id, t.name, t.source_id,tt.user_id as target_id, t.transaction_type,t.project_id,t.amount
	FROM transactions as t
	LEFT JOIN transaction_targets as tt
	ON t.id=tt.transaction_id 
	WHERE source_id=$1
	`
	var transactionElements []transactionQueryElement
	err := sqlscan.Select(ctx, c.conn.DB, &transactionElements, sqlQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}

	logger.Info(ctx).String("transactions", fmt.Sprint(transactionElements)).Msg("test")
	// merge transactions
	transactions := mergeTransactionElements(transactionElements)
	return transactions, nil
}

func (c *Client) GetAllIncomingTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
	sqlQuery := `
	SELECT t.id, t.name, t.source_id,tt.user_id as target_id, t.transaction_type,t.project_id,t.amount
	FROM transactions as t
	JOIN transaction_targets as tt
	ON t.id=tt.transaction_id 
	WHERE tt.user_id=$1
	`
	var transactionElements []transactionQueryElement
	err := sqlscan.Select(ctx, c.conn.DB, &transactionElements, sqlQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}

	// merge transactions
	transactions := mergeTransactionElements(transactionElements)
	return transactions, nil
}

func (c *Client) GetProjectUsers(ctx context.Context, projectID uuid.UUID) ([]User, error) {
	sqlQuery := `
	SELECT user_id as id
	FROM project_memberships
	WHERE project_id=$1
	`
	var users []User
	err := sqlscan.Select(ctx, c.conn.DB, &users, sqlQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (c *Client) GetUser(ctx context.Context, userID string) (User, error) {
	sqlQuery := `
	SELECT id
	FROM members
	WHERE id=$1
	`
	var user User
	err := sqlscan.Get(ctx, c.conn.DB, &user, sqlQuery, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("select users: %w", err)
	}

	return user, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	sqlQuery := `
	SELECT id
	FROM members
	`
	var users []User
	err := sqlscan.Select(ctx, c.conn.DB, &users, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (c *Client) AddUser(ctx context.Context, user User) error {
	sqlQuery := `
	INSERT INTO members (id)
	VALUES ($1)
	`
	_, err := c.conn.DB.ExecContext(ctx, sqlQuery, user.ID)
	if err != nil {
		return fmt.Errorf("insert member: %w", err)
	}

	return nil
}

func withTransaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) error {
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

func mergeTransactionElements(transactionElements []transactionQueryElement) []Transaction {
	transactions := []Transaction{}
	for _, te := range transactionElements {
		index := slices.IndexFunc(transactions, func(t Transaction) bool {
			return t.ID == te.ID
		})
		var transaction Transaction
		if index == -1 {
			transaction = Transaction{
				ID:              te.ID,
				Name:            te.Name,
				Amount:          te.Amount,
				TransactionType: te.TransactionType,
				SourceID:        te.SourceID,
				ProjectID:       te.ProjectID,
				TargetIDs:       []string{},
			}
			transactions = append(transactions, transaction)
			index = len(transactions) - 1
		} else {
			transaction = transactions[index]
		}

		if te.TargetID != "" {
			transaction.TargetIDs = append(transaction.TargetIDs, te.TargetID)
		}
		transactions[index] = transaction
	}
	return transactions
}

func mergeProject(projectElements []projectQueryElement) []Project {
	projects := []Project{}

	for i := 0; i < len(projectElements); i++ {
		pe := projectElements[i]
		index := slices.IndexFunc(projects, func(p Project) bool {
			return pe.ProjectID.String == p.ID.String()
		})
		var project Project
		if index == -1 {
			project = Project{
				ID:   uuid.MustParse(pe.ProjectID.String),
				Name: pe.ProjectName.String, Transactions: []Transaction{},
			}
			projects = append(projects, project)
			index = len(projects) - 1
		} else {
			project = projects[index]
		}
		if !pe.TransactionID.Valid {
			continue
		}

		var transaction Transaction
		transIndex := slices.IndexFunc(project.Transactions, func(t Transaction) bool {
			return pe.TransactionID.String == t.ID.String()
		})
		if transIndex == -1 {
			transaction = Transaction{
				ID:              uuid.MustParse(pe.TransactionID.String),
				Name:            pe.TransactionName.String,
				Amount:          int(pe.Amount.Int64),
				SourceID:        pe.SourceID.String,
				TargetIDs:       []string{},
				TransactionType: pe.TransactionType.String,
			}
			project.Transactions = append(project.Transactions, transaction)
			transIndex = len(project.Transactions) - 1
		} else {
			transaction = project.Transactions[transIndex]
		}

		if pe.TargetID.Valid {
			transaction.TargetIDs = append(transaction.TargetIDs, pe.TargetID.String)
		}

		project.Transactions[transIndex] = transaction
		projects[index] = project
	}
	return projects
}
