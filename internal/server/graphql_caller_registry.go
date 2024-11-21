package server

import (
	"context"

	"github.com/go-kod/grpc-gateway/internal/config"
	"github.com/go-kod/grpc-gateway/pkg/protographql"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/registry"
	"github.com/jhump/protoreflect/v2/grpcdynamic"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type graphqlCallerRegistry struct {
	kod.Implements[GraphqlCallerRegistry]
	config     kod.Ref[config.Config]
	reflection kod.Ref[GraphqlReflection]

	serviceStub map[string]*grpcdynamic.Stub
	schema      *protographql.SchemaDescriptor
}

func (c *graphqlCallerRegistry) Init(ctx context.Context) error {
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

	return c.setFileDescriptors(descs)
}

func (r *graphqlCallerRegistry) setFileDescriptors(files []protoreflect.FileDescriptor) error {
	schema := protographql.New()
	for _, file := range files {
		err := schema.RegisterFileDescriptor(r.config.Get().Config().Server.GraphQL.GenerateUnboundMethods, file)
		if err != nil {
			return err
		}
	}
	r.schema = schema

	return nil
}

func (r *graphqlCallerRegistry) FindMethodByName(op ast.Operation, name string) protoreflect.MethodDescriptor {
	return r.schema.MethodsByName[op][name]
}

func (r *graphqlCallerRegistry) GraphQLSchema() *ast.Schema {
	return r.schema.AsGraphQL()
}

func (r *graphqlCallerRegistry) Marshal(proto proto.Message, field *ast.Field) (interface{}, error) {
	return r.schema.Marshal(proto, field)
}

func (r *graphqlCallerRegistry) Unmarshal(desc protoreflect.MessageDescriptor, field *ast.Field, vars map[string]interface{}) (proto.Message, error) {
	return r.schema.Unmarshal(desc, field, vars)
}

func (r *graphqlCallerRegistry) GetCallerStub(service string) *grpcdynamic.Stub {
	return r.serviceStub[service]
}
