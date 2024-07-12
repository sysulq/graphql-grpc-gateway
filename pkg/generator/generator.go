package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jhump/protoreflect/v2/protowrap"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	descriptor "google.golang.org/protobuf/types/descriptorpb"

	gqlpb "github.com/sysulq/graphql-gateway/api/graphql/v1"
)

const (
	fieldPrefix        = "Field"
	inputSuffix        = "Input"
	typeSep            = "_"
	packageSep         = "."
	anyTypeDescription = "Any is any json type"
	scalarBytes        = "Bytes"
	goFieldDirective   = "goField"

	DefaultExtension = "graphql"
)

func NewSchemas(descs []protoreflect.FileDescriptor, mergeSchemas, genServiceDesc bool, plugin *protogen.Plugin) (schemas SchemaDescriptorList, err error) {
	var files []*descriptor.FileDescriptorProto
	for _, d := range descs {
		files = append(files, protowrap.ProtoFromFileDescriptor(d))
	}
	var goref GoRef
	if plugin != nil {
		goref, err = NewGoRef(plugin)
		if err != nil {
			return nil, err
		}
	}

	if mergeSchemas {
		schema := NewSchemaDescriptor(genServiceDesc, goref)
		for _, file := range descs {
			err := generateFile(file, schema)
			if err != nil {
				return nil, err
			}
		}

		return []*SchemaDescriptor{schema}, nil
	}

	for _, file := range descs {
		schema := NewSchemaDescriptor(genServiceDesc, goref)
		err := generateFile(file, schema)
		if err != nil {
			return nil, err
		}

		schemas = append(schemas, schema)
	}

	return
}

func generateFile(file protoreflect.FileDescriptor, schema *SchemaDescriptor) error {
	schema.FileDescriptors = append(schema.FileDescriptors, file)

	for i := 0; i < file.Services().Len(); i++ {
		svc := file.Services().Get(i)
		svcOpts := GraphqlServiceOptions(svc.Options())
		if svcOpts != nil && svcOpts.Ignore {
			continue
		}
		for j := 0; j < svc.Methods().Len(); j++ {
			rpc := svc.Methods().Get(j)
			rpcOpts := GraphqlMethodOptions(rpc.Options())
			if rpcOpts != nil && rpcOpts.Ignore {
				continue
			}
			in, err := schema.CreateObjects(rpc.Input(), true)
			if err != nil {
				return err
			}

			out, err := schema.CreateObjects(rpc.Output(), false)
			if err != nil {
				return err
			}

			if rpc.IsStreamingServer() && rpc.IsStreamingClient() {
				schema.GetMutation().addMethod(svc, rpc, in, out)
			}

			if rpc.IsStreamingServer() {
				schema.GetSubscription().addMethod(svc, rpc, in, out)
			} else {
				switch GetRequestType(rpcOpts, svcOpts) {
				case gqlpb.Type_QUERY:
					schema.GetQuery().addMethod(svc, rpc, in, out)
				default:
					fmt.Println(rpc.FullName())
					schema.GetMutation().addMethod(svc, rpc, in, out)
				}
			}
		}
	}

	return nil
}

type SchemaDescriptorList []*SchemaDescriptor

func (s SchemaDescriptorList) AsGraphql() (astSchema []*ast.Schema) {
	for _, ss := range s {
		astSchema = append(astSchema, ss.AsGraphql())
	}
	return
}

func (s SchemaDescriptorList) GetForDescriptor(file *protogen.File) *SchemaDescriptor {
	for _, schema := range s {
		for _, d := range schema.FileDescriptors {
			if protowrap.ProtoFromFileDescriptor(d) == file.Proto {
				return schema
			}
		}
	}
	return nil
}

func NewSchemaDescriptor(genServiceDesc bool, goref GoRef) *SchemaDescriptor {
	sd := &SchemaDescriptor{
		Directives:                 map[string]*ast.DirectiveDefinition{},
		reservedNames:              map[string]protoreflect.Descriptor{},
		createdObjects:             map[createdObjectKey]*ObjectDescriptor{},
		generateServiceDescriptors: genServiceDesc,
		goRef:                      goref,
	}
	for _, name := range graphqlReservedNames {
		sd.reservedNames[name] = nil
	}
	return sd
}

