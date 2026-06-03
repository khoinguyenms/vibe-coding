package module

import (
	"go.uber.org/fx"

	"github.com/vibe-be/internal/config"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
	fx.Provide(func(c *config.Config) config.DatabaseConfig { return c.Database }),
)
