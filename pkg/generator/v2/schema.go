package v2

import (
	"strings"

	graphqlv1 "github.com/sysulq/graphql-grpc-gateway/api/graphql/v1"
	"github.com/sysulq/graphql-grpc-gateway/pkg/generator"
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
	ProcessedTypes  map[string]struct{}
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
		Types:          make(map[string]*ast.Definition),
		ProcessedTypes: make(map[string]struct{}),
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
	if msgDesc == nil {
		return nil, nil
	}

	typeName := s.msgFullName(msgDesc)
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

	if isInput {
		definition.Kind = ast.InputObject
	}
	if generator.IsEmptyV2(msgDesc) {
		return definition, nil
	}

	s.Types[typeName] = definition

	for i := 0; i < msgDesc.Fields().Len(); i++ {
		field := msgDesc.Fields().Get(i)

		if field.Kind() == protoreflect.MessageKind && generator.IsEmptyV2(field.Message()) {
			continue
		}

		fieldType, err := s.getGraphQLFieldType(field, isInput)
		if err != nil {
			return nil, err
		}
		fieldDef := &ast.FieldDefinition{
			Name: field.JSONName(),
			Type: fieldType,
		}
		definition.Fields = append(definition.Fields, fieldDef)
	}

	// 处理 oneof 字段
	for i := 0; i < msgDesc.Oneofs().Len(); i++ {
		oneof := msgDesc.Oneofs().Get(i)

		if !isInput {
			fieldDef, err := s.createUnion(oneof)
			if err != nil {
				return nil, err
			}
			definition.Fields = append(definition.Fields, fieldDef)
		}

	}

	return definition, nil
}

// createUnion 处理 oneof 字段，创建 GraphQL 联合类型和关联的对象结构
func (s *SchemaDescriptor) createUnion(oneof protoreflect.OneofDescriptor) (*ast.FieldDefinition, error) {
	var types []string
	for i := 0; i < oneof.Fields().Len(); i++ {
		choice := oneof.Fields().Get(i)
		choiceName := s.uniqueName(choice)
		objDef, err := s.createObjectFromChoice(choice)
		if err != nil {
			return nil, err
		}

		objDef.Name = choiceName
		s.Types[choiceName] = objDef
		types = append(types, choiceName)
	}

	unionName := s.uniqueName(oneof)
	unionDef := &ast.Definition{
		Kind:  ast.Union,
		Name:  unionName,
		Types: types,
	}

	s.Types[unionName] = unionDef

	fieldName := string(oneof.Name())
	return &ast.FieldDefinition{
			Name: fieldName,
			Type: &ast.Type{NamedType: unionName},
		},
		nil
}

// createObjectFromChoice 生成 oneof 中的选择项对象
func (s *SchemaDescriptor) createObjectFromChoice(choice protoreflect.FieldDescriptor) (*ast.Definition, error) {
	fieldType, err := s.getGraphQLFieldType(choice, false)
	if err != nil {
		return nil, err
	}

	fields := []*ast.FieldDefinition{
		{
			Name: string(choice.Name()),
			Type: fieldType,
		},
	}

	return &ast.Definition{
		Kind:   ast.Object,
		Name:   s.uniqueName(choice),
		Fields: fields,
	}, nil
}

