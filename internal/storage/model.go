package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	Name            string
	TransactionType string
	Amount          int
	SourceID        string
	TargetIDs       []string
}

type projectQueryElement struct {
	ProjectID       sql.NullString
	ProjectName     sql.NullString
	TransactionID   sql.NullString
	TransactionName sql.NullString
	TransactionType sql.NullString
	Amount          sql.NullInt64
	SourceID        sql.NullString
	TargetID        sql.NullString
}

type transactionQueryElement struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	Name            string
	TransactionType string
	Amount          int
	SourceID        string
	TargetID        string
}

type Project struct {
	ID           uuid.UUID
	Name         string
	Transactions []Transaction
	Members      []string
}