type SchemaDescriptor struct {
	Directives      map[string]*ast.DirectiveDefinition
	FileDescriptors []protoreflect.FileDescriptor

	files []protoreflect.FileDescriptor

	query        *RootDefinition
	mutation     *RootDefinition
	subscription *RootDefinition

	objects []*ObjectDescriptor

	reservedNames  map[string]protoreflect.Descriptor
	createdObjects map[createdObjectKey]*ObjectDescriptor

	generateServiceDescriptors bool

	goRef GoRef
}

type createdObjectKey struct {
	desc  protoreflect.FullName
	input bool
}

func (s *SchemaDescriptor) AsGraphql() *ast.Schema {
	queryDef := *s.GetQuery().Definition
	mutationDef := *s.GetMutation().Definition
	subscriptionsDef := *s.GetSubscription().Definition
	schema := &ast.Schema{Types: map[string]*ast.Definition{}, Directives: s.Directives}
	schema.Query = &queryDef
	schema.Types["Query"] = &queryDef
	if s.query.methods == nil {
		schema.Query.Fields = append(schema.Query.Fields, &ast.FieldDefinition{
			Name: "dummy",
			Type: ast.NamedType("Boolean", &ast.Position{}),
		})
	}
	if s.mutation.methods != nil {
		schema.Mutation = &mutationDef
		schema.Types["Mutation"] = &mutationDef
	}
	if s.subscription.methods != nil {
		schema.Subscription = &subscriptionsDef
		schema.Types["Subscription"] = &subscriptionsDef
	}

	for _, o := range s.objects {
		def := o.AsGraphql()
		schema.Types[def.Name] = def
	}
	return schema
}

func (s *SchemaDescriptor) Objects() []*ObjectDescriptor {
	return s.objects
}

func (s *SchemaDescriptor) GetMutation() *RootDefinition {
	if s.mutation == nil {
		s.mutation = NewRootDefinition(Mutation, s)
	}
	return s.mutation
}

func (s *SchemaDescriptor) GetSubscription() *RootDefinition {
	if s.subscription == nil {
		s.subscription = NewRootDefinition(Subscription, s)
	}
	return s.subscription
}

func (s *SchemaDescriptor) GetQuery() *RootDefinition {
	if s.query == nil {
		s.query = NewRootDefinition(Query, s)
	}

	return s.query
}

// make name be unique
// just create a map and register every name
func (s *SchemaDescriptor) uniqueName(d protoreflect.Descriptor, input bool) (name string) {
	var collisionPrefix string
	var suffix string
	if _, ok := d.(protoreflect.MessageDescriptor); input && ok {
		suffix = inputSuffix
	}
	name = strings.Title(CamelCaseSlice(strings.Split(strings.TrimPrefix(string(d.FullName()), string(d.ParentFile().Package())+packageSep), packageSep)) + suffix)

	if _, ok := d.(protoreflect.FieldDescriptor); ok {
		collisionPrefix = fieldPrefix
		name = CamelCaseSlice(strings.Split(strings.Trim(string(d.Parent().Name())+packageSep+strings.Title(string(d.Name())), packageSep), packageSep))
	} else {
		collisionPrefix = CamelCaseSlice(strings.Split(string(d.ParentFile().Package()), packageSep))
	}

	originalName := name
	for uniqueSuffix := 0; ; uniqueSuffix++ {
		d2, ok := s.reservedNames[name]
		if !ok {
			break
		}
		if d2 == d {
			return name
		}
		if uniqueSuffix == 0 {
			name = collisionPrefix + typeSep + originalName
			continue
		}
		name = collisionPrefix + typeSep + originalName + strconv.Itoa(uniqueSuffix)
	}

	s.reservedNames[name] = d
	return
}

