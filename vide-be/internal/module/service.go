package module

import (
	"go.uber.org/fx"

	"github.com/vibe-be/internal/service"
)

var ServiceModule = fx.Module("service",
	fx.Provide(service.NewUserService),
)
