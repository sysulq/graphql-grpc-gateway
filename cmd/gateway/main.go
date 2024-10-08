package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/samber/lo"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
)

type app struct {
	kod.Implements[kod.Main]

	config kod.Ref[config.Config]
	server kod.Ref[server.Gateway]
}

func run(ctx context.Context, app *app) error {
	cfg := app.config.Get().Config()

	l, err := net.Listen("tcp", cfg.GraphQL.Address)
	lo.Must0(err)
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler, err := app.server.Get().BuildServer()
	lo.Must0(err)

	lo.Must0(http.Serve(l, handler))

	return nil
}

func main() {
	_ = kod.Run(context.Background(), run,
		kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()),
	)
}
