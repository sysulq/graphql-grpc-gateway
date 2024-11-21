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

	l, err := net.Listen("tcp", cfg.Server.GraphQL.Address)
	lo.Must0(err)
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler, err := app.server.Get().BuildServer()
	lo.Must0(err)
	go func() { lo.Must0(http.Serve(l, handler)) }()

	l, err = net.Listen("tcp", cfg.Server.HTTP.Address)
	lo.Must0(err)
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler, err = app.server.Get().BuildHTTPServer()
	lo.Must0(err)
	lo.Must0(http.Serve(l, handler))

	return nil
}

func main() {
	kod.MustRun(context.Background(), run,
		kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()),
	)
}
