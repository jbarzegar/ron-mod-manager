package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Pak holds the schema definition for the Pak entity.
type Pak struct {
	ent.Schema
}

// Fields of the Pak.
func (Pak) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("installed").Default(false),
	}
}

// Edges of the Pak.
func (Pak) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mod", Mod.Type).
			Ref("paks").Unique(),
	}
}
