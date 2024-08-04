package schema

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Mod holds the schema definition for the Mod entity.
type Mod struct {
	ent.Schema
}

// Fields of the Mod.
func (Mod) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("state").Match(regexp.MustCompile("^(active|inactive)$")),
	}
}

// Edges of the Mod.
func (Mod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("archive", Archive.Type).Unique().Required(),
		edge.To("paks", Pak.Type),
	}
}