// getGraphQLFieldType 将 proto 字段类型转换为 GraphQL 字段类型
func (s *SchemaDescriptor) getGraphQLFieldType(field protoreflect.FieldDescriptor, isInput bool) (*ast.Type, error) {
	var astType *ast.Type

	switch field.Kind() {
	case protoreflect.BoolKind:
		astType = &ast.Type{NamedType: "Boolean"}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		astType = &ast.Type{NamedType: "Int"}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		astType = &ast.Type{NamedType: "Int"}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		astType = &ast.Type{NamedType: "Int"}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		astType = &ast.Type{NamedType: "Int"}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		astType = &ast.Type{NamedType: "Float"}
	case protoreflect.StringKind:
		astType = &ast.Type{NamedType: "String"}
	case protoreflect.BytesKind:
		astType = &ast.Type{NamedType: "String"}
	case protoreflect.EnumKind:
		astType = s.getGraphQLEnumType(field.Enum())
	case protoreflect.MessageKind, protoreflect.GroupKind:
		nestedType, err := s.CreateObjects(field.Message(), isInput)
		if err != nil {
			return nil, err
		}

		astType = &ast.Type{NamedType: nestedType.Name}
	default:
		return &ast.Type{NamedType: "String"}, nil
	}

	if field.IsList() {
		astType.NonNull = true
		astType = ast.ListType(astType, &ast.Position{})
	} else if field.IsMap() {
		astType.NonNull = true
		astType = ast.ListType(astType, &ast.Position{})
	}

	return astType, nil
}

// getGraphQLEnumType 将 proto 枚举类型转换为 GraphQL 枚举类型
func (s *SchemaDescriptor) getGraphQLEnumType(enumDesc protoreflect.EnumDescriptor) *ast.Type {
	typeName := strings.ReplaceAll(strings.TrimPrefix(string(enumDesc.FullName()), string(enumDesc.ParentFile().Package())+"."), ".", "_")

	if _, exists := s.Types[typeName]; !exists {
		enumDef := &ast.Definition{
			Kind:       ast.Enum,
			Name:       typeName,
			EnumValues: make([]*ast.EnumValueDefinition, enumDesc.Values().Len()),
		}

		for i := 0; i < enumDesc.Values().Len(); i++ {
			enumValue := enumDesc.Values().Get(i)
			enumDef.EnumValues[i] = &ast.EnumValueDefinition{
				Name: string(enumValue.Name()),
			}
		}

		s.Types[typeName] = enumDef
	}

	return &ast.Type{NamedType: typeName}
}

func (s *SchemaDescriptor) addMethod(def *ast.Definition, svc protoreflect.ServiceDescriptor, rpc protoreflect.MethodDescriptor, in, out *ast.Definition) {
	field := &ast.FieldDefinition{
		Name: generator.ToLowerFirst(string(rpc.Parent().Name() + rpc.Name())),
	}

	field.Type = ast.NamedType("Boolean", &ast.Position{})
	if rpc.Output() != nil && !generator.IsEmptyV2(rpc.Output()) {
		field.Type = &ast.Type{
			NamedType: out.Name,
		}
	}

	if rpc.Input() != nil && !generator.IsEmptyV2(rpc.Input()) {
		field.Arguments = []*ast.ArgumentDefinition{
			{
				Name: "in",
				Type: &ast.Type{
					NamedType: in.Name,
				},
			},
		}
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

			if !rpc.IsStreamingClient() && rpc.IsStreamingServer() {
				schema.addMethod(schema.Subscription, svc, rpc, in, out)
				return nil
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
	schema := &ast.Schema{Types: map[string]*ast.Definition{}, Directives: s.Directives}
	schema.Types = s.Types

	if len(s.Query.Fields) > 0 {
		schema.Query = s.Query
		schema.Types["Query"] = s.Query
	}

	if len(s.Mutation.Fields) > 0 {
		schema.Mutation = s.Mutation
		schema.Types["Mutation"] = s.Mutation
	}

	if len(s.Subscription.Fields) > 0 {
		schema.Subscription = s.Subscription
		schema.Types["Subscription"] = s.Subscription
	}

	return schema
}

// uniqueName 生成唯一名称
func (s *SchemaDescriptor) uniqueName(desc protoreflect.Descriptor) string {
	return strings.ReplaceAll(string(desc.FullName()), ".", "_")
}

func (s *SchemaDescriptor) msgFullName(msg protoreflect.MessageDescriptor) string {
	name := strings.ReplaceAll(string(msg.FullName()), ".", "")
	return strings.TrimPrefix(name, string(msg.ParentFile().Package()))
}
