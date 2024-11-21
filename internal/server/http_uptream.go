package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-kod/grpc-gateway/internal/config"
	"github.com/go-kod/grpc-gateway/pkg/protojson"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/jhump/protoreflect/grpcreflect"

	// nolint
	"github.com/golang/protobuf/jsonpb"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

type httpUpstream struct {
	kod.Implements[HttpUpstream]

	invoker kod.Ref[HttpUpstreamInvoker]
	config  kod.Ref[config.Config]

	upstreams map[string]upstreamInfo
}

type upstreamInfo struct {
	target   string
	conn     grpc.ClientConnInterface
	source   grpcurl.DescriptorSource
	resovler jsonpb.AnyResolver
	methods  []protojson.Method
}

func (u *httpUpstream) Init(ctx context.Context) error {
	u.upstreams = make(map[string]upstreamInfo)

	lo.ForEach(u.config.Get().Config().Grpc.Services, func(item kgrpc.Config, index int) {
		registry := lo.Must(u.config.Get().Config().Grpc.Etcd.Build(ctx))

		conn := item.WithRegistry(registry).Build()

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		client := grpcreflect.NewClientAuto(ctx, conn)

		source := grpcurl.DescriptorSourceFromServer(ctx, client)

		methods := lo.Must(protojson.GetMethods(source))

		u.upstreams[item.Target] = upstreamInfo{
			target:   item.Target,
			conn:     conn,
			source:   source,
			resovler: grpcurl.AnyResolverFromDescriptorSource(source),
			methods:  methods,
		}
	})

	return nil
}

func (u *httpUpstream) Register(ctx context.Context, router *http.ServeMux) {
	for _, upstream := range u.upstreams {
		for _, v := range upstream.methods {
			if v.HttpPath == "" {
				continue
			}

			u.L(ctx).Info("register upstream", "http", v.HttpPath, "rpc", v.RpcPath)
			router.Handle(fmt.Sprintf("%s %s", v.HttpMethod, v.HttpPath), u.buildHandler(ctx, upstream, v.RpcPath, v.PathNames))
		}
	}
}

func (u *httpUpstream) buildHandler(_ context.Context, upstream upstreamInfo, rpcPath string, pathNames []string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		u.invoker.Get().Invoke(r.Context(), rw, r, upstream, rpcPath, pathNames)
	}
}
