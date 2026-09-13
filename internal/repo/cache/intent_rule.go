package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ecodeclub/ai-gateway-go/ent"
	"github.com/redis/go-redis/v9"
)

const (
	rulesCacheKey   = "intent:rules:all"
	rulesVersionKey = "intent:rules:version"
	rulesCacheTTL   = 5 * time.Minute
)

// intentRuleItem Redis 中存储的 schema 字段快照
type intentRuleItem struct {
	RuleID     string            `json:"rule_id"`
	IntentName string            `json:"intent_name"`
	Patterns   []string          `json:"patterns"`
	MatchType  string            `json:"match_type"`
	Params     map[string]string `json:"params,omitempty"`
	Priority   int               `json:"priority"`
	Confidence float64           `json:"confidence"`
	Status     int8              `json:"status"`
}

// IntentRuleCache Redis 规则缓存（接受/返回 schema）
type IntentRuleCache struct {
	rdb *redis.Client
}

func NewIntentRuleCache(rdb *redis.Client) *IntentRuleCache {
	return &IntentRuleCache{rdb: rdb}
}

func (c *IntentRuleCache) GetAll(ctx context.Context) ([]*ent.IntentRule, error) {
	data, err := c.rdb.HGetAll(ctx, rulesCacheKey).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("redis cache empty")
	}

	rows := make([]*ent.IntentRule, 0, len(data))
	for _, jsonStr := range data {
		var item intentRuleItem
		if err := json.Unmarshal([]byte(jsonStr), &item); err != nil {
			continue
		}
		rows = append(rows, itemToEnt(item))
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("redis cache empty")
	}
	return rows, nil
}

func (c *IntentRuleCache) SetAll(ctx context.Context, rows []*ent.IntentRule) error {
	if len(rows) == 0 {
		return nil
	}

	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, rulesCacheKey)

	for _, row := range rows {
		b, err := json.Marshal(entToItem(row))
		if err != nil {
			continue
		}
		pipe.HSet(ctx, rulesCacheKey, row.RuleID, string(b))
	}
	pipe.Expire(ctx, rulesCacheKey, rulesCacheTTL)

	_, err := pipe.Exec(ctx)
	return err
}

func (c *IntentRuleCache) Invalidate(ctx context.Context) {
	c.rdb.Del(ctx, rulesCacheKey)
	_ = c.rdb.Publish(ctx, "intent:rules:reload", "update").Err()
}

func (c *IntentRuleCache) PublishReload(ctx context.Context) error {
	return c.rdb.Incr(ctx, rulesVersionKey).Err()
}

func (c *IntentRuleCache) SubscribeReload(ctx context.Context, onReload func()) {
	go func() {
		sub := c.rdb.Subscribe(ctx, "intent:rules:reload")
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				onReload()
			}
		}
	}()
}

func entToItem(row *ent.IntentRule) intentRuleItem {
	return intentRuleItem{
		RuleID:     row.RuleID,
		IntentName: row.IntentName,
		Patterns:   row.Patterns,
		MatchType:  row.MatchType,
		Params:     row.Params,
		Priority:   row.Priority,
		Confidence: row.Confidence,
		Status:     row.Status,
	}
}

func itemToEnt(item intentRuleItem) *ent.IntentRule {
	return &ent.IntentRule{
		RuleID:     item.RuleID,
		IntentName: item.IntentName,
		Patterns:   item.Patterns,
		MatchType:  item.MatchType,
		Params:     item.Params,
		Priority:   item.Priority,
		Confidence: item.Confidence,
		Status:     item.Status,
	}
}
