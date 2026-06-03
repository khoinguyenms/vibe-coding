package module

import (
	"context"

	"go.uber.org/fx"

	"github.com/vibe-be/internal/config"
	"github.com/vibe-be/pkg/logger"
)

var LoggerModule = fx.Module("logger",
	fx.Provide(func(cfg *config.Config) (*logger.Logger, error) {
		return logger.New(cfg.LogLevel, cfg.AppEnv)
	}),
	// fx.WithLogger(func(l *logger.Logger) fxevent.Logger {
	// 	return &fxevent.ZapLogger{Logger: l.Named("fx").Raw()}
	// }),
	fx.Invoke(func(lc fx.Lifecycle, l *logger.Logger) {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				_ = l.Sync()
				return nil
			},
		})
	}),
)
