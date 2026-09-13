package schema

type RecognizeRequest struct {
	Text      string            `json:"text" binding:"required"`
	Context   map[string]string `json:"context,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
}

type ClarifyRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}

type Intent struct {
	Name       string            `json:"name"`
	Confidence float64           `json:"confidence"`
	Params     map[string]string `json:"params,omitempty"`
}

type Clarification struct {
	Question string   `json:"question"`
	Options  []string `json:"options,omitempty"`
	Hint     string   `json:"hint,omitempty"`
}

type IntentResponse struct {
	Intents        []Intent          `json:"intents"`
	MatchedLayer   string            `json:"matched_layer"`
	NeedClarify    bool              `json:"need_clarify"`
	Clarification  *Clarification    `json:"clarification,omitempty"`
	AutoCompleted  map[string]string `json:"auto_completed,omitempty"`
	TotalLatencyMs int64             `json:"total_latency_ms"`
}
