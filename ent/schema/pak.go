package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Pak holds the schema definition for the Pak entity.
type Pak struct {
	ent.Schema
}

// Fields of the Pak.
func (Pak) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique(),
		field.String("name"),
		field.String("path"),
		field.Bool("active"),
	}
}

// Edges of the Pak.
func (Pak) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modVersion", ModVersion.Type).
			Ref("paks").
			Unique(),
	}
}
