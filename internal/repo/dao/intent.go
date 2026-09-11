package dao

import (
	"context"

	"github.com/ecodeclub/ai-gateway-go/ent"
)

type IntentDao struct {
	client *ent.Client
}

func NewIntentDao(client *ent.Client) *IntentDao {
	return &IntentDao{client: client}
}

func (intent *IntentDao) Create(ctx context.Context) error {
	return nil
}

func (intent *IntentDao) Update(ctx context.Context) error {
	return nil
}
