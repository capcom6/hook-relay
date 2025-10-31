package internal

import (
	"context"

	"github.com/capcom6/hook-relay/internal/config"
	"github.com/capcom6/hook-relay/internal/db"
	"github.com/capcom6/hook-relay/internal/server"
	"github.com/capcom6/hook-relay/internal/subscriptions"
	"github.com/go-core-fx/bunfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/validator"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		logger.Module(),
		logger.WithFxDefaultLogger(),
		fiberfx.Module(),
		sqlfx.Module(),
		bunfx.Module(),
		validator.Module(),
		module(),
	).Run()
}

func module() fx.Option {
	return fx.Module(
		"app",
		config.Module(),
		server.Module(),
		db.Module(),
		subscriptions.Module(),
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
