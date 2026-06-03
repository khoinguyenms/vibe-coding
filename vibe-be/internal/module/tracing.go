package module

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vibe-be/internal/config"
	"github.com/vibe-be/pkg/logger"
	"github.com/vibe-be/pkg/tracing"
)

var TracingModule = fx.Module("tracing",
	fx.Provide(newTracerProvider),
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)

func newTracerProvider(lc fx.Lifecycle, cfg *config.Config, log *logger.Logger) (*sdktrace.TracerProvider, error) {
	if !cfg.Tracing.Enabled {
		log.Info("tracing disabled")
		tp := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		))
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error { return tp.Shutdown(ctx) },
		})
		return tp, nil
	}

	tp, err := tracing.NewTracerProvider(context.Background(), tracing.ProviderConfig{
		Endpoint:     cfg.Tracing.Endpoint,
		ServiceName:  cfg.Tracing.ServiceName,
		Environment:  cfg.AppEnv,
		SamplerRatio: cfg.Tracing.SamplerRatio,
	})
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Error("otel error", zap.Error(err))
	}))

	log.Info("tracing enabled",
		zap.String("endpoint", cfg.Tracing.Endpoint),
		zap.String("service", cfg.Tracing.ServiceName),
		zap.Float64("sampler_ratio", cfg.Tracing.SamplerRatio),
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			log.Info("flushing traces")
			return tp.Shutdown(shutdownCtx)
		},
	})

	return tp, nil
}
