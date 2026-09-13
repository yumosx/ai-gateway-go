package domain

// Clarification 反问澄清
type Clarification struct {
	Question string
	Options  []string
	Hint     string
}

// ClarifyRequest 用户回答澄清问题
type ClarifyRequest struct {
	SessionID string
	Answer    string
}

// AutoCompleteRule 参数自动补全规则
type AutoCompleteRule struct {
	IntentName string
	ParamName  string
	Source     string // context / default / history
	Default    string
	Prompt     string // 补全失败时的提示
}
