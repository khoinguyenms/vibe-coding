package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"github.com/vibe-be/internal/module"
)

func main() {
	_ = godotenv.Load()

	fx.New(
		module.ConfigModule,
		module.LoggerModule,
		module.DatabaseModule,
		module.RepositoryModule,
		module.ServiceModule,
		module.HTTPModule,
	).Run()
}
