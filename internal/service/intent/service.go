package intent

import (
	"context"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

// Service 意图识别业务服务
type Service struct {
	pipeline  *Pipeline
	clarifier *Clarifier
}

func NewService(pipeline *Pipeline, clarifier *Clarifier) *Service {
	return &Service{
		pipeline:  pipeline,
		clarifier: clarifier,
	}
}

func (s *Service) Recognize(ctx context.Context, req domain.IntentRequest) domain.IntentResponse {
	resp := s.pipeline.Recognize(ctx, req)
	s.clarifier.Evaluate(&resp, req.Context)
	return resp
}

func (s *Service) Clarify(ctx context.Context, sessionID, answer string) domain.IntentResponse {
	resp := s.pipeline.Recognize(ctx, domain.IntentRequest{
		Text:      answer,
		SessionID: sessionID,
	})
	s.clarifier.HandleClarify(sessionID, answer, &resp)
	return resp
}