func (s *SchemaDescriptor) CreateObjects(d protoreflect.Descriptor, input bool) (obj *ObjectDescriptor, err error) {
	// the case if trying to resolve a primitive as a object. In this case we just return nil
	if d == nil || d.Name() == "FieldMask" {
		return
	}

	if obj, ok := s.createdObjects[createdObjectKey{d.FullName(), input}]; ok {
		return obj, nil
	}

	obj = &ObjectDescriptor{
		Definition: &ast.Definition{
			Description: getDescription(d),
			Name:        s.uniqueName(d, input),
			Position:    &ast.Position{},
		},
		Descriptor: d,
	}

	s.createdObjects[createdObjectKey{d.FullName(), input}] = obj

	switch dd := d.(type) {
	case protoreflect.MessageDescriptor:
		fmt.Println("----", dd.FullName())
		if IsEmpty(dd) {
			return obj, nil
		}

		if IsAny(dd) {
			// TODO find a better way to handle any types
			delete(s.createdObjects, createdObjectKey{d.FullName(), input})
			any := s.createScalar(s.uniqueName(dd, false), anyTypeDescription)
			return any, nil
		}

		kind := ast.Object
		if input {
			kind = ast.InputObject
		}
		fields := FieldDescriptorList{}
		outputOneofRegistrar := map[protoreflect.OneofDescriptor]struct{}{}

		for i := 0; i < dd.Fields().Len(); i++ {

			df := dd.Fields().Get(i)

			fieldOpts := GraphqlFieldOptions(df.Options())
			if fieldOpts != nil && fieldOpts.Ignore {
				continue
			}
			var fieldDirective []*ast.Directive
			if IsEmpty(dd) {
				continue
			}

			// Internally `optional` fields are represented as a oneof, and as such should be skipped.
			if oneof := df.ContainingOneof(); oneof != nil && !protowrap.ProtoFromFieldDescriptor(df).GetProto3Optional() {
				opts := GraphqlOneofOptions(oneof.Options())
				if opts.GetIgnore() {
					continue
				}
				if !input {
					if _, ok := outputOneofRegistrar[oneof]; ok {
						continue
					}
					outputOneofRegistrar[oneof] = struct{}{}
					field, err := s.createUnion(oneof)
					if err != nil {
						return nil, err
					}
					fields = append(fields, field)
					continue
				}

				// create oneofs as directives for input objects
				directive := &ast.DirectiveDefinition{
					Description: getDescription(oneof),
					Name:        s.uniqueName(oneof, input),
					Locations:   []ast.DirectiveLocation{ast.LocationInputFieldDefinition},
					Position:    &ast.Position{Src: &ast.Source{}},
				}
				s.Directives[directive.Name] = directive
				fieldDirective = append(fieldDirective, &ast.Directive{
					Name:     directive.Name,
					Position: &ast.Position{Src: &ast.Source{}},
					// ParentDefinition: obj.Definition, TODO
					Definition: directive,
					Location:   ast.LocationInputFieldDefinition,
				})
			}

			fieldObj, err := s.CreateObjects(resolveFieldType(df), input)
			if err != nil {
				return nil, err
			}
			if fieldObj == nil && df.Message() != nil && !df.IsMap() {
				continue
			}
			if df.Message() != nil && df.Message().FullName() == "google.protobuf.FieldMask" {
				continue
			}
			f, err := s.createField(df, fieldObj)
			if err != nil {
				return nil, err
			}
			f.Directives = append(f.Directives, fieldDirective...)
			fields = append(fields, f)
		}

		obj.Definition.Fields = fields.AsGraphql()
		obj.Definition.Kind = kind
		obj.fields = fields
	case protoreflect.EnumDescriptor:
		obj.Definition.Kind = ast.Enum
		vv := make([]protoreflect.EnumValueDescriptor, 0, dd.Values().Len())
		for i := 0; i < dd.Values().Len(); i++ {
			vv = append(vv, dd.Values().Get(i))
		}
		obj.Definition.EnumValues = enumValues(vv)
	default:
		panic(fmt.Sprintf("received unexpected value %v of type %T", dd, dd))
	}

	s.objects = append(s.objects, obj)
	return obj, nil
}

func resolveFieldType(field protoreflect.FieldDescriptor) protoreflect.Descriptor {
	msgType := field.Message()
	enumType := field.Enum()
	if msgType != nil {
		return msgType
	}
	if enumType != nil {
		return enumType
	}
	return nil
}

