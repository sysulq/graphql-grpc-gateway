package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/go-kod/grpc-gateway/internal/config"
	"github.com/go-kod/grpc-gateway/internal/server"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/samber/lo"
)

type app struct {
	kod.Implements[kod.Main]

	config kod.Ref[config.Config]
	server kod.Ref[server.Gateway]
}

func run(ctx context.Context, app *app) error {
	cfg := app.config.Get().Config()

	l := lo.Must(net.Listen("tcp", cfg.Server.GraphQL.Address))
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler := lo.Must(app.server.Get().BuildServer())
	go func() { lo.Must0(http.Serve(l, handler)) }()

	l = lo.Must(net.Listen("tcp", cfg.Server.HTTP.Address))
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler = lo.Must(app.server.Get().BuildHTTPServer())
	lo.Must0(http.Serve(l, handler))

	return nil
}

func main() {
	kod.MustRun(context.Background(), run,
		kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()),
	)
}
