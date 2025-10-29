package config

import (
	"fmt"

	"github.com/capcom6/hook-relay/internal/server"
	"github.com/go-core-fx/config"
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
		fx.Provide(func(c Config) server.Config {
			return server.Config{
				Address:     c.Server.Address,
				ProxyHeader: c.Server.ProxyHeader,
				Proxies:     c.Server.Proxies,
			}
		}),
	)
}
