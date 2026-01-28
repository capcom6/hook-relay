package events

import (
	"database/sql"
	"fmt"

	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/capcom6/hook-relay/pkg/watermillfx"
	"go.uber.org/zap"
)

type Publisher struct {
	client message.Publisher

	logger *zap.Logger
}

func NewPublisher(sql *sql.DB, logger *zap.Logger) (*Publisher, error) {
	client, err := wmsql.NewPublisher(
		wmsql.BeginnerFromStdSQL(sql),
		wmsql.PublisherConfig{
			SchemaAdapter:        wmsql.MySQLQueueSchema{},
			AutoInitializeSchema: true,
		},
		watermillfx.NewLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &Publisher{
		client: client,

		logger: logger,
	}, nil
}

func (p *Publisher) Publish(id string, event Event) error {
	payload, err := event.Marshal()
	if err != nil {
		return err
	}

	msg := message.NewMessage(id, payload)

	if pubErr := p.client.Publish("events", msg); pubErr != nil {
		return fmt.Errorf("failed to publish event: %w", pubErr)
	}

	return nil
}

func (p *Publisher) Close() error {
	if err := p.client.Close(); err != nil {
		return fmt.Errorf("failed to close publisher: %w", err)
	}
	return nil
}
