package costcalc

import (
	"context"

	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/diezfx/split-app-backend/pkg/logger"
)

type Calculator struct {
	edges []Edge
}

type (
	Debitor  string
	Creditor string
)

func New(txs []*service.Transaction) *Calculator {
	return &Calculator{edges: TransformTransactionsToCostEdges(txs)}
}

func (c *Calculator) CalculateCostForUser(userID string) *money.Money {
	costSum := money.New(0, money.EUR)

	for _, tx := range c.edges {
		newSum, _ := costSum.Add(calculateBalanceForUser(tx, userID))
		costSum = newSum
	}

	return costSum
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

func calculateBalanceForUser(tx Edge, userID string) *money.Money {
	amount := money.New(0, money.EUR)
	if userID == tx.Source {
		amount = tx.Amount
	}
	if userID == tx.Target {
		amount = tx.Amount.Negative()
	}
	return amount
}
