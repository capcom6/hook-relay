package subscriptions

import (
	"context"

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

func (s *Service) Replace(ctx context.Context, draft SubscriptionIn) (*Subscription, error) {
	item, err := s.subscriptions.Replace(ctx, draft)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) SelectByUUID(ctx context.Context, uuid ...string) ([]Subscription, error) {
	items, err := s.subscriptions.SelectByUUID(ctx, uuid...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*Subscription, error) {
	subscriptions, err := s.SelectByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if len(subscriptions) == 0 {
		return nil, ErrNotFound
	}

	return &subscriptions[0], nil
}

func (s *Service) DeleteByUUID(ctx context.Context, uuid string) error {
	return s.subscriptions.DeleteByUUID(ctx, uuid)
}
