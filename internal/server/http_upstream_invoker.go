package server

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/sysulq/graphql-grpc-gateway/pkg/protojson"
)

type invoker struct {
	kod.Implements[Invoker]
}

func (i *invoker) Invoke(ctx context.Context, rw http.ResponseWriter, r *http.Request, upstream upstreamInfo, rpcPath string, pathNames []string) {
	parser, err := protojson.NewRequestParser(r, pathNames, upstream.resovler)
	if err != nil {
		i.L(ctx).Error("parse request", "error", err)
		return
	}

	handler := protojson.NewEventHandler(rw, upstream.resovler)

	err = grpcurl.InvokeRPC(ctx, upstream.source, upstream.conn,
		rpcPath,
		protojson.ProcessHeaders(r.Header),
		handler, parser.Next)
	if err != nil {
		i.L(ctx).Error("invoke rpc", "error", err)
	}
}

func (invoker) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kaccesslog.Interceptor(),
		kmetric.Interceptor(),
		ktrace.Interceptor(),
	}
}
