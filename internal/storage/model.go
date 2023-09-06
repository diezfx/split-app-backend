package storage

import (
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

type Project struct {
	ID           uuid.UUID
	Name         string
	Transactions []Transaction
	Members      []string
}
