package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// IntentRule 对应 intent_rules 表。
type IntentRule struct {
	ent.Schema
}

func (IntentRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "intent_rules"},
	}
}

func (IntentRule) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Positive().
			Immutable().
			Comment("主键"),
		field.String("rule_id").
			MaxLen(64).
			Unique().
			Comment("规则唯一ID"),
		field.String("intent_name").
			MaxLen(128).
			Comment("意图名称"),
		field.JSON("patterns", []string{}).
			Comment("匹配模式列表"),
		field.String("match_type").
			MaxLen(16).
			Default("contains").
			Comment("exact/prefix/contains/regex"),
		field.JSON("params", map[string]string{}).
			Optional().
			Comment("规则携带的默认参数"),
		field.Int("priority").
			Default(0).
			Comment("优先级, 越大越优先"),
		field.Float("confidence").
			Default(1.0).
			Comment("该规则的固定置信度"),
		field.Int8("status").
			Default(1).
			Comment("1=启用 0=禁用"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (IntentRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("intent_name"),
		index.Fields("status", "priority"),
	}
}

func (IntentRule) Edges() []ent.Edge {
	return nil
}
