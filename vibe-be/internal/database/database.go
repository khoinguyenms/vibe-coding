package database

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vibe-be/internal/config"
	"github.com/vibe-be/pkg/logger"
)

func NewPool(lc fx.Lifecycle, cfg config.DatabaseConfig, log *logger.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	poolCfg.MaxConns = 25
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTrimSQLInSpanName(),
	)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := pool.Ping(pingCtx); err != nil {
				return fmt.Errorf("ping db: %w", err)
			}
			log.Info("database connected", zap.String("host", cfg.Host), zap.String("db", cfg.Name))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("closing database pool")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
