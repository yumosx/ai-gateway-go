package schema

type CreateRuleRequest struct {
	IntentName string            `json:"intent_name" binding:"required"`
	Patterns   []string          `json:"patterns" binding:"required"`
	MatchType  string            `json:"match_type" binding:"required"`
	Params     map[string]string `json:"params,omitempty"`
	Priority   int               `json:"priority"`
	Confidence float64           `json:"confidence"`
}

type Rule struct {
	ID         string            `json:"id"`
	IntentName string            `json:"intent_name"`
	Patterns   []string          `json:"patterns"`
	MatchType  string            `json:"match_type"`
	Params     map[string]string `json:"params,omitempty"`
	Priority   int               `json:"priority"`
	Confidence float64           `json:"confidence"`
}

type RulesResponse struct {
	Rules []Rule `json:"rules"`
}

type RuleResponse struct {
	Rule Rule `json:"rule"`
}

type DeleteRuleResponse struct {
	Deleted string `json:"deleted"`
}
