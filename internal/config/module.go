package config

import (
	"fmt"

	"github.com/capcom6/hook-relay/internal/delivery"
	"github.com/capcom6/hook-relay/internal/events"
	"github.com/go-core-fx/config"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/sqlfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(func() (Config, error) {
			c := Default()

			if err := config.Load(&c); err != nil {
				return Config{}, fmt.Errorf("failed to load config: %w", err)
			}

			return c, nil
		}, fx.Private),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
			func(c Config) sqlfx.Config {
				return sqlfx.Config{
					URL:             c.Database.URL,
					ConnMaxIdleTime: c.Database.ConnMaxIdleTime,
					ConnMaxLifetime: c.Database.ConnMaxLifetime,
					MaxOpenConns:    c.Database.MaxOpenConns,
					MaxIdleConns:    c.Database.MaxIdleConns,
				}
			},
		),

		fx.Provide(
			func(c Config) events.Config {
				return events.Config{
					Timeout: c.Events.Timeout,
				}
			},
			func(c Config) delivery.Config {
				return delivery.Config{
					Timeout:   c.Delivery.Timeout,
					UserAgent: c.Delivery.UserAgent,
				}
			},
		),
	)
}
