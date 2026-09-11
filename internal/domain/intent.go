package domain

type IntentRequest struct {
	Text string
}

type Intent struct {
	Name       string
	Confidence float64
}

type IntentResponse struct {
	Intents []Intent
}
