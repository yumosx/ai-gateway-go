package intent

import (
	"fmt"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

const (
	// ClarifyLowThreshold 低于此阈值不做任何处理
	ClarifyLowThreshold = 0.3
	// ClarifyHighThreshold 高于此阈值不需要澄清
	ClarifyHighThreshold = 0.8
)

// Clarifier 反问澄清器
type Clarifier struct {
	autoCompleteRules []domain.AutoCompleteRule
	// intentQuestions 意图对应的澄清问题模板
	intentQuestions map[string]string
}

// NewClarifier 创建澄清器
func NewClarifier(rules []domain.AutoCompleteRule, questions map[string]string) *Clarifier {
	return &Clarifier{
		autoCompleteRules: rules,
		intentQuestions:   questions,
	}
}

// Evaluate 评估是否需要反问澄清，并补全参数
func (c *Clarifier) Evaluate(resp *domain.IntentResponse, ctx map[string]string) {
	if len(resp.Intents) == 0 {
		return
	}

	best := resp.Intents[0]

	// 参数自动补全
	completed := c.autoComplete(best.Name, best.Params, ctx)
	if len(completed) > 0 {
		if resp.AutoCompleted == nil {
			resp.AutoCompleted = make(map[string]string)
		}
		for k, v := range completed {
			if best.Params == nil {
				best.Params = make(map[string]string)
			}
			best.Params[k] = v
			resp.AutoCompleted[k] = v
		}
		resp.Intents[0] = best
	}

	// 判断是否需要澄清
	if best.Confidence >= ClarifyHighThreshold || best.Confidence <= ClarifyLowThreshold {
		resp.NeedClarify = false
		return
	}

	resp.NeedClarify = true
	resp.Clarification = c.buildClarification(best.Name, best.Params)
}

// HandleClarify 处理用户的澄清回答
func (c *Clarifier) HandleClarify(sessionID, answer string, resp *domain.IntentResponse) {
	if resp == nil || len(resp.Intents) == 0 {
		return
	}

	// 简单策略：用户回答作为额外上下文合并到 params
	best := resp.Intents[0]
	if best.Params == nil {
		best.Params = make(map[string]string)
	}
	best.Params["user_clarification"] = answer
	// 置信度提升
	best.Confidence = min(best.Confidence+0.15, 1.0)
	resp.Intents[0] = best
	resp.NeedClarify = false
	resp.Clarification = nil
}

// autoComplete 自动补全参数
func (c *Clarifier) autoComplete(intentName string, currentParams, ctx map[string]string) map[string]string {
	completed := make(map[string]string)

	for _, rule := range c.autoCompleteRules {
		if rule.IntentName != intentName {
			continue
		}

		// 如果参数已有值，跳过
		if currentParams != nil {
			if _, exists := currentParams[rule.ParamName]; exists {
				continue
			}
		}

		switch rule.Source {
		case "context":
			if ctx != nil {
				if val, ok := ctx[rule.ParamName]; ok {
					completed[rule.ParamName] = val
					continue
				}
			}
			// context 中没有，尝试 default
			if rule.Default != "" {
				completed[rule.ParamName] = rule.Default
			}
		case "default":
			if rule.Default != "" {
				completed[rule.ParamName] = rule.Default
			}
		}
	}

	return completed
}

// buildClarification 构建澄清问题
func (c *Clarifier) buildClarification(intentName string, params map[string]string) *domain.Clarification {
	question, ok := c.intentQuestions[intentName]
	if !ok {
		question = fmt.Sprintf("请确认您是要执行「%s」操作吗？", intentName)
	}

	return &domain.Clarification{
		Question: question,
		Hint:     "请提供更详细的信息以完成识别。",
	}
}
