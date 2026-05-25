package module

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	sqlcdb "github.com/vibe-be/db/sqlc"
	"github.com/vibe-be/internal/database"
)

var DatabaseModule = fx.Module("database",
	fx.Provide(database.NewPool),
	fx.Provide(func(p *pgxpool.Pool) sqlcdb.DBTX { return p }),
	fx.Provide(
		fx.Annotate(
			sqlcdb.New,
			fx.As(new(sqlcdb.Querier)),
		),
	),
)
