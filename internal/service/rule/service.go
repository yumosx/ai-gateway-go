package rule

import (
	"context"
	"fmt"
	"time"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
	"github.com/ecodeclub/ai-gateway-go/internal/repo"
)

// Service 意图规则业务服务
type Service struct {
	repo *repo.IntentRuleRepo
}

func NewService(r *repo.IntentRuleRepo) *Service {
	return &Service{repo: r}
}

func (s *Service) List(ctx context.Context) ([]domain.IntentRule, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) Create(ctx context.Context, rule domain.IntentRule) (domain.IntentRule, error) {
	if rule.ID == "" {
		rule.ID = fmt.Sprintf("r%d", time.Now().UnixNano())
	}
	if rule.Confidence == 0 {
		rule.Confidence = 1.0
	}
	if err := s.repo.Create(ctx, rule); err != nil {
		return domain.IntentRule{}, err
	}
	return rule, nil
}

func (s *Service) Update(ctx context.Context, rule domain.IntentRule) error {
	return s.repo.Update(ctx, rule)
}

func (s *Service) Delete(ctx context.Context, ruleID string) error {
	return s.repo.Delete(ctx, ruleID)
}
