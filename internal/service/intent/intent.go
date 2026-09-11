package intent

import (
	"context"

	"github.com/ecodeclub/ai-gateway-go/internal/domain"
)

type IntentHandler interface {
	Recognize(ctx context.Context, req domain.IntentRequest) (domain.IntentResponse, error)
}

type IntentServiceImpl struct {
}

func (i *IntentServiceImpl) Recognize(ctx context.Context, req domain.IntentRequest) (domain.IntentResponse, error) {
	return domain.IntentResponse{}, nil
}
