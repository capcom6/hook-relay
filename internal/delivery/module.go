package delivery

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"delivery",
		logger.WithNamedLogger("delivery"),
		fx.Provide(NewService),
		fx.Invoke(func(lc fx.Lifecycle, svc *Service) {
			ctx, cancel := context.WithCancel(context.Background())
			waitCh := make(chan struct{})
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						defer close(waitCh)
						if err := svc.Run(ctx); err != nil {
							cancel()
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					cancel()
					select {
					case <-waitCh:
					case <-ctx.Done():
					}
					return nil
				},
			})
		}),
	)
}
