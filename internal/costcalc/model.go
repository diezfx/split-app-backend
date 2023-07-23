package costcalc

import (
	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/service"
)

type Edge struct {
	Source string
	Target string
	Amount *money.Money
}

func TransformTransactionsToCostEdges(txs []*service.Transaction) []Edge {
	edges := []Edge{}

	for _, tx := range txs {
		splitVals, _ := tx.Amount.Split(len(tx.TargetIDs))
		for i, splitValue := range splitVals {
			edges = append(edges, Edge{Source: tx.SourceID, Target: tx.TargetIDs[i], Amount: splitValue})
		}
	}

	return edges
}
