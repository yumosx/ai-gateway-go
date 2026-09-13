package rule

import (
	"context"
	"testing"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

func TestMatcher_Exact(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "greeting", Patterns: []string{"你好"}, MatchType: domain.MatchExact, Priority: 10, Confidence: 1.0},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "你好")
	if len(hits) != 1 || hits[0].Name != "greeting" {
		t.Fatalf("expected greeting, got %v", hits)
	}
	// 精确匹配不应命中 "你好啊"
	hits = m.Match(context.Background(), "你好啊")
	if len(hits) != 0 {
		t.Fatalf("expected no match for '你好啊', got %v", hits)
	}
}

func TestMatcher_Contains(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "weather", Patterns: []string{"天气"}, MatchType: domain.MatchContains, Priority: 10, Confidence: 0.98},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "今天天气怎么样")
	if len(hits) != 1 || hits[0].Name != "weather" {
		t.Fatalf("expected weather, got %v", hits)
	}
}

func TestMatcher_Prefix(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "cmd", Patterns: []string{"/help"}, MatchType: domain.MatchPrefix, Priority: 10, Confidence: 1.0},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "/help me")
	if len(hits) != 1 || hits[0].Name != "cmd" {
		t.Fatalf("expected cmd, got %v", hits)
	}
	hits = m.Match(context.Background(), "show /help")
	if len(hits) != 0 {
		t.Fatalf("expected no match, got %v", hits)
	}
}

func TestMatcher_Regex(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "date", Patterns: []string{`\d{4}-\d{2}-\d{2}`}, MatchType: domain.MatchRegex, Priority: 10, Confidence: 0.9},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "2024-01-15 的日程")
	if len(hits) != 1 || hits[0].Name != "date" {
		t.Fatalf("expected date, got %v", hits)
	}
}

func TestMatcher_Priority(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "low", Patterns: []string{"test"}, MatchType: domain.MatchContains, Priority: 1, Confidence: 0.5},
		{ID: "2", IntentName: "high", Patterns: []string{"test"}, MatchType: domain.MatchContains, Priority: 100, Confidence: 0.9},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "test input")
	if len(hits) != 2 || hits[0].Name != "high" {
		t.Fatalf("expected high first, got %v", hits)
	}
}

func TestMatcher_Params(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "weather", Patterns: []string{"天气"}, MatchType: domain.MatchContains, Priority: 10, Confidence: 0.98, Params: map[string]string{"date": "today"}},
	}
	m := NewMatcher(rules)
	hits := m.Match(context.Background(), "今天天气怎么样")
	if len(hits) != 1 || hits[0].Params["date"] != "today" {
		t.Fatalf("expected params with date=today, got %v", hits)
	}
}

func TestMatcher_ContextCancel(t *testing.T) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "greeting", Patterns: []string{"你好"}, MatchType: domain.MatchExact, Priority: 10, Confidence: 1.0},
	}
	m := NewMatcher(rules)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消
	hits := m.Match(ctx, "你好")
	if len(hits) != 0 {
		t.Fatalf("expected empty on cancelled context, got %v", hits)
	}
}

func BenchmarkMatcher_Contains(b *testing.B) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "weather", Patterns: []string{"天气", "下雨", "温度"}, MatchType: domain.MatchContains, Priority: 10, Confidence: 0.98},
		{ID: "2", IntentName: "greeting", Patterns: []string{"你好", "hello", "hi"}, MatchType: domain.MatchContains, Priority: 10, Confidence: 1.0},
		{ID: "3", IntentName: "code", Patterns: []string{"代码", "编程", "算法"}, MatchType: domain.MatchContains, Priority: 10, Confidence: 0.95},
	}
	m := NewMatcher(rules)
	ctx := context.Background()
	input := "今天天气怎么样，顺便帮我写个代码"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match(ctx, input)
	}
}

func BenchmarkMatcher_Regex(b *testing.B) {
	rules := []domain.IntentRule{
		{ID: "1", IntentName: "date", Patterns: []string{`\d{4}-\d{2}-\d{2}`}, MatchType: domain.MatchRegex, Priority: 10, Confidence: 0.9},
	}
	m := NewMatcher(rules)
	ctx := context.Background()
	input := "2024-01-15 的日程安排"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match(ctx, input)
	}
}
