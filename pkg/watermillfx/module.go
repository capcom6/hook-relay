package watermillfx

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"watermill",
		logger.WithNamedLogger("watermill"),
		fx.Provide(NewLogger),
	)
}
