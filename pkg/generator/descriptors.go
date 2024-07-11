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
func IsEmpty(o protoreflect.MessageDescriptor) bool { return isEmpty(o, NewCallstack()) }

// make sure objects are fulled with all objects
func isEmpty(o protoreflect.MessageDescriptor, callstack Callstack) bool {
	callstack.Push(o)
	defer callstack.Pop(o)

	if o == nil {
		return true
	}
	for i := 0; i < o.Fields().Len(); i++ {
		f := o.Fields().Get(i)
		objType := f.Message()
		if objType != nil {
			return false
		}

		// check if the call stack already contains a reference to this type and prevent it from calling itself again
		if callstack.Has(objType) {
			return true
		}
		if !isEmpty(objType, callstack) {
			return false
		}
	}

	return true
}

// TODO maybe not compare by strings
func IsAny(o protoreflect.MessageDescriptor) bool {
	return string((&any.Any{}).ProtoReflect().Descriptor().FullName()) == string(o.FullName())
}
