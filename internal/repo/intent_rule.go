package repo

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ai-gateway-go/ent"
	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	"github.com/ecodeclub/ai-gateway-go/internal/repo/cache"
	"github.com/ecodeclub/ai-gateway-go/internal/repo/dao"
)

// IntentRuleRepo 聚合 cache + dao，对外接受/返回 domain
type IntentRuleRepo struct {
	dao   *dao.IntentRuleDAO
	cache *cache.IntentRuleCache
}

func NewIntentRuleRepo(d *dao.IntentRuleDAO, c *cache.IntentRuleCache) *IntentRuleRepo {
	return &IntentRuleRepo{dao: d, cache: c}
}

func (r *IntentRuleRepo) GetAll(ctx context.Context) ([]domain.IntentRule, error) {
	if rows, err := r.cache.GetAll(ctx); err == nil && len(rows) > 0 {
		return toDomainRules(rows), nil
	}

	rows, err := r.dao.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAll: %w", err)
	}
	_ = r.cache.SetAll(ctx, rows)
	return toDomainRules(rows), nil
}

func (r *IntentRuleRepo) GetByIntentName(ctx context.Context, name string) ([]domain.IntentRule, error) {
	all, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []domain.IntentRule
	for _, rule := range all {
		if rule.IntentName == name {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *IntentRuleRepo) Create(ctx context.Context, rule domain.IntentRule) error {
	if err := r.dao.Create(ctx, toEntRule(rule)); err != nil {
		return err
	}
	r.cache.Invalidate(ctx)
	return nil
}

func (r *IntentRuleRepo) Update(ctx context.Context, rule domain.IntentRule) error {
	if err := r.dao.Update(ctx, toEntRule(rule)); err != nil {
		return err
	}
	r.cache.Invalidate(ctx)
	return nil
}

func (r *IntentRuleRepo) Delete(ctx context.Context, ruleID string) error {
	if err := r.dao.Delete(ctx, ruleID); err != nil {
		return err
	}
	r.cache.Invalidate(ctx)
	return nil
}

func (r *IntentRuleRepo) SubscribeReload(ctx context.Context, onReload func()) {
	r.cache.SubscribeReload(ctx, onReload)
}

func toDomainRules(rows []*ent.IntentRule) []domain.IntentRule {
	rules := make([]domain.IntentRule, 0, len(rows))
	for _, row := range rows {
		rules = append(rules, toDomainRule(row))
	}
	return rules
}

func toDomainRule(row *ent.IntentRule) domain.IntentRule {
	return domain.IntentRule{
		ID:         row.RuleID,
		IntentName: row.IntentName,
		Patterns:   row.Patterns,
		MatchType:  domain.RuleMatchType(row.MatchType),
		Params:     row.Params,
		Priority:   row.Priority,
		Confidence: row.Confidence,
	}
}

func toEntRule(rule domain.IntentRule) *ent.IntentRule {
	return &ent.IntentRule{
		RuleID:     rule.ID,
		IntentName: rule.IntentName,
		Patterns:   rule.Patterns,
		MatchType:  string(rule.MatchType),
		Params:     rule.Params,
		Priority:   rule.Priority,
		Confidence: rule.Confidence,
		Status:     1,
	}
}
