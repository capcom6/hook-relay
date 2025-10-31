package db

import (
	"github.com/go-core-fx/logger"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func Module() fx.Option {
	return fx.Module(
		"db",
		logger.WithNamedLogger("db"),
		fx.Provide(func() schema.Dialect {
			return mysqldialect.New()
		}),
	)
}
