package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server 只聚合各 Handler，负责路由注册。
type Server struct {
	llm    *LLMHandler
	intent *IntentHandler
	rule   *RuleHandler
}

func NewServer(llm *LLMHandler, intent *IntentHandler, rule *RuleHandler) *Server {
	return &Server{llm: llm, intent: intent, rule: rule}
}

func (s *Server) Router(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group("/api/v1")
	{
		api.POST("/llm/invoke", s.llm.Invoke)
		api.POST("/llm/stream", s.llm.Stream)

		api.POST("/intent/recognize", s.intent.Recognize)
		api.POST("/intent/clarify", s.intent.Clarify)

		rules := api.Group("/rules")
		{
			rules.GET("", s.rule.List)
			rules.POST("", s.rule.Create)
			rules.PUT("/:id", s.rule.Update)
			rules.DELETE("/:id", s.rule.Delete)
		}
	}
}
