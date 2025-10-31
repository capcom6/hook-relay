package subscriptions

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"subscriptions",
		logger.WithNamedLogger("subscriptions"),
		fx.Provide(newRepository, fx.Private),
		fx.Provide(NewService),
	)
}
