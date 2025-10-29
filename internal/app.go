package internal

import (
	"context"

	"github.com/capcom6/hook-relay/internal/config"
	"github.com/capcom6/hook-relay/internal/server"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		logger.Module(),
		logger.WithFxDefaultLogger(),
		module(),
	).Run()
}

func module() fx.Option {
	return fx.Module(
		"app",
		config.Module(),
		fiberfx.Module(),
		server.Module(),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	)
}
