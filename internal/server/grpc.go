package server

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/registry"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/jhump/protoreflect/v2/grpcdynamic"
	"github.com/samber/lo"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/pkg/protographql"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type caller struct {
	kod.Implements[Caller]

	reflection kod.Ref[Reflection]
	config     kod.Ref[config.Config]

	serviceStub  map[string]*grpcdynamic.Stub
	singleflight singleflight.Group

	schema *protographql.SchemaDescriptor
}

func (c *caller) Init(ctx context.Context) (err error) {
	config := c.config.Get().Config().Grpc

	serviceStub := map[string]*grpcdynamic.Stub{}
	descs := make([]protoreflect.FileDescriptor, 0)
	descsconn := map[string]*grpc.ClientConn{}
	var etcd registry.Registry

	if len(c.config.Get().Config().Grpc.Etcd.Endpoints) > 0 {
		etcd = lo.Must(c.config.Get().Config().Grpc.Etcd.Build(ctx))
	}

	for _, service := range config.Services {

		if etcd != nil {
			service = service.WithRegistry(etcd)
		}

		conn := service.Build()

		newDescs, err := c.reflection.Get().ListPackages(ctx, conn)
		if err != nil {
			return err
		}

		for _, d := range newDescs {
			descsconn[string(d.FullName())] = conn
		}
		descs = append(descs, newDescs...)
	}

	for _, d := range descs {
		for i := 0; i < d.Services().Len(); i++ {
			svc := d.Services().Get(i)
			serviceStub[string(svc.FullName())] = grpcdynamic.NewStub(descsconn[string(d.FullName())])
		}
	}

	descs = lo.UniqBy(descs, func(item protoreflect.FileDescriptor) string {
		return string(item.FullName())
	})

	c.serviceStub = serviceStub

	c.schema = protographql.New()

	for _, desc := range descs {
		err := c.schema.RegisterFileDescriptor(c.config.Get().Config().GraphQL.GenerateUnboundMethods, desc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *caller) Call(ctx context.Context, rpc protoreflect.MethodDescriptor, message proto.Message) (proto.Message, error) {
	if c.config.Get().Config().GraphQL.SingleFlight {
		if enable, ok := ctx.Value(allowSingleFlightKey).(bool); ok && enable {
			hash := Hash64.Get()
			defer Hash64.Put(hash)

			md, ok := metadata.FromOutgoingContext(ctx)
			if ok {
				hd := make([]string, 0, len(md))
				for k, v := range md {
					// skip grpc gateway prefixed metadata
					if strings.Contains(k, MetadataPrefix) {
						continue
					}
					hd = append(hd, k+strings.Join(v, ","))
				}
				slices.Sort(hd)
				for _, v := range hd {
					_, err := hash.Write([]byte(v))
					if err != nil {
						return nil, err
					}
				}
			}

			msg, err := proto.Marshal(message)
			if err != nil {
				return nil, err
			}

			// generate hash based on rpc pointer
			_, err = hash.Write([]byte(rpc.FullName()))
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
				return c.serviceStub[string(rpc.Parent().FullName())].InvokeRpc(ctx, rpc, message)
			})
			if err != nil {
				return nil, err
			}

			return res.(proto.Message), nil
		}
	}

	res, err := c.serviceStub[string(rpc.Parent().FullName())].InvokeRpc(ctx, rpc, message)
	return res, err
}

func (c *caller) Interceptors() []interceptor.Interceptor {
	if c.config.Get().Config().Engine.CircuitBreaker {
		return []interceptor.Interceptor{
			kcircuitbreaker.Interceptor(),
		}
	}
	return nil
}

var allowSingleFlightKey struct{}

func (r *caller) FindMethodByName(op ast.Operation, name string) protoreflect.MethodDescriptor {
	return r.schema.MethodsByName[op][name]
}

func (r *caller) GraphQLSchema() *ast.Schema {
	return r.schema.AsGraphQL()
}

func (r *caller) Marshal(proto proto.Message, field *ast.Field) (interface{}, error) {
	return r.schema.Marshal(proto, field)
}

func (r *caller) Unmarshal(desc protoreflect.MessageDescriptor, field *ast.Field, vars map[string]interface{}) (proto.Message, error) {
	return r.schema.Unmarshal(desc, field, vars)
}
