package api

import (
	"errors"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/google/uuid"
)

type InvalidArgumentError struct {
	Argument string
}

func (a *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument %s", a.Argument)
}

func NewInvalidArgumentError(arg string) *InvalidArgumentError {
	return &InvalidArgumentError{Argument: arg}
}

type AddProject struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Members []string `json:"members,omitempty"`
}

type AddTransaction struct {
	ID              uuid.UUID `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	TransactionType string    `json:"transactionType,omitempty"`
	Amount          float64   `json:"amount,omitempty"`
	SourceID        string    `json:"sourceID,omitempty"`
	TargetIDs       []string  `json:"targetIDs,omitempty"`
}

type GetProjectsQueryParams struct{}

func (t *AddTransaction) Validate() (service.Transaction, error) {
	var err error
	if t.Name == "" {
		err = errors.Join(err, NewInvalidArgumentError("Name"))
	}
	transactionType := service.ConvertToTransactionType(t.TransactionType)
	if transactionType == service.Undefined {
		err = errors.Join(err, NewInvalidArgumentError("TransactionType"))
	}

	if t.Amount <= 0 {
		err = errors.Join(err, NewInvalidArgumentError("Amount"))
	}
	amount := money.NewFromFloat(t.Amount, money.EUR)

	if t.SourceID == "" {
		err = errors.Join(err, NewInvalidArgumentError("SourceID"))
	}

	if len(t.TargetIDs) < 1 {
		err = errors.Join(err, NewInvalidArgumentError("TargetIDs"))
	}

	return service.Transaction{
		ID:              t.ID,
		Name:            t.Name,
		TransactionType: transactionType,
		Amount:          amount,
		SourceID:        t.SourceID,
		TargetIDs:       t.TargetIDs,
	}, err
}

type Transaction struct {
	ID              uuid.UUID               `json:"id,omitempty"`
	Name            string                  `json:"name,omitempty"`
	TransactionType service.TransactionType `json:"transactionType,omitempty"`
	Amount          float64                 `json:"amount,omitempty"`
	SourceID        string                  `json:"sourceId,omitempty"`
	TargetIDs       []string                `json:"targetIds,omitempty"`
}

func TransactionFromServiceTransaction(t service.Transaction) Transaction {
	return Transaction{
		ID:              t.ID,
		Name:            t.Name,
		TransactionType: t.TransactionType,
		Amount:          t.Amount.AsMajorUnits(),
		SourceID:        t.SourceID,
		TargetIDs:       t.TargetIDs,
	}
}

type Project struct {
	ID           uuid.UUID     `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Members      []string      `json:"members,omitempty"`
}

func ProjectFromServiceProject(p service.Project) Project {
	transactions := make([]Transaction, 0, len(p.Transactions))
	for _, t := range p.Transactions {
		transactions = append(transactions, TransactionFromServiceTransaction(t))
	}
	return Project{ID: p.ID, Name: p.Name, Transactions: transactions, Members: p.Members}
}
