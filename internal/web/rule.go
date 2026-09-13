package web

import (
	"net/http"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	rulesvc "github.com/ecodeclub/ai-gateway-go/internal/service/rule"
	"github.com/ecodeclub/ai-gateway-go/internal/web/schema"
	"github.com/gin-gonic/gin"
)

// RuleHandler 对应 rule.Service
type RuleHandler struct {
	svc *rulesvc.Service
}

func NewRuleHandler(svc *rulesvc.Service) *RuleHandler {
	return &RuleHandler{svc: svc}
}

func (h *RuleHandler) List(c *gin.Context) {
	rules, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema.RulesResponse{Rules: toRules(rules)})
}

func (h *RuleHandler) Create(c *gin.Context) {
	var req schema.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.svc.Create(c.Request.Context(), domain.IntentRule{
		IntentName: req.IntentName,
		Patterns:   req.Patterns,
		MatchType:  domain.RuleMatchType(req.MatchType),
		Params:     req.Params,
		Priority:   req.Priority,
		Confidence: req.Confidence,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, schema.RuleResponse{Rule: toRule(created)})
}

func (h *RuleHandler) Update(c *gin.Context) {
	var req schema.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := domain.IntentRule{
		ID:         c.Param("id"),
		IntentName: req.IntentName,
		Patterns:   req.Patterns,
		MatchType:  domain.RuleMatchType(req.MatchType),
		Params:     req.Params,
		Priority:   req.Priority,
		Confidence: req.Confidence,
	}
	if err := h.svc.Update(c.Request.Context(), r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema.RuleResponse{Rule: toRule(r)})
}

func (h *RuleHandler) Delete(c *gin.Context) {
	ruleID := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema.DeleteRuleResponse{Deleted: ruleID})
}

func toRule(r domain.IntentRule) schema.Rule {
	return schema.Rule{
		ID:         r.ID,
		IntentName: r.IntentName,
		Patterns:   r.Patterns,
		MatchType:  string(r.MatchType),
		Params:     r.Params,
		Priority:   r.Priority,
		Confidence: r.Confidence,
	}
}

func toRules(rules []domain.IntentRule) []schema.Rule {
	out := make([]schema.Rule, 0, len(rules))
	for _, r := range rules {
		out = append(out, toRule(r))
	}
	return out
}
