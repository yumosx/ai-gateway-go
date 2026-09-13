package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ecodeclub/ai-gateway-go/internal/repo"
	"github.com/ecodeclub/ai-gateway-go/internal/service/intent/rule"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func StartHotReload(lc fx.Lifecycle, matcher *rule.Matcher, ruleRepo *repo.IntentRuleRepo) {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go runHotReload(ctx, matcher, ruleRepo)
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})
}

func runHotReload(ctx context.Context, matcher *rule.Matcher, ruleRepo *repo.IntentRuleRepo) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ruleRepo.SubscribeReload(ctx, func() {
		log.Println("rule reload triggered by pub/sub")
		if err := matcher.ReloadFromRepo(ruleRepo); err != nil {
			log.Printf("rule reload failed: %v", err)
		}
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := matcher.ReloadFromRepo(ruleRepo); err != nil {
				log.Printf("periodic rule reload failed: %v", err)
			}
		}
	}
}

func StartHTTPServer(lc fx.Lifecycle, cfg Config, r *gin.Engine) {
	srv := &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Printf("listening on %s", cfg.HTTP.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("http server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
