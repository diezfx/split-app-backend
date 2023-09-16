package costcalc

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type Edge struct {
	Source string
	Target string
	Amount *money.Money
}

func TransformTransactionsToCostEdges(txs []Transaction) []Edge {
	edges := []Edge{}

	for _, tx := range txs {
		splitVals, _ := tx.Amount.Split(len(tx.TargetIDs))
		for i, splitValue := range splitVals {
			edges = append(edges, Edge{Source: tx.SourceID, Target: tx.TargetIDs[i], Amount: splitValue})
		}
	}

	return edges
}

type Transaction struct {
	ProjectID uuid.UUID
	ID        uuid.UUID
	Amount    *money.Money
	SourceID  string
	TargetIDs []string
}

type Cost struct {
	Expenses *money.Money
	Income   *money.Money
	Balance  *money.Money
}

type ProjectCost struct {
	TotalCost *money.Money

	CostPerUser map[string]*Cost
}
