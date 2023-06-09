package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// BattleValues holds the schema definition for the BattleValues entity.
type Transaction struct {
	ent.Schema
}

// Fields of the BattleValues.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.String("name"),
		field.Float("amount"),
		field.String("source_id"),
		field.JSON("target_ids", []string{}),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type),
	}
}

/*
type Transaction struct {
	ID              string
	Name            string
	TransactionType TransactionType
	Amount          float64
	SourceID        string
	TargetIDs       []string
}

*/
