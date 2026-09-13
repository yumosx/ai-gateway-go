package dao

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ai-gateway-go/ent"
	"github.com/ecodeclub/ai-gateway-go/ent/intentrule"
)

// IntentRuleDAO 基于 ent 的意图规则持久化（接受/返回 schema）
type IntentRuleDAO struct {
	client *ent.Client
}

func NewIntentRuleDAO(client *ent.Client) *IntentRuleDAO {
	return &IntentRuleDAO{client: client}
}

func (d *IntentRuleDAO) GetAll(ctx context.Context) ([]*ent.IntentRule, error) {
	rows, err := d.client.IntentRule.Query().
		Where(intentrule.StatusEQ(1)).
		Order(ent.Desc(intentrule.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query rules: %w", err)
	}
	return rows, nil
}

func (d *IntentRuleDAO) GetByIntentName(ctx context.Context, name string) ([]*ent.IntentRule, error) {
	rows, err := d.client.IntentRule.Query().
		Where(
			intentrule.IntentNameEQ(name),
			intentrule.StatusEQ(1),
		).
		Order(ent.Desc(intentrule.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query rules by intent: %w", err)
	}
	return rows, nil
}

func (d *IntentRuleDAO) Create(ctx context.Context, row *ent.IntentRule) error {
	builder := d.client.IntentRule.Create().
		SetRuleID(row.RuleID).
		SetIntentName(row.IntentName).
		SetPatterns(row.Patterns).
		SetMatchType(row.MatchType).
		SetPriority(row.Priority).
		SetConfidence(row.Confidence).
		SetStatus(1)

	if row.Params != nil {
		builder.SetParams(row.Params)
	}

	_, err := builder.Save(ctx)
	return err
}

func (d *IntentRuleDAO) Update(ctx context.Context, row *ent.IntentRule) error {
	builder := d.client.IntentRule.Update().
		Where(intentrule.RuleIDEQ(row.RuleID)).
		SetIntentName(row.IntentName).
		SetPatterns(row.Patterns).
		SetMatchType(row.MatchType).
		SetPriority(row.Priority).
		SetConfidence(row.Confidence)

	if row.Params != nil {
		builder.SetParams(row.Params)
	} else {
		builder.ClearParams()
	}

	n, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("rule %s not found", row.RuleID)
	}
	return nil
}

func (d *IntentRuleDAO) Delete(ctx context.Context, ruleID string) error {
	n, err := d.client.IntentRule.Update().
		Where(intentrule.RuleIDEQ(ruleID)).
		SetStatus(0).
		Save(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("rule %s not found", ruleID)
	}
	return nil
}
