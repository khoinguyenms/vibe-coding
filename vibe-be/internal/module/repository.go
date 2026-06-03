package module

import (
	"go.uber.org/fx"

	"github.com/vibe-be/internal/repository"
)

var RepositoryModule = fx.Module("repository",
	fx.Provide(repository.NewUserRepository),
)
