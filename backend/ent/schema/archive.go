package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Archive holds the schema definition for the Archive entity.
type Archive struct {
	ent.Schema
}

// Fields of the Archive.
func (Archive) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("archivePath").Unique(),
		field.String("md5Sum"),
		field.Bool("installed").Default(false),
	}
}

// Edges of the Archive.
func (Archive) Edges() []ent.Edge {
	return nil
}
