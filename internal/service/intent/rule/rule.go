package rule

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

// Engine 规则引擎接口
type Engine interface {
	// Match 对输入文本进行规则匹配，返回命中的意图列表（按 priority 降序）
	Match(ctx context.Context, text string) []domain.Intent
	// Reload 热加载规则
	Reload(rules []domain.IntentRule)
}

// Matcher 规则匹配器实现
type Matcher struct {
	mu    sync.RWMutex
	rules []domain.IntentRule
	// 预编译的正则缓存
	regexCache map[string]*regexp.Regexp
}

// NewMatcher 创建规则匹配器
func NewMatcher(rules []domain.IntentRule) *Matcher {
	m := &Matcher{
		rules:      rules,
		regexCache: make(map[string]*regexp.Regexp),
	}
	m.buildRegexCache()
	return m
}
// RuleSource 规则数据源（由 repo 实现）
type RuleSource interface {
	GetAll(ctx context.Context) ([]domain.IntentRule, error)
}

// NewMatcherFromRepo 从数据源加载规则创建匹配器
func NewMatcherFromRepo(src RuleSource) (*Matcher, error) {
	rules, err := src.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("load rules from repo: %w", err)
	}
	return NewMatcher(rules), nil
}

// ReloadFromRepo 从数据源重新加载规则到内存
func (m *Matcher) ReloadFromRepo(src RuleSource) error {
	rules, err := src.GetAll(context.Background())
	if err != nil {
		return fmt.Errorf("reload rules from repo: %w", err)
	}
	m.Reload(rules)
	return nil
}
func (m *Matcher) Match(ctx context.Context, text string) []domain.Intent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lowerText := strings.ToLower(text)
	type scored struct {
		intent   domain.Intent
		priority int
	}
	var hits []scored

	for _, rule := range m.rules {
		if ctx.Err() != nil {
			break
		}
		for _, pattern := range rule.Patterns {
			if m.matchPattern(lowerText, pattern, rule.MatchType) {
				hits = append(hits, scored{
					intent: domain.Intent{
						Name:       rule.IntentName,
						Confidence: rule.Confidence,
						Params:     cloneMap(rule.Params),
					},
					priority: rule.Priority,
				})
				break // 一个规则只命中一次
			}
		}
	}

	// 按 priority 降序排列
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].priority > hits[j].priority
	})

	intents := make([]domain.Intent, 0, len(hits))
	for _, h := range hits {
		intents = append(intents, h.intent)
	}
	return intents
}

func (m *Matcher) Reload(rules []domain.IntentRule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = rules
	m.regexCache = make(map[string]*regexp.Regexp)
	m.buildRegexCache()
}

func (m *Matcher) matchPattern(text, pattern string, matchType domain.RuleMatchType) bool {
	pattern = strings.ToLower(pattern)
	switch matchType {
	case domain.MatchExact:
		return text == pattern
	case domain.MatchPrefix:
		return strings.HasPrefix(text, pattern)
	case domain.MatchContains:
		return strings.Contains(text, pattern)
	case domain.MatchRegex:
		return m.matchRegex(text, pattern)
	default:
		return strings.Contains(text, pattern)
	}
}

func (m *Matcher) matchRegex(text, pattern string) bool {
	re, ok := m.regexCache[pattern]
	if !ok {
		return false
	}
	return re.MatchString(text)
}

func (m *Matcher) buildRegexCache() {
	for _, rule := range m.rules {
		if rule.MatchType == domain.MatchRegex {
			for _, p := range rule.Patterns {
				if _, exists := m.regexCache[p]; !exists {
					if re, err := regexp.Compile(p); err == nil {
						m.regexCache[p] = re
					}
				}
			}
		}
	}
}

func cloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

