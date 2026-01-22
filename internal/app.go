package internal

import (
	"context"

	"github.com/capcom6/hook-relay/internal/config"
	"github.com/capcom6/hook-relay/internal/db"
	"github.com/capcom6/hook-relay/internal/server"
	"github.com/capcom6/hook-relay/internal/subscriptions"
	"github.com/go-core-fx/bunfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/validator"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		sqlfx.Module(),
		goosefx.Module(),
		bunfx.Module(),
		fiberfx.Module(),
		validator.Module(),
		//
		// APP MODULES
		config.Module(),
		server.Module(),
		db.Module(),
		//
		// BUSINESS MODULES
		subscriptions.Module(),
		//
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
	).Run()
}
