package service

import (
	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type TransactionType string

const (
	UndefinedTransactionType TransactionType = "Undefined"
	ExpenseTransactionType   TransactionType = "Expense"
	TransferTransactionType  TransactionType = "Transfer"
)

func ParseTransactionType(trans string) TransactionType {
	switch trans {
	case string(ExpenseTransactionType):
		return ExpenseTransactionType
	case string(TransferTransactionType):
		return TransferTransactionType
	default:
		return UndefinedTransactionType
	}
}

type User struct {
	ID string
}

type Transaction struct {
	ID              uuid.UUID
	Name            string
	TransactionType TransactionType
	Amount          *money.Money
	SourceID        string
	TargetIDs       []string
}

type Project struct {
	ID           uuid.UUID
	Name         string
	Transactions []Transaction
	Members      []string
}

func FromStorageProject(project storage.Project) Project {
	transactions := make([]Transaction, len(project.Transactions))
	for i, t := range project.Transactions {
		transactions[i] = FromStorageTransaction(t)
	}

	return Project{
		ID:           project.ID,
		Name:         project.Name,
		Transactions: transactions,
		Members:      project.Members,
	}
}

func ToStorageProject(proj Project) storage.Project {
	transactions := make([]storage.Transaction, len(proj.Transactions))
	for i, t := range proj.Transactions {
		transactions[i] = ToStorageTransaction(t)
	}

	return storage.Project{
		ID:           proj.ID,
		Name:         proj.Name,
		Transactions: transactions,
		Members:      proj.Members,
	}
}

func ToStorageTransaction(trans Transaction) storage.Transaction {
	return storage.Transaction{
		ID:   trans.ID,
		Name: trans.Name, Amount: int(trans.Amount.Amount()),
		SourceID:        trans.SourceID,
		TargetIDs:       trans.TargetIDs,
		TransactionType: string(trans.TransactionType),
	}
}

func FromStorageTransaction(trans storage.Transaction) Transaction {
	return Transaction{
		ID:   trans.ID,
		Name: trans.Name, Amount: money.New(int64(trans.Amount), money.EUR),
		SourceID:        trans.SourceID,
		TargetIDs:       trans.TargetIDs,
		TransactionType: ParseTransactionType(trans.TransactionType),
	}
}
