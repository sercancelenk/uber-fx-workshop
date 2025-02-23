package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net"
	"net/http"
	"uber-fx-workshop/route"
	"uber-fx-workshop/service"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewHttpServer,
			fx.Annotate(
				route.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			route.AsRoute(route.NewEchoHandler),
			route.AsRoute(route.NewHelloHandler),
			zap.NewExample,
			fx.Annotate(
				service.NewConsumer1,
			),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
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
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					panic("Server could not started.")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
