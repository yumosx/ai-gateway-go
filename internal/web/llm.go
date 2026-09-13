package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	"github.com/ecodeclub/ai-gateway-go/internal/service"
	"github.com/ecodeclub/ai-gateway-go/internal/web/schema"
	"github.com/gin-gonic/gin"
)

// LLMHandler 对应 AIService
type LLMHandler struct {
	svc *service.AIService
}

func NewLLMHandler(svc *service.AIService) *LLMHandler {
	return &LLMHandler{svc: svc}
}

func (h *LLMHandler) Invoke(c *gin.Context) {
	var req schema.LLMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Invoke(c.Request.Context(), domain.LLMRequest{
		Id:   req.ID,
		Text: req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema.LLMResponse{
		ReasoningContent: resp.ReasoningContent,
		Content:          resp.Content,
	})
}

func (h *LLMHandler) Stream(c *gin.Context) {
	var req schema.LLMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ch, err := h.svc.Stream(c.Request.Context(), domain.LLMRequest{
		Id:   req.ID,
		Text: req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			_ = writeSSE(w, schema.StreamEvent{Final: true, Err: ctx.Err().Error()})
			flusher.Flush()
			return
		case e, ok := <-ch:
			if !ok || e.Done {
				_ = writeSSE(w, schema.StreamEvent{Final: true})
				flusher.Flush()
				return
			}
			if e.Error != nil {
				_ = writeSSE(w, schema.StreamEvent{Final: true, Err: e.Error.Error()})
				flusher.Flush()
				return
			}
			if err := writeSSE(w, schema.StreamEvent{
				Final:            false,
				ReasoningContent: e.ReasoningContent,
				Content:          e.Content,
			}); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func writeSSE(w io.Writer, evt schema.StreamEvent) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", b)
	return err
}
