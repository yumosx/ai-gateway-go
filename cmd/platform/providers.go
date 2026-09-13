package main

import (
	ds "github.com/cohesion-org/deepseek-go"
	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	"github.com/ecodeclub/ai-gateway-go/internal/repo"
	intent "github.com/ecodeclub/ai-gateway-go/internal/service/intent"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/llmlayer"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/rule"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/semantic"
	"github.com/ecodeclub/ai-gateway-go/internal/service/llm"
	"github.com/ecodeclub/ai-gateway-go/internal/service/llm/platform/deepseek"
)

func NewRuleMatcher(ruleRepo *repo.IntentRuleRepo) (*rule.Matcher, error) {
	return rule.NewMatcherFromRepo(ruleRepo)
}

func AsRuleEngine(m *rule.Matcher) rule.Engine { return m }

func NewLLMRecognizer(client *ds.Client) *llmlayer.LLMRecognizer {
	return llmlayer.NewLLMRecognizer(client, "", 3000)
}

func AsLLMEngine(r *llmlayer.LLMRecognizer) llmlayer.Engine { return r }

func NewDeepSeekHandler(client *ds.Client) *deepseek.Handler {
	return deepseek.NewHandler(client)
}

func AsLLMHandler(h *deepseek.Handler) llm.LLMHandler { return h }

func NewSemanticEngine() *semantic.Retriever {
	engine := semantic.NewRetriever()
	_ = engine.AddIntent("greeting", []string{"你好", "hello", "hi", "早上好"})
	_ = engine.AddIntent("weather_query", []string{"今天天气怎么样", "明天会下雨吗", "天气预报"})
	_ = engine.AddIntent("code_help", []string{"帮我写一个排序算法", "实现一个HTTP服务器", "写代码"})
	_ = engine.AddIntent("translation", []string{"翻译成英文", "translate to Chinese", "翻译这段话"})
	_ = engine.AddIntent("knowledge_qa", []string{"什么是量子计算", "解释一下机器学习", "区块链是什么"})
	_ = engine.BuildIndex()
	return engine
}

func AsSemanticEngine(r *semantic.Retriever) semantic.Engine { return r }

func NewIntentDefinitions() []llmlayer.IntentDefinition {
	return []llmlayer.IntentDefinition{
		{Name: "greeting", Desc: "问候语", Examples: []string{"你好", "hello", "早上好"}},
		{Name: "weather_query", Desc: "查询天气", Examples: []string{"今天天气怎么样", "明天会下雨吗"}},
		{Name: "code_help", Desc: "代码相关帮助", Examples: []string{"帮我写一个排序算法", "实现一个HTTP服务器"}},
		{Name: "translation", Desc: "翻译请求", Examples: []string{"翻译成英文", "translate to Chinese"}},
		{Name: "knowledge_qa", Desc: "知识问答", Examples: []string{"什么是量子计算", "解释一下机器学习"}},
	}
}

func NewClarifier() *intent.Clarifier {
	autoCompleteRules := []domain.AutoCompleteRule{
		{IntentName: "weather_query", ParamName: "date", Source: "default", Default: "today"},
		{IntentName: "weather_query", ParamName: "city", Source: "context"},
		{IntentName: "translation", ParamName: "target_lang", Source: "context", Default: "english"},
	}
	questions := map[string]string{
		"weather_query": "请问您想查询哪个城市、哪一天的天气？",
		"translation":   "请问您想翻译成什么语言？",
		"code_help":     "请问您需要什么编程语言的代码帮助？",
		"knowledge_qa":  "请问您想了解哪个方面的知识？",
	}
	return intent.NewClarifier(autoCompleteRules, questions)
}
