package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
	"uber-fx-workshop/config"
	"uber-fx-workshop/route"
	"uber-fx-workshop/service"
)

func main() {
	// Create logger
	logger := zap.NewExample()
	defer logger.Sync()

	// Uber Fx already listens for OS signals, so no need for signal.NotifyContext()
	app := fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewHttpServer,
			route.Module,
			route.AsRoute(route.NewEchoHandler),
			route.AsRoute(route.NewHelloHandler),
			zap.NewExample,
			fx.Annotate(
				service.NewConsumer1,
			),
			config.LoadConfig,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	// Use app.Run() to properly block and handle signals (Uber Fx does this automatically)
	app.Run()

	logger.Info("Application shutdown complete.")
}

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":8087", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))

			// Run the server in a goroutine so it doesn't block
			go func() {
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					log.Fatal("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server: No longer accepting new connections...")

			// 1️⃣ Stop accepting new connections immediately
			shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			// 2️⃣ Attempt graceful shutdown
			err := srv.Shutdown(shutdownCtx)
			if err != nil {
				log.Error("HTTP server shutdown error", zap.Error(err))
				return err
			}

			log.Info("HTTP server shutdown complete.")
			return nil
		},
	})
	return srv
}
