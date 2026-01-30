package events

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/capcom6/hook-relay/pkg/watermillfx"
	"go.uber.org/zap"
)

type Subscriber struct {
	config Config

	client message.Subscriber

	logger *zap.Logger
}

func NewSubscriber(config Config, sql *sql.DB, logger *zap.Logger) (*Subscriber, error) {
	client, err := wmsql.NewSubscriber(
		wmsql.BeginnerFromStdSQL(sql),
		wmsql.SubscriberConfig{
			SchemaAdapter:  wmsql.MySQLQueueSchema{},
			OffsetsAdapter: wmsql.MySQLQueueOffsetsAdapter{DeleteOnAck: true},
		},
		watermillfx.NewLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	return &Subscriber{
		config: config,

		client: client,

		logger: logger,
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan *EventWrapper, error) {
	out := make(chan *EventWrapper)

	in, err := s.client.Subscribe(ctx, "events")
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to events: %w", err)
	}

	go func() {
		defer close(out)

		for msg := range in {
			var event Event
			if jsonErr := event.Unmarshal(msg.Payload); jsonErr != nil {
				s.logger.Error("failed to unmarshal event", zap.String("id", msg.UUID), zap.Error(jsonErr))
				msg.Nack()
				continue
			}

			out <- &EventWrapper{
				ID:    msg.UUID,
				Event: event,
				Ack:   msg.Ack,
				Nack:  msg.Nack,
			}

			select {
			case <-msg.Acked():
				continue
			case <-msg.Nacked():
				continue
			case <-time.After(s.config.Timeout):
				s.logger.Warn(
					"event not acked or nacked within timeout",
					zap.String("id", msg.UUID),
					zap.Duration("timeout", s.config.Timeout),
				)
			}
		}
	}()

	return out, nil
}

func (s *Subscriber) Close() error {
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber: %w", err)
	}
	return nil
}
