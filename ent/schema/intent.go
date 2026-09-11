package schema

import "entgo.io/ent"

// Intent holds the schema definition for the Intent entity.
type Intent struct {
	ent.Schema
}

// Fields of the Intent.
func (Intent) Fields() []ent.Field {
	return nil
}

// Edges of the Intent.
func (Intent) Edges() []ent.Edge {
	return nil
}
