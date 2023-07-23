package service

import (
	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/gen/ent/transaction"
	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type TransactionType string

const (
	Undefined TransactionType = "Undefined"
	Expense   TransactionType = "Expense"
	Transfer  TransactionType = "Transfer"
)

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
	transactionValue := transaction.TransactionTypeExpense
	switch trans.TransactionType {
	case Expense:
		transactionValue = transaction.TransactionTypeExpense
	case Transfer:
		transactionValue = transaction.TransactionTypeTransfer
	}

	return storage.Transaction{
		ID:   trans.ID,
		Name: trans.Name, Amount: trans.Amount,
		SourceID:        trans.SourceID,
		TargetIDs:       trans.TargetIDs,
		TransactionType: transactionValue,
	}
}

func FromStorageTransaction(trans storage.Transaction) Transaction {
	transactionValue := Undefined
	switch trans.TransactionType {
	case transaction.TransactionTypeExpense:
		transactionValue = Expense
	case transaction.TransactionTypeTransfer:
		transactionValue = Transfer
	}

	return Transaction{
		ID:   trans.ID,
		Name: trans.Name, Amount: trans.Amount,
		SourceID:        trans.SourceID,
		TargetIDs:       trans.TargetIDs,
		TransactionType: transactionValue,
	}
}
