package web

import (
	"net/http"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	intent "github.com/ecodeclub/ai-gateway-go/internal/service/intent"
	"github.com/ecodeclub/ai-gateway-go/internal/web/schema"
	"github.com/gin-gonic/gin"
)

// IntentHandler 对应 intent.Service
type IntentHandler struct {
	svc *intent.Service
}

func NewIntentHandler(svc *intent.Service) *IntentHandler {
	return &IntentHandler{svc: svc}
}

func (h *IntentHandler) Recognize(c *gin.Context) {
	var req schema.RecognizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := h.svc.Recognize(c.Request.Context(), domain.IntentRequest{
		Text:      req.Text,
		Context:   req.Context,
		SessionID: req.SessionID,
	})
	c.JSON(http.StatusOK, toIntentResponse(resp))
}

func (h *IntentHandler) Clarify(c *gin.Context) {
	var req schema.ClarifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := h.svc.Clarify(c.Request.Context(), req.SessionID, req.Answer)
	c.JSON(http.StatusOK, toIntentResponse(resp))
}

func toIntentResponse(resp domain.IntentResponse) schema.IntentResponse {
	out := schema.IntentResponse{
		MatchedLayer:   resp.MatchedLayer,
		NeedClarify:    resp.NeedClarify,
		AutoCompleted:  resp.AutoCompleted,
		TotalLatencyMs: resp.TotalLatencyMs,
		Intents:        make([]schema.Intent, 0, len(resp.Intents)),
	}
	for _, in := range resp.Intents {
		out.Intents = append(out.Intents, schema.Intent{
			Name:       in.Name,
			Confidence: in.Confidence,
			Params:     in.Params,
		})
	}
	if resp.Clarification != nil {
		out.Clarification = &schema.Clarification{
			Question: resp.Clarification.Question,
			Options:  resp.Clarification.Options,
			Hint:     resp.Clarification.Hint,
		}
	}
	return out
}
