package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ModVersion holds the schema definition for the ModVersion entity.
type ModVersion struct {
	ent.Schema
}

// Fields of the ModVersion.
func (ModVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
	}
}

// Edges of the ModVersion.
func (ModVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mod", Mod.Type).
			Ref("versions").Unique(),
		edge.To("archives", Archive.Type).Required(),
		edge.To("paks", Pak.Type).Required(),
	}
}
