package domain

// IntentRequest 用户意图识别请求
type IntentRequest struct {
	Text      string
	Context   map[string]string // 上下文信息，用于参数补全
	SessionID string
}

// Intent 单个识别到的意图
type Intent struct {
	Name       string
	Confidence float64
	Params     map[string]string // 自动提取/补全的参数
}

// IntentResponse 意图识别响应
type IntentResponse struct {
	Intents        []Intent
	MatchedLayer   string // rule / llm / semantic
	NeedClarify    bool
	Clarification  *Clarification
	AutoCompleted  map[string]string // 自动补全的参数
	TotalLatencyMs int64
}

// --- 规则引擎相关 ---

// RuleMatchType 规则匹配类型
type RuleMatchType string

const (
	MatchExact    RuleMatchType = "exact"    // 精确匹配
	MatchPrefix   RuleMatchType = "prefix"   // 前缀匹配
	MatchRegex    RuleMatchType = "regex"    // 正则匹配
	MatchContains RuleMatchType = "contains" // 包含匹配
)

// IntentRule 单条意图识别规则
type IntentRule struct {
	ID         string
	IntentName string
	Patterns   []string // 匹配模式列表
	MatchType  RuleMatchType
	Params     map[string]string // 规则携带的默认参数
	Priority   int               // 优先级，数值越大越优先
	Confidence float64           // 该规则的固定置信度
}

// --- 语义检索相关 ---

// IntentEmbedding 意图向量表示
type IntentEmbedding struct {
	IntentName string
	Embedding  []float64
	Examples   []string // few-shot 示例
}

// SemanticCandidate 语义检索候选结果
type SemanticCandidate struct {
	IntentName string
	Score      float64 // cosine similarity
	Example    string  // 最近邻示例
}

// --- Pipeline 相关 ---

// RecognizeLayer 识别层
type RecognizeLayer string

const (
	LayerRule     RecognizeLayer = "rule"
	LayerLLM      RecognizeLayer = "llm"
	LayerSemantic RecognizeLayer = "semantic"
)
