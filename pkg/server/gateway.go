package server

import (
	"context"
	"net/http"

	"github.com/go-kod/kod"
	"github.com/grafana/pyroscope-go"
	"github.com/nautilus/gateway"
	"github.com/nautilus/graphql"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type server struct {
	kod.Implements[ServerComponent]

	profiler *pyroscope.Profiler

	config   kod.Ref[ConfigComponent]
	queryer  kod.Ref[Queryer]
	registry kod.Ref[Registry]
}

func (ins *server) Init(ctx context.Context) error {
	profiler, err := ins.config.Get().Config().Pyroscope.Build(ctx)
	if err != nil {
		return err
	}

	ins.profiler = profiler

	return nil
}

func (ins *server) Shutdown(ctx context.Context) error {
	if ins.profiler != nil {
		ins.profiler.Stop()
	}

	return nil
}

func (s *server) BuildServer(ctx context.Context) (http.Handler, error) {
	queryFactory := gateway.QueryerFactory(func(ctx *gateway.PlanningContext, url string) graphql.Queryer {
		return s.queryer.Get()
	})

	sources := []*graphql.RemoteSchema{{URL: "url1"}}
	sources[0].Schema = s.registry.Get().SchemaDescriptorList().AsGraphql()[0]

	g, err := gateway.New(sources, gateway.WithAutomaticQueryPlanCache(), gateway.WithQueryerFactory(&queryFactory))
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/query", g.GraphQLHandler)

	cfg := s.config.Get().Config()
	if cfg.Playground {
		mux.HandleFunc("/playground", g.PlaygroundHandler)
	}

	var handler http.Handler = mux
	handler = otelhttp.NewMiddleware("graphql-gateway")(handler)

	return cors.New(cfg.Cors).Handler(handler), nil
}