func enumValues(evals []protoreflect.EnumValueDescriptor) (vlist ast.EnumValueList) {
	for _, eval := range evals {
		vlist = append(vlist, &ast.EnumValueDefinition{
			Description: getDescription(eval),
			Name:        string(eval.Name()),
			Position:    &ast.Position{},
		})
	}

	return vlist
}

type FieldDescriptorList []*FieldDescriptor

func (fl FieldDescriptorList) AsGraphql() (dl []*ast.FieldDefinition) {
	for _, f := range fl {
		dl = append(dl, f.FieldDefinition)
	}
	return dl
}

type FieldDescriptor struct {
	*ast.FieldDefinition
	protoreflect.FieldDescriptor

	typ *ObjectDescriptor
}

func (f *FieldDescriptor) GetType() *ObjectDescriptor {
	return f.typ
}

type MethodDescriptor struct {
	protoreflect.ServiceDescriptor
	protoreflect.MethodDescriptor

	*ast.FieldDefinition

	input  *ObjectDescriptor
	output *ObjectDescriptor
}

func (m *MethodDescriptor) AsGraphql() *ast.FieldDefinition {
	return m.FieldDefinition
}

func (m *MethodDescriptor) GetInput() *ObjectDescriptor {
	return m.input
}

func (m *MethodDescriptor) GetOutput() *ObjectDescriptor {
	return m.output
}

type RootDefinition struct {
	*ast.Definition

	Parent *SchemaDescriptor

	methods       []*MethodDescriptor
	reservedNames map[string]ServiceAndMethod
}

type ServiceAndMethod struct {
	svc *descriptor.ServiceDescriptorProto
	rpc *descriptor.MethodDescriptorProto
}

func (r *RootDefinition) UniqueName(svc *descriptor.ServiceDescriptorProto, rpc *descriptor.MethodDescriptorProto) (name string) {
	rpcOpts := GraphqlMethodOptions(rpc.GetOptions())
	svcOpts := GraphqlServiceOptions(svc.GetOptions())
	if rpcOpts != nil && rpcOpts.Name != "" {
		name = rpcOpts.Name
	} else if svcOpts != nil && svcOpts.Name != "" {
		name = svcOpts.Name + strings.Title(rpc.GetName())
	} else {
		name = ToLowerFirst(svc.GetName()) + strings.Title(rpc.GetName())
	}

	originalName := name
	for uniqueSuffix := 0; ; uniqueSuffix++ {
		snm, ok := r.reservedNames[name]
		if !ok {
			break
		}
		if svc == snm.svc && snm.rpc == rpc {
			return name
		}
		name = originalName + strconv.Itoa(uniqueSuffix)
	}

	r.reservedNames[name] = ServiceAndMethod{svc, rpc}
	return
}

func (r *RootDefinition) Methods() []*MethodDescriptor {
	return r.methods
}

func (r *RootDefinition) addMethod(svc protoreflect.ServiceDescriptor, rpc protoreflect.MethodDescriptor, in, out *ObjectDescriptor) {
	var args ast.ArgumentDefinitionList

	if in != nil && (in.Descriptor != nil && !IsEmpty(in.Descriptor.(protoreflect.MessageDescriptor)) || in.Definition.Kind == ast.Scalar) {
		args = append(args, &ast.ArgumentDefinition{
			Name:     "in",
			Type:     ast.NamedType(in.Definition.Name, &ast.Position{}),
			Position: &ast.Position{},
		})
	}
	objType := ast.NamedType("Boolean", &ast.Position{})
	if out != nil && (out.Descriptor != nil && !IsEmpty(out.Descriptor.(protoreflect.MessageDescriptor)) || in.Definition.Kind == ast.Scalar) {
		objType = ast.NamedType(out.Definition.Name, &ast.Position{})
	}

	svcDir := &ast.DirectiveDefinition{
		Description: getDescription(svc),
		Name:        string(svc.Name()),
		Locations:   []ast.DirectiveLocation{ast.LocationFieldDefinition},
		Position:    &ast.Position{Src: &ast.Source{}},
	}
	r.Parent.Directives[svcDir.Name] = svcDir

	m := &MethodDescriptor{
		ServiceDescriptor: svc,
		MethodDescriptor:  rpc,
		FieldDefinition: &ast.FieldDefinition{
			Description: getDescription(rpc),
			Name:        r.UniqueName(protowrap.ProtoFromServiceDescriptor(svc), protowrap.ProtoFromMethodDescriptor(rpc)),
			Arguments:   args,
			Type:        objType,
			Position:    &ast.Position{},
		},
		input:  in,
		output: out,
	}
	if r.Parent.generateServiceDescriptors {
		m.Directives = []*ast.Directive{{
			Name:       svcDir.Name,
			Position:   &ast.Position{},
			Definition: svcDir,
			Location:   svcDir.Locations[0],
		}}
	}

	r.methods = append(r.methods, m)
	// TODO maybe not do it here?
	r.Definition.Fields = append(r.Definition.Fields, m.FieldDefinition)
}

