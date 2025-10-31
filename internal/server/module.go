package server

import (
	"github.com/capcom6/hook-relay/internal/server/api/subscriptions"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/statuscode"
	"github.com/go-core-fx/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		logger.WithNamedLogger("server"),
		fx.Provide(func(c Config) fiberfx.Config {
			return fiberfx.Config{
				Address:     c.Address,
				ProxyHeader: c.ProxyHeader,
				Proxies:     c.Proxies,
			}
		}),
		fx.Provide(func() fiberfx.Options {
			return fiberfx.Options{}
		}),
		fx.Provide(
			subscriptions.NewHandler,
			fx.Private,
		),
		fx.Invoke(func(app *fiber.App, subscriptions *subscriptions.Handler) error {
			subscriptions.Register(app.Group("/subscriptions"))

			app.Use(statuscode.New())
			return nil
		}),
	)
}
