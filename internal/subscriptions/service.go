package subscriptions

import (
	"context"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Service struct {
	subscriptions *repository

	logger *zap.Logger
}

func NewService(subscriptions *repository, logger *zap.Logger) *Service {
	return &Service{
		subscriptions: subscriptions,
		logger:        logger,
	}
}

func (s *Service) Replace(ctx context.Context, subscription SubscriptionIn) (*Subscription, error) {
	model := newSubscriptionModel(subscription)

	err := s.subscriptions.Replace(ctx, model)
	if err != nil {
		return nil, err
	}

	return &Subscription{
		UUID:      model.UUID,
		URL:       model.URL,
		Events:    lo.Map(model.Events, func(event eventModel, _ int) string { return event.Event }),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (s *Service) SelectByUUID(ctx context.Context, uuid ...string) ([]Subscription, error) {
	models, err := s.subscriptions.SelectByUUID(ctx, uuid...)
	if err != nil {
		return nil, err
	}

	if len(models) == 0 {
		return nil, nil
	}

	return lo.Map(models, func(model subscriptionModel, _ int) Subscription {
		return Subscription{
			UUID:      model.UUID,
			URL:       model.URL,
			Events:    lo.Map(model.Events, func(event eventModel, _ int) string { return event.Event }),
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}), nil
}
