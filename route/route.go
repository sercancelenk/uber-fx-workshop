package route

import (
	"go.uber.org/fx"
	"net/http"
)

type Route interface {
	http.Handler
	Pattern() string
}

func AsRoute(f any) interface{} {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

var Module = fx.Annotate(
	NewServeMux,
	fx.ParamTags(`group:"routes"`),
)

// NewServeMux builds a ServeMux that will route requests
// to the given EchoHandler.
func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
