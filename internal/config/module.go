package config

import (
	"fmt"

	"github.com/capcom6/hook-relay/internal/server"
	"github.com/go-core-fx/config"
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
		fx.Provide(func(c Config) server.Config {
			return server.Config{
				Address:     c.Server.Address,
				ProxyHeader: c.Server.ProxyHeader,
				Proxies:     c.Server.Proxies,
			}
		}),
		fx.Provide(func(c Config) sqlfx.Config {
			return sqlfx.Config{
				URL:             c.Database.URL,
				ConnMaxIdleTime: c.Database.ConnMaxIdleTime,
				ConnMaxLifetime: c.Database.ConnMaxLifetime,
				MaxOpenConns:    c.Database.MaxOpenConns,
				MaxIdleConns:    c.Database.MaxIdleConns,
			}
		}),
	)
}
