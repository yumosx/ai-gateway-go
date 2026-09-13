package llmlayer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

// IntentDefinition 意图定义（用于构建 few-shot prompt）
type IntentDefinition struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Examples []string `json:"examples"`
}

// Engine LLM 意图识别引擎接口
type Engine interface {
	Recognize(ctx context.Context, text string, intents []IntentDefinition) (domain.Intent, error)
}

// LLMRecognizer 大模型意图识别器
type LLMRecognizer struct {
	client    *deepseek.Client
	model     string
	timeoutMs int
}

// NewLLMRecognizer 创建大模型意图识别器
func NewLLMRecognizer(client *deepseek.Client, model string, timeoutMs int) *LLMRecognizer {
	if model == "" {
		model = deepseek.DeepSeekChat
	}
	if timeoutMs <= 0 {
		timeoutMs = 3000
	}
	return &LLMRecognizer{
		client:    client,
		model:     model,
		timeoutMs: timeoutMs,
	}
}

// llmOutput 大模型结构化输出
type llmOutput struct {
	Intent   string            `json:"intent"`
	Confidence float64         `json:"confidence"`
	Params   map[string]string `json:"params,omitempty"`
}

// Recognize 调用大模型识别意图
func (r *LLMRecognizer) Recognize(ctx context.Context, text string, intents []IntentDefinition) (domain.Intent, error) {
	if len(intents) == 0 {
		return domain.Intent{}, fmt.Errorf("no intent definitions provided")
	}

	// 构建 few-shot prompt
	prompt := r.buildPrompt(text, intents)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	request := &deepseek.ChatCompletionRequest{
		Model: r.model,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleSystem,
				Content: "你是一个意图识别助手。根据用户输入，从给定的意图列表中选择最匹配的意图。请严格以JSON格式返回。",
			},
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &deepseek.ResponseFormat{Type: "json_object"},
	}

	response, err := r.client.CreateChatCompletion(timeoutCtx, request)
	if err != nil {
		return domain.Intent{}, fmt.Errorf("llm call failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return domain.Intent{}, fmt.Errorf("llm returned no choices")
	}

	content := response.Choices[0].Message.Content
	return r.parseOutput(content)
}

// buildPrompt 构建 few-shot prompt
func (r *LLMRecognizer) buildPrompt(text string, intents []IntentDefinition) string {
	var sb strings.Builder

	sb.WriteString("## 可选意图列表\n\n")
	for _, intent := range intents {
		sb.WriteString(fmt.Sprintf("### %s\n", intent.Name))
		sb.WriteString(fmt.Sprintf("描述: %s\n", intent.Desc))
		if len(intent.Examples) > 0 {
			sb.WriteString("示例:\n")
			for _, ex := range intent.Examples {
				sb.WriteString(fmt.Sprintf("- \"%s\"\n", ex))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 用户输入\n\n")
	sb.WriteString(fmt.Sprintf("\"%s\"\n\n", text))

	sb.WriteString("## 输出要求\n\n")
	sb.WriteString("请返回JSON格式:\n")
	sb.WriteString("```json\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"intent\": \"意图名称\",\n")
	sb.WriteString("  \"confidence\": 0.0到1.0之间的置信度,\n")
	sb.WriteString("  \"params\": {\"参数名\": \"参数值\"}\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("如果没有匹配的意图，intent设为\"unknown\"，confidence设为0。\n")

	return sb.String()
}

// parseOutput 解析大模型输出
func (r *LLMRecognizer) parseOutput(content string) (domain.Intent, error) {
	// 尝试提取 JSON 部分
	jsonStr := extractJSON(content)

	var out llmOutput
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		return domain.Intent{}, fmt.Errorf("failed to parse llm output: %w", err)
	}

	// 限制 confidence 范围
	if out.Confidence < 0 {
		out.Confidence = 0
	}
	if out.Confidence > 1 {
		out.Confidence = 1
	}

	return domain.Intent{
		Name:       out.Intent,
		Confidence: out.Confidence,
		Params:     out.Params,
	}, nil
}

// extractJSON 从可能包含 markdown 的文本中提取 JSON
func extractJSON(s string) string {
	// 尝试找 ```json ... ```
	if idx := strings.Index(s, "```json"); idx != -1 {
		start := idx + len("```json")
		if end := strings.Index(s[start:], "```"); end != -1 {
			return strings.TrimSpace(s[start : start+end])
		}
	}
	// 尝试找 { ... }
	if start := strings.Index(s, "{"); start != -1 {
		if end := strings.LastIndex(s, "}"); end != -1 {
			return s[start : end+1]
		}
	}
	return s
}
