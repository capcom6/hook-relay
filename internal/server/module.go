package server

import (
	"github.com/capcom6/hook-relay/internal/server/api/events"
	"github.com/capcom6/hook-relay/internal/server/api/subscriptions"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/statuscode"
	"github.com/go-core-fx/fiberfx/validation"
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
			fx.Annotate(subscriptions.NewHandler, fx.ResultTags(`group:"handlers"`)), fx.Private,
			fx.Annotate(events.NewHandler, fx.ResultTags(`group:"handlers"`)), fx.Private,
		),

		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, app *fiber.App) {
					// Health endpoint
					// healthHandler.Register(app)

					// Version 1 API group
					v1 := app.Group("api/v1")
					// openapiHandler.Register(v1.Group("/docs"))

					v1.Use(validation.Middleware)

					for _, h := range handlers {
						h.Register(v1)
					}

					app.Use(statuscode.New())
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
