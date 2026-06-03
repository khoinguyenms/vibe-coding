package module

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vibe-be/internal/config"
	"github.com/vibe-be/internal/handler"
	"github.com/vibe-be/internal/router"
	"github.com/vibe-be/pkg/logger"
)

var HTTPModule = fx.Module("http",
	fx.Provide(handler.NewUserHandler),
	fx.Provide(router.New),
	fx.Provide(newHTTPServer),
	fx.Invoke(func(*http.Server) {}),
)

func newHTTPServer(lc fx.Lifecycle, cfg *config.Config, engine *gin.Engine, log *logger.Logger) *http.Server {
	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("http server starting", zap.String("addr", srv.Addr))
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("http server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("http server stopping")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
