package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Mod holds the schema definition for the Mod entity.
type Mod struct {
	ent.Schema
}

// Fields of the Mod.
// A mod is an entry of a specific mod release.
// Every mod has a list of versions.
// A mod can only have one active version at a time.
// Of the mod's state is "inactive" no version should be active
func (Mod) Fields() []ent.Field {
	return []ent.Field{
		// Display name of the overall mod.
		// This is not to be confused with the archive name
		field.String("name").Unique(),
		// determines if the mod is considered "installed"
		field.Enum("state").
			Values("activate", "inactive"),
		// activeVersion records the current version of the installed mod
		// when the "state" field is inactive "activeVersion" must be nil
		field.UUID("activeVersion", uuid.UUID{}).
			Optional().
			Nillable(),
		// A reference to where the mod originated from
		// This could be a nexus url.
		field.String("origin").
			Optional().
			Nillable(),
	}
}

// Edges of the Mod.
func (Mod) Edges() []ent.Edge {
	return []ent.Edge{
		// versions should be a one to many relationship
		// "versions" contains all possible versions of a mod
		edge.To("versions", ModVersion.Type),
	}
}
