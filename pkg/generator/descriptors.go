package generator

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
	any "google.golang.org/protobuf/types/known/anypb"
)

type ObjectDescriptor struct {
	*ast.Definition
	protoreflect.Descriptor

	types      []*ObjectDescriptor
	fields     []*FieldDescriptor
	fieldNames map[string]*FieldDescriptor
}

func (o *ObjectDescriptor) AsGraphql() *ast.Definition {
	return o.Definition
}

func (o *ObjectDescriptor) uniqueName(f protoreflect.FieldDescriptor) string {
	return strings.Title(string(f.Name()))
}

func (o *ObjectDescriptor) IsInput() bool {
	return o.Kind == ast.InputObject
}

func (o *ObjectDescriptor) GetFields() []*FieldDescriptor {
	return o.fields
}

func (o *ObjectDescriptor) GetTypes() []*ObjectDescriptor {
	return o.types
}

func (o *ObjectDescriptor) IsMessage() bool {
	_, ok := o.Descriptor.(protoreflect.MessageDescriptor)
	return ok
}

// same isEmpty but for mortals
func IsEmpty(o protoreflect.MessageDescriptor) bool { return isEmpty(o) }

// make sure objects are fulled with all objects
func isEmpty(o protoreflect.MessageDescriptor) bool {
	if o == nil {
		return true
	}

	if o.FullName() == "google.protobuf.Empty" {
		return true
	}

	if o.Fields().Len() == 0 {
		return true
	}

	return false
}

// TODO maybe not compare by strings
func IsAny(o protoreflect.MessageDescriptor) bool {
	// fmt.Println(string((&any.Any{}).ProtoReflect().Descriptor().FullName()), string(o.FullName()))
	return string((&any.Any{}).ProtoReflect().Descriptor().FullName()) == string(o.FullName())
}
