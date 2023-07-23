package costcalc

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/internal/service"
)

func TestCalculator(t *testing.T) {
	const userID = "test_user1"

	tests := []struct {
		name         string
		transactions []*service.Transaction

		expectedAmount *money.Money
	}{{
		name:           "outgoing",
		transactions:   []*service.Transaction{{SourceID: userID, TargetIDs: []string{"u2", "u3"}, Amount: money.New(24, money.EUR)}},
		expectedAmount: money.New(24, money.EUR),
	}, {
		name:           "incoming",
		transactions:   []*service.Transaction{{SourceID: "u2", TargetIDs: []string{userID, "u3"}, Amount: money.New(24, money.EUR)}},
		expectedAmount: money.New(-12, money.EUR),
	}, {
		name: "complex",
		transactions: []*service.Transaction{
			{SourceID: userID, TargetIDs: []string{"u1", "u4", "u3"}, Amount: money.New(25, money.EUR)},
			{SourceID: userID, TargetIDs: []string{"u1", "u3"}, Amount: money.New(25, money.EUR)},
			{SourceID: "u2", TargetIDs: []string{userID, "u3", "u4", "u5"}, Amount: money.New(100, money.EUR)},
			{SourceID: "u2", TargetIDs: []string{userID}, Amount: money.New(25, money.EUR)},
		},
		expectedAmount: money.New(0, money.EUR),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calculator := New(test.transactions)
			result := calculator.CalculateCostForUser(userID)

			if same, _ := result.Equals(test.expectedAmount); !same {
				t.Errorf("expected amount %s got %s", test.expectedAmount.Display(), result.Display())
			}
		})
	}
}

func TestMinCostFlow(t *testing.T) {
	const userID = "test_user1"

	tests := []struct {
		name         string
		transactions []*service.Transaction

		expectedCashFlow []Edge
	}{{
		name:         "outgoing",
		transactions: []*service.Transaction{{SourceID: userID, TargetIDs: []string{"u2", "u3"}, Amount: money.New(24, money.EUR)}},
		expectedCashFlow: []Edge{
			{Source: "u2", Target: userID, Amount: money.New(12, money.EUR)},
			{Source: "u3", Target: userID, Amount: money.New(12, money.EUR)},
		},
	}, {
		name:         "incoming",
		transactions: []*service.Transaction{{SourceID: "u2", TargetIDs: []string{userID, "u3"}, Amount: money.New(24, money.EUR)}},
		expectedCashFlow: []Edge{
			{Source: userID, Target: "u2", Amount: money.New(12, money.EUR)},
			{Source: "u3", Target: "u2", Amount: money.New(12, money.EUR)},
		},
	}, {
		name: "complex",
		transactions: []*service.Transaction{
			{SourceID: userID, TargetIDs: []string{"u1", "u4", "u3"}, Amount: money.New(25, money.EUR)},
			{SourceID: userID, TargetIDs: []string{"u1", "u3"}, Amount: money.New(25, money.EUR)},
			{SourceID: "u2", TargetIDs: []string{userID, "u3", "u4", "u5"}, Amount: money.New(100, money.EUR)},
			{SourceID: "u2", TargetIDs: []string{userID}, Amount: money.New(25, money.EUR)},
		},
		expectedCashFlow: []Edge{
			{Source: "u3", Target: "u2", Amount: money.New(45, money.EUR)},
			{Source: "u4", Target: "u2", Amount: money.New(33, money.EUR)},
			{Source: "u5", Target: "u2", Amount: money.New(25, money.EUR)},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calculator := New(test.transactions)
			resultFlow := calculator.CalculateMinCostFlow()

			compareEdges(t, resultFlow, test.expectedCashFlow)
		})
	}
}

func compareEdges(t *testing.T, actual, expected []Edge) {
	// find same edge
	for _, expectedEdge := range expected {
		foundEdge := false
		for _, actualEdge := range actual {
			if expectedEdge.Source == actualEdge.Source && expectedEdge.Target == actualEdge.Target {
				if expectedEdge.Amount.Amount() != actualEdge.Amount.Amount() {
					t.Errorf("expected amount %s for edge (%s:%s) got %s", expectedEdge.Amount.Display(), expectedEdge.Source, expectedEdge.Target, actualEdge.Amount.Display())
				}
				foundEdge = true
				break
			}
		}
		if !foundEdge {
			t.Errorf("expected edge (%s:%s) not found", expectedEdge.Source, expectedEdge.Target)
		}
	}
}
