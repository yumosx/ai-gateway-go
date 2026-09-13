package schema

type LLMRequest struct {
	ID   string `json:"id"`
	Text string `json:"text" binding:"required"`
}

type LLMResponse struct {
	ReasoningContent string `json:"reasoning_content,omitempty"`
	Content          string `json:"content"`
}

type StreamEvent struct {
	Final            bool   `json:"final"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
	Content          string `json:"content,omitempty"`
	Err              string `json:"err,omitempty"`
}
