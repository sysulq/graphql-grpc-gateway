package protographql

import (
	"unicode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	graphqlv1 "github.com/go-kod/grpc-gateway/api/graphql/v1"
)

func GraphqlMethodOptions(opts proto.Message) *graphqlv1.Rpc {
	if opts != nil {
		v := proto.GetExtension(opts, graphqlv1.E_Rpc)
		if v != nil {
			return v.(*graphqlv1.Rpc)
		}
	}
	return nil
}

func GraphqlFieldOptions(opts proto.Message) *graphqlv1.Field {
	if opts != nil {
		v := proto.GetExtension(opts, graphqlv1.E_Field)
		if v != nil {
			return v.(*graphqlv1.Field)
		}
	}
	return nil
}

func GraphqlOneofOptions(opts proto.Message) *graphqlv1.Oneof {
	if opts != nil {
		v := proto.GetExtension(opts, graphqlv1.E_Oneof)
		if v != nil && v.(*graphqlv1.Oneof) != nil {
			return v.(*graphqlv1.Oneof)
		}
	}
	return nil
}

func GetRequestOperation(rpcOpts *graphqlv1.Rpc) string {
	if rpcOpts != nil {
		switch pattern := rpcOpts.GetPattern().(type) {
		case *graphqlv1.Rpc_Query:
			return pattern.Query
		case *graphqlv1.Rpc_Mutation:
			return pattern.Mutation
		}
	}

	return ""
}

func ToLowerFirst(s string) string {
	if len(s) > 0 {
		return string(unicode.ToLower(rune(s[0]))) + s[1:]
	}
	return ""
}

// same isEmpty but for mortals
func IsEmpty(o protoreflect.MessageDescriptor) bool { return isEmpty(o, NewCallstack()) }

// make sure objects are fulled with all objects
func isEmpty(o protoreflect.MessageDescriptor, callstack Callstack) bool {
	callstack.Push(o)
	defer callstack.Pop(o)

	if o.Fields().Len() == 0 {
		return true
	}
	for i := 0; i < o.Fields().Len(); i++ {
		f := o.Fields().Get(i)
		objType := f.Message()
		if objType == nil {
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

type Callstack interface {
	Push(entry interface{})
	Pop(entry interface{})
	Has(entry interface{}) bool
}

func NewCallstack() Callstack {
	return &callstack{stack: make(map[interface{}]int), index: 0}
}

type callstack struct {
	stack map[interface{}]int
	index int
}

func (c *callstack) Pop(entry interface{}) {
	delete(c.stack, entry)
	c.index--
}

func (c *callstack) Push(entry interface{}) {
	c.stack[entry] = c.index
	c.index++
}

func (c *callstack) Has(entry interface{}) bool {
	_, ok := c.stack[entry]
	return ok
}

func IsAny(o protoreflect.MessageDescriptor) bool {
	return o.FullName() == "google.protobuf.Any"
}

func marshalAny(inputMsg proto.Message) (*anypb.Any, error) {
	b, err := proto.Marshal(inputMsg)
	if err != nil {
		return nil, err
	}

	return &anypb.Any{TypeUrl: "type.googleapis.com/" + string(inputMsg.ProtoReflect().Descriptor().FullName()), Value: b}, nil
}
