package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Project struct {
	ent.Schema
}

// Fields of the BattleValues.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.String("name"),
		field.JSON("members", []string{}),
	}
}

// Edges of the BattleValues.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("transactions", Transaction.Type),
	}
}

/*
type Project struct {
	ID           string
	Name         string
	Transactions []Transaction
	Members      []string
}

*/
