package server

import (
	"context"
	"strconv"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/samber/lo"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
)

type caller struct {
	kod.Implements[Caller]

	singleflight singleflight.Group

	reflection  kod.Ref[Reflection]
	config      kod.Ref[config.Config]
	serviceStub map[string]grpcdynamic.Stub

	descs []*desc.FileDescriptor
}

func (c *caller) Init(ctx context.Context) (err error) {
	config := c.config.Get().Config().Grpc

	serviceStub := map[string]grpcdynamic.Stub{}
	descs := make([]*desc.FileDescriptor, 0)
	descsconn := map[string]*grpc.ClientConn{}

	for _, e := range config.Services {

		options := []grpc.DialOption{
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		}
		if e.Authentication != nil && e.Authentication.Tls != nil {
			cred, err := credentials.NewServerTLSFromFile(e.Authentication.Tls.Certificate, e.Authentication.Tls.PrivateKey)
			if err != nil {
				return err
			}
			options = append(options, grpc.WithTransportCredentials(cred))
		} else {
			options = append(options, grpc.WithCredentialsBundle(insecure.NewBundle()))
		}
		conn, err := grpc.NewClient(e.Address, options...)
		if err != nil {
			return err
		}

		var newDescs []*desc.FileDescriptor
		if e.Reflection {
			newDescs, err = c.reflection.Get().ListPackages(ctx, conn)
			if err != nil {
				return err
			}
		}

		for _, d := range newDescs {
			descsconn[d.GetFullyQualifiedName()] = conn
		}
		descs = append(descs, newDescs...)
	}

	for _, d := range descs {
		for _, svc := range d.GetServices() {
			serviceStub[svc.GetFullyQualifiedName()] = grpcdynamic.NewStub(descsconn[d.GetFullyQualifiedName()])
		}
	}

	descs = lo.UniqBy(descs, func(item *desc.FileDescriptor) string {
		return item.GetFile().GetFullyQualifiedName()
	})

	c.descs = descs
	c.serviceStub = serviceStub

	return nil
}

func (c *caller) GetDescs() []*desc.FileDescriptor {
	return c.descs
}

func (c *caller) Call(ctx context.Context, rpc *desc.MethodDescriptor, message protoadapt.MessageV1) (protoadapt.MessageV1, error) {
	if enable, ok := ctx.Value(allowSingleFlightKey).(bool); ok && enable {
		hash := Hash64.Get()
		defer Hash64.Put(hash)

		msg, err := proto.Marshal(protoadapt.MessageV2Of(message))
		if err != nil {
			return nil, err
		}

		// generate hash based on rpc pointer
		_, err = hash.Write([]byte(rpc.GetFullyQualifiedName()))
		if err != nil {
			return nil, err
		}
		_, err = hash.Write(msg)
		if err != nil {
			return nil, err
		}
		sum := hash.Sum64()
		key := strconv.FormatUint(sum, 10)

		res, err, _ := c.singleflight.Do(key, func() (interface{}, error) {
			return c.serviceStub[rpc.GetService().GetFullyQualifiedName()].InvokeRpc(ctx, rpc, message)
		})
		if err != nil {
			return nil, err
		}

		return res.(protoadapt.MessageV1), nil
	}

	res, err := c.serviceStub[rpc.GetService().GetFullyQualifiedName()].InvokeRpc(ctx, rpc, message)
	return res, err
}

func (c *caller) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kcircuitbreaker.Interceptor(),
	}
}

var allowSingleFlightKey struct{}
