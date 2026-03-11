package delivery

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/capcom6/go-restkit"
	"github.com/capcom6/hook-relay/internal/events"
	"github.com/capcom6/hook-relay/internal/subscriptions"
	"github.com/go-core-fx/healthfx"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	subscriptionsSvc *subscriptions.Service

	subscriber *events.Subscriber

	client *restkit.Client

	logger *zap.Logger
}

func NewService(
	config Config,
	version healthfx.Version,
	subscriptionsSvc *subscriptions.Service,
	subscriber *events.Subscriber,
	logger *zap.Logger,
) (*Service, error) {
	client, err := restkit.NewClient(restkit.Config{
		Client:  nil,
		BaseURL: "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	config.UserAgent = fmt.Sprintf("%s/%s", config.UserAgent, version.Version)

	return &Service{
		config: config,

		subscriptionsSvc: subscriptionsSvc,

		subscriber: subscriber,

		client: client,

		logger: logger,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	ch, err := s.subscriber.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-ch:
			if !ok {
				return nil // channel closed
			}

			s.logger.Info("event", zap.String("event_id", event.ID), zap.String("event_type", event.Event.Type))

			if handleErr := s.handleEvent(ctx, event); handleErr != nil {
				s.logger.Error("failed to process event", zap.Error(handleErr))
				event.Nack()
			} else {
				event.Ack()
			}
		}
	}
}

func (s *Service) handleEvent(ctx context.Context, event *events.EventWrapper) error {
	subscriptions, err := s.subscriptionsSvc.SelectByEvent(ctx, event.Event.Type)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		return nil
	}

	for _, subscription := range subscriptions {
		if deliverErr := s.deliverEvent(ctx, subscription, event); deliverErr != nil {
			s.logger.Error(
				"failed to deliver event",
				zap.Error(deliverErr),
				zap.String("subscription", subscription.UUID),
				zap.String("event_id", event.ID),
			)
		}
	}

	return nil
}

func (s *Service) signPayload(secret string, payload []byte) (string, error) {
	if secret == "" {
		return "", nil
	}

	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write(payload)
	if err != nil {
		return "", fmt.Errorf("failed to sign payload: %w", err)
	}
	return fmt.Sprintf("sha256=%x", h.Sum(nil)), nil
}

func (s *Service) deliverEvent(
	ctx context.Context,
	subscription subscriptions.Subscription,
	event *events.EventWrapper,
) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	body := event.Event.Payload
	signature, err := s.signPayload(subscription.Secret, body)
	if err != nil {
		return fmt.Errorf("failed to sign payload: %w", err)
	}

	headers := http.Header{
		"Content-Type":        []string{"application/json"},
		"User-Agent":          []string{s.config.UserAgent},
		"X-Request-ID":        []string{event.ID},
		"X-Event-Type":        []string{event.Event.Type},
		"X-Subscription-UUID": []string{subscription.UUID},
	}

	if signature != "" {
		headers.Set("X-Signature", signature)
	}

	if doErr := s.client.DoRAW(
		ctx,
		http.MethodPost,
		subscription.URL,
		headers,
		bytes.NewReader(body),
		nil,
	); doErr != nil {
		return fmt.Errorf("failed to deliver event: %w", doErr)
	}

	return nil
}
