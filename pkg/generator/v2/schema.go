package v2

import (
	"strings"

	graphqlv1 "github.com/sysulq/graphql-grpc-gateway/api/graphql/v1"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// SchemaDescriptor 存储生成的 GraphQL schema 信息
type SchemaDescriptor struct {
	FileDescriptors []protoreflect.FileDescriptor
	Query           *ast.Definition
	Mutation        *ast.Definition
	Subscription    *ast.Definition
	Directives      map[string]*ast.DirectiveDefinition
	Types           map[string]*ast.Definition
}

func New() *SchemaDescriptor {
	return &SchemaDescriptor{
		Query: &ast.Definition{
			Kind:   ast.Object,
			Name:   "Query",
			Fields: []*ast.FieldDefinition{},
		},
		Mutation: &ast.Definition{
			Kind:   ast.Object,
			Name:   "Mutation",
			Fields: []*ast.FieldDefinition{},
		},
		Subscription: &ast.Definition{
			Kind:   ast.Object,
			Name:   "Subscription",
			Fields: []*ast.FieldDefinition{},
		},
		Types: make(map[string]*ast.Definition),
	}
}

func GraphqlMethodOptions(opts proto.Message) *graphqlv1.Rpc {
	if opts != nil {
		v := proto.GetExtension(opts, graphqlv1.E_Rpc)
		if v != nil {
			return v.(*graphqlv1.Rpc)
		}
	}
	return nil
}

// CreateObjects 创建 GraphQL 对象类型定义
func (s *SchemaDescriptor) CreateObjects(msgDesc protoreflect.MessageDescriptor, isInput bool) (*ast.Definition, error) {
	typeName := string(msgDesc.Name())
	if isInput {
		typeName += "Input"
	}

	if def, exists := s.Types[typeName]; exists {
		return def, nil
	}

	definition := &ast.Definition{
		Kind:   ast.Object,
		Name:   typeName,
		Fields: []*ast.FieldDefinition{},
	}

	for i := 0; i < msgDesc.Fields().Len(); i++ {
		field := msgDesc.Fields().Get(i)
		fieldDef := &ast.FieldDefinition{
			Name: string(field.Name()),
			Type: &ast.Type{
				NamedType: getGraphQLType(field),
			},
		}
		definition.Fields = append(definition.Fields, fieldDef)
	}

	s.Types[typeName] = definition
	return definition, nil
}

// getGraphQLType 将 proto 字段类型转换为 GraphQL 字段类型
func getGraphQLType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return "Boolean"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "Int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "Int"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "Int"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "Int"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "Float"
	case protoreflect.StringKind:
		return "String"
	case protoreflect.BytesKind:
		return "String"
	case protoreflect.EnumKind:
		return strings.ReplaceAll(string(field.Enum().FullName()), ".", "")
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return strings.ReplaceAll(string(field.Message().FullName()), ".", "")
	default:
		return "String"
	}
}

func (s *SchemaDescriptor) addMethod(def *ast.Definition, svc protoreflect.ServiceDescriptor, rpc protoreflect.MethodDescriptor, in, out *ast.Definition) {
	field := &ast.FieldDefinition{
		Name: string(rpc.Name()),
		Type: &ast.Type{
			NamedType: out.Name,
		},
		Arguments: []*ast.ArgumentDefinition{
			{
				Name: "input",
				Type: &ast.Type{
					NamedType: in.Name,
				},
			},
		},
	}
	def.Fields = append(def.Fields, field)
}

func (schema *SchemaDescriptor) GenerateFile(generateUnboundMethods bool, file protoreflect.FileDescriptor) error {
	for i := 0; i < file.Services().Len(); i++ {
		svc := file.Services().Get(i)
		for j := 0; j < svc.Methods().Len(); j++ {
			rpc := svc.Methods().Get(j)
			rpcOpts := GraphqlMethodOptions(rpc.Options())

			in, err := schema.CreateObjects(rpc.Input(), true)
			if err != nil {
				return err
			}

			out, err := schema.CreateObjects(rpc.Output(), false)
			if err != nil {
				return err
			}

			switch rpcOpts.GetPattern().(type) {
			case *graphqlv1.Rpc_Query:
				schema.addMethod(schema.Query, svc, rpc, in, out)
			default:
				schema.addMethod(schema.Mutation, svc, rpc, in, out)
			}
		}
	}

	return nil
}

func (s *SchemaDescriptor) AsGraphQL() *ast.Schema {
	queryDef := *s.Query
	mutationDef := *s.Mutation
	subscriptionsDef := *s.Subscription
	schema := &ast.Schema{Types: map[string]*ast.Definition{}, Directives: s.Directives}
	schema.Query = &queryDef
	schema.Types["Query"] = &queryDef

	schema.Mutation = &mutationDef
	schema.Types["Mutation"] = &mutationDef
	schema.Subscription = &subscriptionsDef
	schema.Types["Subscription"] = &subscriptionsDef

	schema.Types = s.Types

	return schema
}