type rootName string

const (
	Mutation     rootName = "Mutation"
	Query        rootName = "Query"
	Subscription rootName = "Subscription"
)

func NewRootDefinition(name rootName, parent *SchemaDescriptor) *RootDefinition {
	return &RootDefinition{Definition: &ast.Definition{
		Kind:     ast.Object,
		Name:     string(name),
		Position: &ast.Position{},
	}, Parent: parent, reservedNames: map[string]ServiceAndMethod{}}
}

func getDescription(descs ...protoreflect.Descriptor) string {
	var description []string
	for _, d := range descs {
		info := d.ParentFile().SourceLocations()
		if info.Len() == 0 {
			continue
		}
		for i := 0; i < info.Len(); i++ {
			loc := info.Get(i)

			if loc.LeadingComments != "" {
				description = append(description, loc.LeadingComments)
			}

			if loc.TrailingComments != "" {
				description = append(description, loc.TrailingComments)
			}
		}
	}

	return strings.Join(description, "\n")
}

func (s *SchemaDescriptor) createField(field protoreflect.FieldDescriptor, obj *ObjectDescriptor) (_ *FieldDescriptor, err error) {
	fieldAst := &ast.FieldDefinition{
		Description: getDescription(field),
		Name:        ToLowerFirst(CamelCase(string(field.Name()))),
		Type:        &ast.Type{Position: &ast.Position{}},
		Position:    &ast.Position{},
	}
	fieldOpts := GraphqlFieldOptions(field.Options())
	if fieldOpts != nil {
		if fieldOpts.Name != "" {
			fieldAst.Name = fieldOpts.Name
		}
		directive := &ast.DirectiveDefinition{
			Name: goFieldDirective,
			Arguments: []*ast.ArgumentDefinition{{
				Name:     "forceResolver",
				Type:     ast.NamedType("Boolean", &ast.Position{}),
				Position: &ast.Position{},
			}, {
				Name:     "name",
				Type:     ast.NamedType("String", &ast.Position{}),
				Position: &ast.Position{},
			}},
			Locations: []ast.DirectiveLocation{ast.LocationInputFieldDefinition, ast.LocationFieldDefinition},
			Position:  &ast.Position{Src: &ast.Source{}},
		}
		s.Directives[directive.Name] = directive
		if s.goRef != nil {
			fieldAst.Directives = []*ast.Directive{{
				Name: directive.Name,
				Arguments: []*ast.Argument{{
					Name: "name",
					Value: &ast.Value{
						Raw:      s.goRef.FindGoField(string(field.FullName())).GoName,
						Kind:     ast.StringValue,
						Position: &ast.Position{},
					},
					Position: &ast.Position{},
				}},
				Position: &ast.Position{},
				// ParentDefinition: nil, TODO
				Definition: directive,
			}}
		}
	}
	switch field.Kind() {
	case protoreflect.DoubleKind,
		protoreflect.FloatKind:
		fieldAst.Type.NamedType = ScalarFloat

	case protoreflect.BytesKind:
		scalar := s.createScalar(scalarBytes, "")
		fieldAst.Type.NamedType = scalar.Definition.Name

	case protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind:
		fieldAst.Type.NamedType = ScalarInt

	case protoreflect.BoolKind:
		fieldAst.Type.NamedType = ScalarBoolean

	case protoreflect.StringKind:
		fieldAst.Type.NamedType = ScalarString

	case protoreflect.GroupKind:
		return nil, fmt.Errorf("proto2 groups are not supported please use proto3 syntax")

	case protoreflect.EnumKind:
		fieldAst.Type.NamedType = obj.Definition.Name

	case protoreflect.MessageKind:
		fieldAst.Type.NamedType = obj.Definition.Name

	default:
		panic("unknown proto field type")
	}

	if isRepeated(field) {
		fieldAst.Type = ast.ListType(fieldAst.Type, &ast.Position{})
		fieldAst.Type.Elem.NonNull = true
	}
	if isRequired(field) {
		fieldAst.Type.NonNull = true
	}

	return &FieldDescriptor{
		FieldDefinition: fieldAst,
		FieldDescriptor: field,
		typ:             obj,
	}, nil
}

