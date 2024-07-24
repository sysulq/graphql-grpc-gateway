package server

import (
	"context"
	"sync"

	"github.com/go-kod/kod"
	"github.com/jhump/protoreflect/desc"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/pkg/generator"
	"github.com/sysulq/graphql-grpc-gateway/pkg/protographql"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type repository struct {
	kod.Implements[Registry]

	caller kod.Ref[Caller]
	config kod.Ref[config.Config]

	mu *sync.RWMutex

	methodsByName            map[ast.Operation]map[string]protoreflect.MethodDescriptor
	objectsByName            map[string]protoreflect.MessageDescriptor
	objectsByFQN             map[string]*generator.ObjectDescriptor
	graphqlFieldsByName      map[desc.Descriptor]map[string]protoreflect.FieldDescriptor
	protoFieldsByName        map[*ast.Definition]map[string]*ast.FieldDefinition
	graphqlUnionFieldsByName map[string]map[string]protoreflect.FieldDescriptor

	schema *protographql.SchemaDescriptor
}

func (v *repository) Init(ctx context.Context) error {
	v.mu = &sync.RWMutex{}
	v.methodsByName = map[ast.Operation]map[string]protoreflect.MethodDescriptor{}
	v.objectsByName = map[string]protoreflect.MessageDescriptor{}
	v.objectsByFQN = map[string]*generator.ObjectDescriptor{}
	v.graphqlFieldsByName = map[desc.Descriptor]map[string]protoreflect.FieldDescriptor{}
	v.graphqlUnionFieldsByName = map[string]map[string]protoreflect.FieldDescriptor{}
	v.protoFieldsByName = map[*ast.Definition]map[string]*ast.FieldDefinition{}

	descs := v.caller.Get().GetDescs()

	v.schema = protographql.New()

	for _, desc := range descs {
		err := v.schema.GenerateFile(v.config.Get().Config().GraphQL.GenerateUnboundMethods, desc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) SchemaDescriptorList() *protographql.SchemaDescriptor {
	return r.schema
}

func (r *repository) FindMethodByName(op ast.Operation, name string) protoreflect.MethodDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.schema.MethodsByName[op][name]
}
