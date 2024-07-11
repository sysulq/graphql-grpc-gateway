package server

import (
	"context"
	"sync"

	"github.com/go-kod/kod"
	"github.com/sysulq/graphql-gateway/pkg/generator"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type repository struct {
	kod.Implements[Registry]

	caller kod.Ref[Caller]
	files  generator.SchemaDescriptorList
	mu     *sync.RWMutex

	methodsByName            map[ast.Operation]map[string]protoreflect.MethodDescriptor
	objectsByName            map[string]protoreflect.MessageDescriptor
	objectsByFQN             map[string]*generator.ObjectDescriptor
	graphqlFieldsByName      map[protoreflect.Descriptor]map[string]protoreflect.FieldDescriptor
	protoFieldsByName        map[*ast.Definition]map[string]*ast.FieldDefinition
	graphqlUnionFieldsByName map[string]map[string]protoreflect.FieldDescriptor
}

func (v *repository) Init(ctx context.Context) error {
	v.mu = &sync.RWMutex{}
	v.methodsByName = map[ast.Operation]map[string]protoreflect.MethodDescriptor{}
	v.objectsByName = map[string]protoreflect.MessageDescriptor{}
	v.objectsByFQN = map[string]*generator.ObjectDescriptor{}
	v.graphqlFieldsByName = map[protoreflect.Descriptor]map[string]protoreflect.FieldDescriptor{}
	v.graphqlUnionFieldsByName = map[string]map[string]protoreflect.FieldDescriptor{}
	v.protoFieldsByName = map[*ast.Definition]map[string]*ast.FieldDefinition{}

	descs := v.caller.Get().GetDescs()

	gqlDesc, err := generator.NewSchemas(descs, true, true, nil)
	if err != nil {
		return err
	}

	v.files = gqlDesc

	for _, f := range gqlDesc {
		v.methodsByName[ast.Mutation] = map[string]protoreflect.MethodDescriptor{}
		for _, m := range f.GetMutation().Methods() {
			v.methodsByName[ast.Mutation][string(m.FieldDefinition.Name)] = m.MethodDescriptor
		}
		v.methodsByName[ast.Query] = map[string]protoreflect.MethodDescriptor{}
		for _, m := range f.GetQuery().Methods() {
			v.methodsByName[ast.Query][string(m.FieldDefinition.Name)] = m.MethodDescriptor
		}
		v.methodsByName[ast.Subscription] = map[string]protoreflect.MethodDescriptor{}
		for _, m := range f.GetSubscription().Methods() {
			v.methodsByName[ast.Subscription][string(m.FieldDefinition.Name)] = m.MethodDescriptor
		}
	}
	for _, f := range gqlDesc {
		for _, m := range f.Objects() {
			switch m.Kind {
			case ast.Union:
				fqn := string(m.Descriptor.FullName())
				if _, ok := v.graphqlUnionFieldsByName[fqn]; !ok {
					v.graphqlUnionFieldsByName[fqn] = map[string]protoreflect.FieldDescriptor{}
				}
				for _, tt := range m.GetTypes() {
					for _, f := range tt.GetFields() {
						v.graphqlUnionFieldsByName[fqn][f.FieldDefinition.Name] = f.FieldDescriptor
					}
				}
			case ast.Object:
				v.protoFieldsByName[m.Definition] = map[string]*ast.FieldDefinition{}
				for _, f := range m.GetFields() {
					if f.FieldDescriptor != nil {
						v.protoFieldsByName[m.Definition][f.FieldDefinition.Name] = f.FieldDefinition
					}
				}
			case ast.InputObject:
				v.graphqlFieldsByName[m.Descriptor] = map[string]protoreflect.FieldDescriptor{}
				for _, f := range m.GetFields() {
					v.graphqlFieldsByName[m.Descriptor][f.FieldDefinition.Name] = f.FieldDescriptor
				}
			}
			switch msgDesc := m.Descriptor.(type) {
			case protoreflect.MessageDescriptor:
				v.objectsByFQN[string(m.FullName())] = m
				v.objectsByName[m.Definition.Name] = msgDesc
			}
		}
	}

	return nil
}

func (r *repository) SchemaDescriptorList() generator.SchemaDescriptorList {
	return r.files
}

func (r *repository) FindMethodByName(op ast.Operation, name string) protoreflect.MethodDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.methodsByName[op][name]
}

func (r *repository) FindObjectByName(name string) protoreflect.MessageDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.objectsByName[name]
}

func (r *repository) FindObjectByFullyQualifiedName(fqn string) (protoreflect.MessageDescriptor, *ast.Definition) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := r.objectsByFQN[fqn]
	msg, _ := o.Descriptor.(protoreflect.MessageDescriptor)
	return msg, o.Definition
}

func (r *repository) FindFieldByName(msg protoreflect.Descriptor, name string) protoreflect.FieldDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.graphqlFieldsByName[msg][name]
}

func (r *repository) FindUnionFieldByMessageFQNAndName(fqn, name string) protoreflect.FieldDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.graphqlUnionFieldsByName[fqn][name]
}

func (r *repository) FindGraphqlFieldByProtoField(msg *ast.Definition, name string) *ast.FieldDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.protoFieldsByName[msg][name]
}
