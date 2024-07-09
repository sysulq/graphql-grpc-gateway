package server

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/sysulq/graphql-gateway/pkg/protoparser"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"

	"github.com/sysulq/graphql-gateway/pkg/reflection"
)

type Grpc struct {
	Services    []*Service
	ImportPaths []string
}

type caller struct {
	kod.Implements[Caller]

	config      kod.Ref[ConfigComponent]
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
			newDescs, err = reflection.NewClientWithImportsResolver(conn).ListPackages()
			if err != nil {
				return err
			}
		}
		if e.ProtoFiles != nil {
			descs, err := protoparser.Parse(config.ImportPaths, e.ProtoFiles)
			if err != nil {
				return err
			}
			newDescs = append(newDescs, descs...)
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

	c.descs = descs
	c.serviceStub = serviceStub

	return nil
}

func (c *caller) GetDescs() []*desc.FileDescriptor {
	return c.descs
}

func (c *caller) Call(ctx context.Context, rpc *desc.MethodDescriptor, message proto.Message) (proto.Message, error) {
	res, err := c.serviceStub[rpc.GetService().GetFullyQualifiedName()].InvokeRpc(ctx, rpc, message)
	return res, err
}

func (c *caller) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kcircuitbreaker.Interceptor(),
	}
}
