package semantic

import (
	"context"
	"testing"
)

func TestRetriever_Basic(t *testing.T) {
	r := NewRetriever()
	_ = r.AddIntent("greeting", []string{"你好", "hello", "早上好"})
	_ = r.AddIntent("weather", []string{"今天天气怎么样", "明天会下雨吗"})
	_ = r.AddIntent("code", []string{"帮我写代码", "实现一个算法"})
	_ = r.BuildIndex()

	candidates, err := r.Retrieve(context.Background(), "你好世界", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}
	// "你好" 应该和 greeting 最相似
	if candidates[0].IntentName != "greeting" {
		t.Fatalf("expected greeting as top candidate, got %s (score=%.3f)", candidates[0].IntentName, candidates[0].Score)
	}
}

func TestRetriever_TopK(t *testing.T) {
	r := NewRetriever()
	_ = r.AddIntent("a", []string{"apple"})
	_ = r.AddIntent("b", []string{"banana"})
	_ = r.AddIntent("c", []string{"cherry"})
	_ = r.BuildIndex()

	candidates, err := r.Retrieve(context.Background(), "fruit", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) > 2 {
		t.Fatalf("expected at most 2 candidates, got %d", len(candidates))
	}
}

func TestRetriever_Empty(t *testing.T) {
	r := NewRetriever()
	_ = r.BuildIndex()
	candidates, err := r.Retrieve(context.Background(), "test", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 0 {
		t.Fatalf("expected 0 candidates, got %d", len(candidates))
	}
}

func TestCosineSimilarity(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{1, 0, 0}
	sim := cosineSimilarity(a, b)
	if sim < 0.999 {
		t.Fatalf("expected ~1.0 for identical vectors, got %f", sim)
	}

	c := []float64{0, 1, 0}
	sim = cosineSimilarity(a, c)
	if sim > 0.001 {
		t.Fatalf("expected ~0.0 for orthogonal vectors, got %f", sim)
	}
}

func BenchmarkRetriever_Retrieve(b *testing.B) {
	r := NewRetriever()
	intents := []struct {
		name     string
		examples []string
	}{
		{"greeting", []string{"你好", "hello", "hi", "嗨", "早上好", "下午好"}},
		{"weather", []string{"今天天气怎么样", "明天会下雨吗", "天气预报", "温度多少"}},
		{"code", []string{"帮我写代码", "实现一个算法", "代码调试", "编程问题"}},
		{"translation", []string{"翻译成英文", "translate to Chinese", "翻译这段话"}},
		{"qa", []string{"什么是机器学习", "解释量子计算", "区块链原理", "人工智能定义"}},
	}
	for _, it := range intents {
		_ = r.AddIntent(it.name, it.examples)
	}
	_ = r.BuildIndex()

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Retrieve(ctx, "请帮我翻译一下这段话", 3)
	}
}
