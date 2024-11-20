package server

import (
	"context"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/pkg/protojson"

	// nolint
	"github.com/golang/protobuf/jsonpb"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

type upstream struct {
	kod.Implements[Upstream]

	invoker kod.Ref[Invoker]
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

func (u *upstream) Init(ctx context.Context) error {
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

func (u *upstream) Register(ctx context.Context, router *gin.Engine) {
	for _, upstream := range u.upstreams {
		for _, v := range upstream.methods {
			if v.HttpPath == "" {
				continue
			}

			u.L(ctx).Info("register upstream", "http", v.HttpPath, "rpc", v.RpcPath)
			router.Handle(v.HttpMethod, v.HttpPath, u.buildHandler(ctx, upstream, v.RpcPath))
		}
	}
}

func (u *upstream) buildHandler(_ context.Context, upstream upstreamInfo, rpcPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u.invoker.Get().Invoke(c.Request.Context(), c, upstream, rpcPath)
	}
}
