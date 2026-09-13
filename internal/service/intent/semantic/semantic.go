package semantic

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
)

// Engine 语义检索引擎接口
type Engine interface {
	// Retrieve 根据输入文本检索最相似的意图候选
	Retrieve(ctx context.Context, text string, topK int) ([]Candidate, error)
	// AddIntent 添加意图及其示例到索引
	AddIntent(name string, examples []string) error
	// BuildIndex 构建索引（预计算向量）
	BuildIndex() error
}

// Candidate 检索候选结果
type Candidate struct {
	IntentName string  `json:"intent_name"`
	Score      float64 `json:"score"`
	Example    string  `json:"example,omitempty"`
}

// IntentEntry 意意图谱条目
type IntentEntry struct {
	Name     string
	Examples []string
}

// Retriever 向量检索器
type Retriever struct {
	entries    []IntentEntry
	embeddings [][]float64   // 每条 example 对应的 embedding
	exampleIdx []int         // embeddings[i] 属于 entries[exampleIdx[i]]
	built      bool
}

// NewRetriever 创建检索器
func NewRetriever() *Retriever {
	return &Retriever{}
}

// AddIntent 添加意图
func (r *Retriever) AddIntent(name string, examples []string) error {
	r.entries = append(r.entries, IntentEntry{
		Name:     name,
		Examples: examples,
	})
	r.built = false
	return nil
}

// BuildIndex 构建向量索引
func (r *Retriever) BuildIndex() error {
	r.embeddings = r.embeddings[:0]
	r.exampleIdx = r.exampleIdx[:0]

	for i, entry := range r.entries {
		for _, ex := range entry.Examples {
			vec := simpleEmbed(ex)
			r.embeddings = append(r.embeddings, vec)
			r.exampleIdx = append(r.exampleIdx, i)
		}
	}
	r.built = true
	return nil
}

// Retrieve 检索 top-K 候选
func (r *Retriever) Retrieve(ctx context.Context, text string, topK int) ([]Candidate, error) {
	if !r.built {
		if err := r.BuildIndex(); err != nil {
			return nil, fmt.Errorf("build index failed: %w", err)
		}
	}

	if len(r.embeddings) == 0 {
		return nil, nil
	}

	queryVec := simpleEmbed(text)

	type scored struct {
		idx   int
		score float64
		examp string
	}

	results := make([]scored, 0, len(r.embeddings))
	for i, vec := range r.embeddings {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		sim := cosineSimilarity(queryVec, vec)
		results = append(results, scored{
			idx:   i,
			score: sim,
			examp: r.entries[r.exampleIdx[i]].Examples[0], // 取第一个示例
		})
	}

	// 按 score 降序排列
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if topK <= 0 || topK > len(results) {
		topK = len(results)
	}

	candidates := make([]Candidate, 0, topK)
	for i := 0; i < topK; i++ {
		candidates = append(candidates, Candidate{
			IntentName: r.entries[r.exampleIdx[results[i].idx]].Name,
			Score:      results[i].score,
			Example:    results[i].examp,
		})
	}

	return candidates, nil
}

// simpleEmbed 简单的文本向量化（基于字符 n-gram + 哈希）
// 生产环境应替换为真实的 embedding 模型
func simpleEmbed(text string) []float64 {
	lower := strings.ToLower(text)
	dim := 128
	vec := make([]float64, dim)

	// unigram
	for i := 0; i < len(lower); i++ {
		idx := int(lower[i]) % dim
		vec[idx] += 1.0
	}

	// bigram
	for i := 0; i+1 < len(lower); i++ {
		idx := int(lower[i])*31 + int(lower[i+1])
		idx = ((idx % dim) + dim) % dim
		vec[idx] += 1.0
	}

	// trigram
	for i := 0; i+2 < len(lower); i++ {
		idx := int(lower[i])*31*31 + int(lower[i+1])*31 + int(lower[i+2])
		idx = ((idx % dim) + dim) % dim
		vec[idx] += 0.5
	}

	// 归一化
	var norm float64
	for _, v := range vec {
		norm += v * v
	}
	if norm > 0 {
		norm = math.Sqrt(norm)
		for i := range vec {
			vec[i] /= norm
		}
	}

	return vec
}

// cosineSimilarity 计算余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
