package server

import (
	"context"
	"sync"

	"github.com/go-kod/kod"
	"github.com/jhump/protoreflect/desc"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/pkg/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

type repository struct {
	kod.Implements[Registry]

	caller kod.Ref[Caller]
	config kod.Ref[config.Config]

	files generator.SchemaDescriptorList
	mu    *sync.RWMutex

	methodsByName            map[ast.Operation]map[string]*desc.MethodDescriptor
	objectsByName            map[string]*desc.MessageDescriptor
	objectsByFQN             map[string]*generator.ObjectDescriptor
	graphqlFieldsByName      map[desc.Descriptor]map[string]*desc.FieldDescriptor
	protoFieldsByName        map[*ast.Definition]map[string]*ast.FieldDefinition
	graphqlUnionFieldsByName map[string]map[string]*desc.FieldDescriptor
}

func (v *repository) Init(ctx context.Context) error {
	v.mu = &sync.RWMutex{}
	v.methodsByName = map[ast.Operation]map[string]*desc.MethodDescriptor{}
	v.objectsByName = map[string]*desc.MessageDescriptor{}
	v.objectsByFQN = map[string]*generator.ObjectDescriptor{}
	v.graphqlFieldsByName = map[desc.Descriptor]map[string]*desc.FieldDescriptor{}
	v.graphqlUnionFieldsByName = map[string]map[string]*desc.FieldDescriptor{}
	v.protoFieldsByName = map[*ast.Definition]map[string]*ast.FieldDefinition{}

	descs := v.caller.Get().GetDescs()

	gqlDesc, err := generator.NewSchemas(descs, generator.Options{
		GenServiceDesc:         true,
		MergeSchemas:           true,
		GenerateUnboundMethods: v.config.Get().Config().GraphQL.GenerateUnboundMethods,
	})
	if err != nil {
		return err
	}

	v.files = gqlDesc

	for _, f := range gqlDesc {
		v.methodsByName[ast.Mutation] = map[string]*desc.MethodDescriptor{}
		for _, m := range f.Methods(ast.Mutation) {
			v.methodsByName[ast.Mutation][m.Name] = m.MethodDescriptor
		}
		v.methodsByName[ast.Query] = map[string]*desc.MethodDescriptor{}
		for _, m := range f.Methods(ast.Query) {
			v.methodsByName[ast.Query][m.Name] = m.MethodDescriptor
		}
		v.methodsByName[ast.Subscription] = map[string]*desc.MethodDescriptor{}
		for _, m := range f.Methods(ast.Subscription) {
			v.methodsByName[ast.Subscription][m.Name] = m.MethodDescriptor
		}
	}
	for _, f := range gqlDesc {
		for _, m := range f.Objects() {
			switch m.Kind {
			case ast.Union:
				fqn := m.Descriptor.GetParent().GetFullyQualifiedName()
				if _, ok := v.graphqlUnionFieldsByName[fqn]; !ok {
					v.graphqlUnionFieldsByName[fqn] = map[string]*desc.FieldDescriptor{}
				}
				for _, tt := range m.GetTypes() {
					for _, f := range tt.GetFields() {
						v.graphqlUnionFieldsByName[fqn][f.Name] = f.FieldDescriptor
					}
				}
			case ast.Object:
				v.protoFieldsByName[m.Definition] = map[string]*ast.FieldDefinition{}
				for _, f := range m.GetFields() {
					if f.FieldDescriptor != nil {
						v.protoFieldsByName[m.Definition][f.FieldDescriptor.GetName()] = f.FieldDefinition
					}
				}
			case ast.InputObject:
				v.graphqlFieldsByName[m.Descriptor] = map[string]*desc.FieldDescriptor{}
				for _, f := range m.GetFields() {
					v.graphqlFieldsByName[m.Descriptor][f.Name] = f.FieldDescriptor
				}
			}
			switch msgDesc := m.Descriptor.(type) {
			case *desc.MessageDescriptor:
				v.objectsByFQN[m.GetFullyQualifiedName()] = m
				v.objectsByName[m.Name] = msgDesc
			}
		}
	}

	return nil
}

func (r *repository) SchemaDescriptorList() generator.SchemaDescriptorList {
	return r.files
}

func (r *repository) FindMethodByName(op ast.Operation, name string) *desc.MethodDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.methodsByName[op][name]
}

func (r *repository) FindObjectByName(name string) *desc.MessageDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.objectsByName[name]
}

func (r *repository) FindObjectByFullyQualifiedName(fqn string) (*desc.MessageDescriptor, *ast.Definition) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := r.objectsByFQN[fqn]
	msg, _ := o.Descriptor.(*desc.MessageDescriptor)
	return msg, o.Definition
}

func (r *repository) FindFieldByName(msg desc.Descriptor, name string) *desc.FieldDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.graphqlFieldsByName[msg][name]
}

func (r *repository) FindUnionFieldByMessageFQNAndName(fqn, name string) *desc.FieldDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.graphqlUnionFieldsByName[fqn][name]
}

func (r *repository) FindGraphqlFieldByProtoField(msg *ast.Definition, name string) *ast.FieldDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.protoFieldsByName[msg][name]
}