func (s *SchemaDescriptor) createScalar(name string, description string) *ObjectDescriptor {
	obj := &ObjectDescriptor{
		Definition: &ast.Definition{
			Kind:        ast.Scalar,
			Description: description,
			Name:        name,
			Position:    &ast.Position{},
		},
	}
	s.objects = append(s.objects, obj)
	return obj
}

func (s *SchemaDescriptor) createUnion(oneof protoreflect.OneofDescriptor) (*FieldDescriptor, error) {
	var types []string
	var objTypes []*ObjectDescriptor
	for i := 0; i < oneof.Fields().Len(); i++ {
		choice := oneof.Fields().Get(i)
		obj, err := s.CreateObjects(resolveFieldType(choice), false)
		if err != nil {
			return nil, err
		}
		f, err := s.createField(choice, obj)
		if err != nil {
			return nil, err
		}

		obj = &ObjectDescriptor{
			Definition: &ast.Definition{
				Kind:        ast.Object,
				Description: getDescription(f.FieldDescriptor),
				Name:        s.uniqueName(choice, false),
				Fields:      ast.FieldList{f.FieldDefinition},
				Position:    &ast.Position{},
			},
			Descriptor: f.FieldDescriptor,
			fields:     []*FieldDescriptor{f},
			fieldNames: map[string]*FieldDescriptor{},
		}
		s.objects = append(s.objects, obj)
		types = append(types, obj.Definition.Name)
		objTypes = append(objTypes, obj)
	}
	obj := &ObjectDescriptor{
		Definition: &ast.Definition{
			Kind:        ast.Union,
			Description: getDescription(oneof),
			Name:        s.uniqueName(oneof, false),
			Types:       types,
			Position:    &ast.Position{},
		},
		Descriptor: oneof,
		types:      objTypes,
	}
	s.objects = append(s.objects, obj)
	name := ToLowerFirst(CamelCase(string(oneof.Name())))
	opts := GraphqlOneofOptions(oneof.Options())
	if opts.GetName() != "" {
		name = opts.GetName()
	}
	return &FieldDescriptor{
		FieldDefinition: &ast.FieldDefinition{
			Description: getDescription(oneof),
			Name:        name,
			Type:        ast.NamedType(obj.Definition.Name, &ast.Position{}),
			Position:    &ast.Position{},
		},
		FieldDescriptor: nil,
		typ:             obj,
	}, nil
}

func isRepeated(field protoreflect.FieldDescriptor) bool {
	return field.Cardinality() == protoreflect.Repeated
}

func isRequired(field protoreflect.FieldDescriptor) bool {
	if v := GraphqlFieldOptions(field.Options()); v != nil {
		return v.GetRequired()
	}
	return false
}

const (
	ScalarInt     = "Int"
	ScalarFloat   = "Float"
	ScalarString  = "String"
	ScalarBoolean = "Boolean"
	ScalarID      = "ID"
)

var graphqlReservedNames = []string{"__Directive", "__Type", "__Field", "__EnumValue", "__InputValue", "__Schema", "Int", "Float", "String", "Boolean", "ID"}
