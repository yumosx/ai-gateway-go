package main

import (
	"github.com/ecodeclub/ai-gateway-go/internal/repo"
	"github.com/ecodeclub/ai-gateway-go/internal/repo/cache"
	"github.com/ecodeclub/ai-gateway-go/internal/repo/dao"
	"github.com/ecodeclub/ai-gateway-go/internal/service"
	intent "github.com/ecodeclub/ai-gateway-go/internal/service/intent"
	rulesvc "github.com/ecodeclub/ai-gateway-go/internal/service/rule"
	"github.com/ecodeclub/ai-gateway-go/internal/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// infra
			NewConfig,
			NewEntClient,
			NewRedisClient,
			NewGinEngine,
			NewDeepSeekClient,

			// repo(cache, dao(ent))
			dao.NewIntentRuleDAO,
			cache.NewIntentRuleCache,
			repo.NewIntentRuleRepo,

			// service(repo) / service internals
			NewRuleMatcher,
			AsRuleEngine,
			NewLLMRecognizer,
			AsLLMEngine,
			NewSemanticEngine,
			AsSemanticEngine,
			NewIntentDefinitions,
			intent.NewPipeline,
			NewClarifier,
			intent.NewService,
			rulesvc.NewService,
			NewDeepSeekHandler,
			AsLLMHandler,
			service.NewAIService,

			// handler(service)
			web.NewLLMHandler,
			web.NewIntentHandler,
			web.NewRuleHandler,
			web.NewServer,
		),
		fx.Invoke(
			func(r *gin.Engine, s *web.Server) { s.Router(r) },
			StartHotReload,
			StartHTTPServer,
		),
	).Run()
}
