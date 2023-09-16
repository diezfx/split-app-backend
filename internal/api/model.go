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
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type AddTransaction struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	TransactionType string   `json:"transactionType"`
	Amount          float64  `json:"amount"`
	SourceID        string   `json:"sourceId"`
	TargetIDs       []string `json:"targetIds"`
}

type GetProjectsQueryParams struct{}

func (t *AddTransaction) Validate() (service.Transaction, error) {
	var err error

	id, err := uuid.Parse(t.ID)
	if err != nil {
		err = errors.Join(err, NewInvalidArgumentError("ID"))
	}

	if t.Name == "" {
		err = errors.Join(err, NewInvalidArgumentError("Name"))
	}
	transactionType := service.ParseTransactionType(t.TransactionType)
	if transactionType == service.UndefinedTransactionType {
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
		ID:              id,
		Name:            t.Name,
		TransactionType: transactionType,
		Amount:          amount,
		SourceID:        t.SourceID,
		TargetIDs:       t.TargetIDs,
	}, err
}

type Transaction struct {
	ID              uuid.UUID               `json:"id"`
	Name            string                  `json:"name"`
	TransactionType service.TransactionType `json:"transactionType"`
	Amount          float64                 `json:"amount"`
	SourceID        string                  `json:"sourceId"`
	TargetIDs       []string                `json:"targetIds"`
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
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Transactions []Transaction `json:"transactions"`
	Members      []string      `json:"members"`
}

func ProjectFromServiceProject(p service.Project) Project {
	transactions := make([]Transaction, 0, len(p.Transactions))
	for _, t := range p.Transactions {
		transactions = append(transactions, TransactionFromServiceTransaction(t))
	}
	return Project{ID: p.ID, Name: p.Name, Transactions: transactions, Members: p.Members}
}

type User struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	ErrorCode int
	Reason    string
}

type UserCosts struct {
	TotalCosts   Cost               `json:"totalCosts"`
	ProjectCosts map[uuid.UUID]Cost `json:"projectCosts"`
}

type ProjectCosts struct {
	TotalCosts  float64         `json:"totalCosts"`
	CostsByUser map[string]Cost `json:"costsByUser"`
}

type Cost struct {
	Expenses float64 `json:"expenses"`
	Income   float64 `json:"income"`
	Balance  float64 `json:"balance"`
}

func UserCostsFromService(c service.UserCosts) UserCosts {
	projectCosts := make(map[uuid.UUID]Cost, len(c.ProjectCosts))
	for p, pc := range c.ProjectCosts {
		projectCosts[p] = CostFromService(pc)
	}

	return UserCosts{
		TotalCosts:   CostFromService(c.TotalCost),
		ProjectCosts: projectCosts,
	}
}

func ProjectCostsFromService(cost service.ProjectCosts) ProjectCosts {
	userCosts := make(map[string]Cost, len(cost.UserCosts))
	for u, c := range cost.UserCosts {
		userCosts[u] = CostFromService(c)
	}

	return ProjectCosts{
		TotalCosts:  cost.TotalCost.AsMajorUnits(),
		CostsByUser: userCosts,
	}
}

func CostFromService(c service.Cost) Cost {
	return Cost{
		Expenses: c.Expenses.AsMajorUnits(),
		Income:   c.Income.AsMajorUnits(),
		Balance:  c.Balance.AsMajorUnits(),
	}
}
