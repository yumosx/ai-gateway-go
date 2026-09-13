package intent

import (
	"context"
	"time"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/llmlayer"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/rule"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/semantic"
)

const (
	// HighConfidenceThreshold 高置信度阈值，命中直接返回
	HighConfidenceThreshold = 0.95
	// LLMTimeoutMs LLM 层超时
	LLMTimeoutMs = 3000
	// SemanticTimeoutMs 语义检索层超时
	SemanticTimeoutMs = 1000
)

// Pipeline 分层意图识别编排器
type Pipeline struct {
	ruleEngine     rule.Engine
	llmEngine      llmlayer.Engine
	semanticEngine semantic.Engine
	intents        []llmlayer.IntentDefinition // 意图定义列表，供 LLM few-shot 使用
}

// NewPipeline 创建 pipeline
func NewPipeline(
	ruleEngine rule.Engine,
	llmEngine llmlayer.Engine,
	semanticEngine semantic.Engine,
	intents []llmlayer.IntentDefinition,
) *Pipeline {
	return &Pipeline{
		ruleEngine:     ruleEngine,
		llmEngine:      llmEngine,
		semanticEngine: semanticEngine,
		intents:        intents,
	}
}

// Recognize 分层意图识别
func (p *Pipeline) Recognize(ctx context.Context, req domain.IntentRequest) domain.IntentResponse {
	start := time.Now()
	resp := domain.IntentResponse{
		Intents: make([]domain.Intent, 0),
	}

	// Layer 1: Rule engine — 快速命中
	if hits := p.ruleEngine.Match(ctx, req.Text); len(hits) > 0 {
		resp.Intents = hits
		resp.MatchedLayer = string(domain.LayerRule)
		resp.TotalLatencyMs = time.Since(start).Milliseconds()

		// 高置信度直接返回
		if hits[0].Confidence >= HighConfidenceThreshold {
			return resp
		}
	}

	// Layer 2: LLM dynamic-example — 超时降级
	llmCtx, llmCancel := context.WithTimeout(ctx, time.Duration(LLMTimeoutMs)*time.Millisecond)
	defer llmCancel()

	llmResult, err := p.llmEngine.Recognize(llmCtx, req.Text, p.intents)
	if err == nil && llmResult.Name != "unknown" && llmResult.Confidence > 0 {
		// 合并：如果 rule 层也有结果，合并到列表
		if len(resp.Intents) > 0 {
			resp.Intents = mergeIntents(resp.Intents, llmResult)
		} else {
			resp.Intents = []domain.Intent{llmResult}
		}
		resp.MatchedLayer = string(domain.LayerLLM)
		resp.TotalLatencyMs = time.Since(start).Milliseconds()

		if llmResult.Confidence >= HighConfidenceThreshold {
			return resp
		}
	}

	// Layer 3: Semantic retrieval — 最终兜底
	semCtx, semCancel := context.WithTimeout(ctx, time.Duration(SemanticTimeoutMs)*time.Millisecond)
	defer semCancel()

	candidates, err := p.semanticEngine.Retrieve(semCtx, req.Text, 3)
	if err == nil && len(candidates) > 0 {
		semIntent := domain.Intent{
			Name:       candidates[0].IntentName,
			Confidence: candidates[0].Score,
			Params:     make(map[string]string),
		}
		if len(resp.Intents) > 0 {
			resp.Intents = mergeIntents(resp.Intents, semIntent)
		} else {
			resp.Intents = []domain.Intent{semIntent}
		}
		resp.MatchedLayer = string(domain.LayerSemantic)
	}

	resp.TotalLatencyMs = time.Since(start).Milliseconds()
	return resp
}

// mergeIntents 合并意图列表，去重取最高 confidence
func mergeIntents(existing []domain.Intent, new domain.Intent) []domain.Intent {
	for i, e := range existing {
		if e.Name == new.Name {
			if new.Confidence > e.Confidence {
				existing[i].Confidence = new.Confidence
				if new.Params != nil {
					if existing[i].Params == nil {
						existing[i].Params = make(map[string]string)
					}
					for k, v := range new.Params {
						existing[i].Params[k] = v
					}
				}
			}
			return existing
		}
	}
	return append(existing, new)
}
