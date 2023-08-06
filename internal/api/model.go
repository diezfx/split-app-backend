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
	ID      string   `json:"ID,omitempty"`
	Name    string   `json:"Name,omitempty"`
	Members []string `json:"Members,omitempty"`
}

type AddTransaction struct {
	ID              uuid.UUID `json:"ID,omitempty"`
	Name            string    `json:"Name,omitempty"`
	TransactionType string    `json:"TransactionType,omitempty"`
	Amount          float64   `json:"Amount,omitempty"`
	SourceID        string    `json:"SourceID,omitempty"`
	TargetIDs       []string  `json:"TargetIDs,omitempty"`
}

func (t *AddTransaction) Validate() (service.Transaction, error) {
	var err error
	if t.Name == "" {
		errors.Join(err, NewInvalidArgumentError("Name"))
	}
	transactionType := service.ConvertToTransactionType(t.TransactionType)
	if transactionType == service.Undefined {
		errors.Join(err, NewInvalidArgumentError("TransactionType"))
	}

	if t.Amount <= 0 {
		errors.Join(err, NewInvalidArgumentError("Amount"))
	}
	amount := money.NewFromFloat(t.Amount, money.EUR)

	if t.SourceID == "" {
		errors.Join(err, NewInvalidArgumentError("SourceID"))
	}

	if len(t.TargetIDs) < 1 {
		errors.Join(err, NewInvalidArgumentError("TargetIDs"))
	}

	return service.Transaction{
		ID:              t.ID,
		Name:            t.Name,
		TransactionType: transactionType,
		Amount:          amount,
		SourceID:        t.SourceID,
		TargetIDs:       t.TargetIDs}, err

}
