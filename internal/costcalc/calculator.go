package costcalc

import (
	"context"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/pkg/logger"
)

type Calculator struct {
	edges []Edge
}

type (
	Debitor  string
	Creditor string
)

func New(txs []Transaction) *Calculator {
	return &Calculator{edges: TransformTransactionsToCostEdges(txs)}
}

func (c *Calculator) CalculateCostForUser(userID string) (*Cost, error) {
	cost := &Cost{
		Expenses: money.New(0, money.EUR),
		Income:   money.New(0, money.EUR),
		Balance:  money.New(0, money.EUR)}

	for _, tx := range c.edges {
		// add to expense when source
		if tx.Source == userID {
			newExpense, err := cost.Expenses.Add(tx.Amount)
			if err != nil {
				return nil, fmt.Errorf("add source: %w", err)
			}
			cost.Expenses = newExpense
		}
		// add to income when target
		if tx.Target == userID {
			newIncome, err := cost.Income.Add(tx.Amount)
			if err != nil {
				return nil, fmt.Errorf("add target: %w", err)
			}
			cost.Income = newIncome
		}
	}

	newBalance, err := cost.Expenses.Subtract(cost.Income)
	if err != nil {
		return nil, fmt.Errorf("calculate new user balance: %w", err)
	}
	cost.Balance = newBalance

	return cost, nil
}

func (c *Calculator) CalculateCostForAllUsers() (*ProjectCost, error) {
	totalCost := money.New(0, money.EUR)

	userCosts := map[string]*Cost{}

	for _, tx := range c.edges {
		newCost, err := totalCost.Add(tx.Amount)
		if err != nil {
			return nil, fmt.Errorf("add to total cost: %w", err)
		}
		totalCost = newCost

		source := userCosts[tx.Source]
		if source == nil {
			source = &Cost{Expenses: money.New(0, money.EUR), Income: money.New(0, money.EUR), Balance: money.New(0, money.EUR)}
			userCosts[tx.Source] = source
		}
		target := userCosts[tx.Target]
		if target == nil {
			target = &Cost{Expenses: money.New(0, money.EUR), Income: money.New(0, money.EUR), Balance: money.New(0, money.EUR)}
			userCosts[tx.Target] = target
		}

		newExpenses, err := source.Expenses.Add(tx.Amount)
		if err != nil {
			return nil, fmt.Errorf("add to source expenses: %w", err)
		}
		source.Expenses = newExpenses
		newIncome, err := target.Income.Add(tx.Amount)
		if err != nil {
			return nil, fmt.Errorf("add to source expenses: %w", err)
		}
		target.Income = newIncome
	}

	// as last step calculate balance
	for _, c := range userCosts {
		newBalance, err := c.Expenses.Subtract(c.Income)
		if err != nil {
			return nil, fmt.Errorf("calculate new user balance: %w", err)
		}
		c.Balance = newBalance
	}

	return &ProjectCost{TotalCost: totalCost, CostPerUser: userCosts}, nil
}

func (c *Calculator) CalculateMinCostFlow() []Edge {
	optEdges := []Edge{}

	userBalances := calculateBalances(c.edges)

	for {
		maxDebitor := getMaxDebitor(userBalances)
		maxCreditor := getMaxCreditor(userBalances)

		if maxDebitor == maxCreditor || userBalances[maxDebitor].Amount() == 0 || userBalances[maxCreditor].Amount() == 0 {
			logger.Debug(context.Background()).Msg("maxDebitor and creditor are the same person -> all debts are settled")
			return optEdges
		}

		amount := userBalances[maxDebitor]
		optEdges = append(optEdges, Edge{Source: maxDebitor, Target: maxCreditor, Amount: amount.Absolute()})
		newVal, _ := userBalances[maxDebitor].Subtract(amount)
		userBalances[maxDebitor] = newVal
		newVal, _ = userBalances[maxCreditor].Add(amount)
		userBalances[maxCreditor] = newVal
	}
}

func calculateBalances(edges []Edge) map[string]*money.Money {
	userBalance := map[string]*money.Money{}

	for _, edge := range edges {
		if userBalance[edge.Source] == nil {
			userBalance[edge.Source] = money.New(0, money.EUR)
		}
		if userBalance[edge.Target] == nil {
			userBalance[edge.Target] = money.New(0, money.EUR)
		}

		newMoneyVal, _ := userBalance[edge.Source].Add(edge.Amount)
		userBalance[edge.Source] = newMoneyVal
		newMoneyVal, _ = userBalance[edge.Target].Subtract(edge.Amount)
		userBalance[edge.Target] = newMoneyVal
	}

	return userBalance
}

func getMaxDebitor(userBalance map[string]*money.Money) string {
	debitor := ""

	for user, balance := range userBalance {
		if debitor == "" {
			debitor = user
		}
		if cmp, _ := balance.Compare(userBalance[debitor]); cmp == -1 {
			debitor = user
		}
	}

	return debitor
}

func getMaxCreditor(userBalance map[string]*money.Money) string {
	creditor := ""

	for user, balance := range userBalance {
		if creditor == "" {
			creditor = user
		}
		if cmp, _ := balance.Compare(userBalance[creditor]); cmp == 1 {
			creditor = user
		}
	}
	return creditor
}
