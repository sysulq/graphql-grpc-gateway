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
	"github.com/sysulq/graphql-gateway/pkg/server"
)

type gateway struct {
	kod.Implements[kod.Main]

	config kod.Ref[server.ConfigComponent]
	server kod.Ref[server.ServerComponent]
}

func run(ctx context.Context, gw *gateway) error {
	cfg := gw.config.Get().Config()

	l, err := net.Listen("tcp", cfg.Address)
	lo.Must0(err)
	log.Printf("[INFO] Gateway listening on address: %s\n", l.Addr())
	handler, err := gw.server.Get().BuildServer()
	lo.Must0(err)
	if cfg.Tls.Enable {
		log.Fatal(http.ServeTLS(l, handler, cfg.Tls.Certificate, cfg.Tls.PrivateKey))
	}

	lo.Must0(http.Serve(l, handler))

	return nil
}

func main() {
	_ = kod.Run(context.Background(), run,
		kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()),
	)
}
