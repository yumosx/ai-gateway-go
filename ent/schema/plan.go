package schema

import "entgo.io/ent"

type Plan struct {
	ent.Schema
}

// Fields of the Intent.
func (Plan) Fields() []ent.Field {
	return nil
}

// Edges of the Intent.
func (Plan) Edges() []ent.Edge {
	return nil
}
