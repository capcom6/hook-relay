package events

import (
	"context"
	"errors"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"events",
		logger.WithNamedLogger("events"),
		fx.Provide(NewPublisher),
		fx.Provide(NewSubscriber),

		fx.Invoke(func(lc fx.Lifecycle, pub *Publisher, sub *Subscriber) {
			lc.Append(fx.StopHook(func(_ context.Context) error {
				pubErr := pub.Close()
				subErr := sub.Close()
				return errors.Join(pubErr, subErr)
			}))
		}),
	)
}
