package server

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/jhump/protoreflect/v2/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/sysulq/graphql-gateway/pkg/reflection"
)

type Grpc struct {
	Services    []*Service
	ImportPaths []string
}

type caller struct {
	kod.Implements[Caller]

	singleflight singleflight.Group

	config      kod.Ref[ConfigComponent]
	serviceStub map[string]*grpcdynamic.Stub

	descs []protoreflect.FileDescriptor
}

func (c *caller) Init(ctx context.Context) (err error) {
	config := c.config.Get().Config().Grpc

	serviceStub := map[string]*grpcdynamic.Stub{}
	descs := make([]protoreflect.FileDescriptor, 0)
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

		var newDescs []protoreflect.FileDescriptor
		if e.Reflection {
			newDescs, err = reflection.NewClient(conn).ListPackages()
			if err != nil {
				return err
			}
		}
		// if e.ProtoFiles != nil {
		// 	descs, err := protoparser.Parse(config.ImportPaths, e.ProtoFiles)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	newDescs = append(newDescs, descs...)
		// }
		for _, d := range newDescs {
			for i := 0; i < d.Services().Len(); i++ {
				svc := d.Services().Get(i)
				descsconn[string(svc.FullName())] = conn
			}
		}
		descs = append(descs, newDescs...)
	}

	for _, d := range descs {
		for i := 0; i < d.Services().Len(); i++ {
			svc := d.Services().Get(i)
			serviceStub[string(svc.FullName())] = grpcdynamic.NewStub(
				descsconn[string(svc.FullName())], grpcdynamic.WithResolver(new(protoregistry.Types)))
		}
	}

	c.descs = descs
	c.serviceStub = serviceStub

	return nil
}

func (c *caller) GetDescs() []protoreflect.FileDescriptor {
	return c.descs
}

func (c *caller) Call(ctx context.Context, rpc protoreflect.MethodDescriptor, message proto.Message) (proto.Message, error) {
	if enable, ok := ctx.Value(allowSingleFlightKey).(bool); ok && enable {
		hash := Hash64.Get()
		defer Hash64.Put(hash)

		msg, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}

		// generate hash based on rpc pointer
		hash.Write([]byte(serviceName(string(rpc.FullName()))))
		hash.Write(msg)
		sum := hash.Sum64()
		key := strconv.FormatUint(sum, 10)

		res, err, _ := c.singleflight.Do(key, func() (interface{}, error) {
			return c.serviceStub[serviceName(string(rpc.FullName()))].InvokeRpc(ctx, rpc, message)
		})
		if err != nil {
			return nil, err
		}

		return res.(proto.Message), nil
	}

	res, err := c.serviceStub[serviceName(string(rpc.FullName()))].InvokeRpc(ctx, rpc, message)
	return res, err
}

func serviceName(rpcFullName string) string {
	return rpcFullName[:len(rpcFullName)-len(rpcFullName[strings.LastIndex(rpcFullName, "."):])]
}

func (c *caller) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kcircuitbreaker.Interceptor(),
	}
}

var allowSingleFlightKey struct{}
