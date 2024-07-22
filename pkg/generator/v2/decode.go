package v2

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// DecodeProto2GraphQL 将 gRPC 返回的 proto.Message 转换为 GraphQL 数据
func (ins *Generator) DecodeProto2GraphQL(msg proto.Message, field *ast.Field) (interface{}, error) {
	message := msg.ProtoReflect()
	result := make(map[string]interface{})

	for _, selection := range field.SelectionSet {
		if f, ok := selection.(*ast.Field); ok {
			fieldName := f.Name
			fieldDesc := message.Descriptor().Fields().ByJSONName(fieldName)
			if fieldDesc == nil {
				return nil, fmt.Errorf("field %s not found in message", fieldName)
			}
			value := message.Get(fieldDesc)
			decodedValue, err := decodeValue(value, fieldDesc, f)
			if err != nil {
				return nil, err
			}
			result[fieldName] = decodedValue
		}
	}

	return result, nil
}

// decodeValue 将 protoreflect.Value 转换为适合 GraphQL 的值
func decodeValue(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	switch {
	case fieldDesc.IsList():
		return decodeList(value.List(), fieldDesc, field)
	case fieldDesc.IsMap():
		return decodeMap(value.Map(), fieldDesc, field)
	default:
		return decodeScalar(value, fieldDesc, field)
	}
}

func decodeScalar(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	switch fieldDesc.Kind() {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(value.Int()), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(value.Uint()), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int(), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint(), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return value.Bytes(), nil
	// case protoreflect.EnumKind:
	// 	return value.Enum().String(), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return decodeMessage(value.Message(), field)
	default:
		return nil, fmt.Errorf("unsupported field type: %v", fieldDesc.Kind())
	}
}

func decodeMessage(message protoreflect.Message, field *ast.Field) (interface{}, error) {
	result := make(map[string]interface{})
	for _, selection := range field.SelectionSet {
		if f, ok := selection.(*ast.Field); ok {
			fieldName := f.Name
			fieldDesc := message.Descriptor().Fields().ByJSONName(fieldName)
			if fieldDesc == nil {
				return nil, fmt.Errorf("field %s not found in message", fieldName)
			}
			value := message.Get(fieldDesc)
			decodedValue, err := decodeValue(value, fieldDesc, f)
			if err != nil {
				return nil, err
			}
			result[fieldName] = decodedValue
		}
	}
	return result, nil
}

func decodeList(list protoreflect.List, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	result := make([]interface{}, list.Len())
	for i := 0; i < list.Len(); i++ {
		value := list.Get(i)
		decodedValue, err := decodeValue(value, fieldDesc, field)
		if err != nil {
			return nil, err
		}
		result[i] = decodedValue
	}
	return result, nil
}

func decodeMap(m protoreflect.Map, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	result := make([]map[string]interface{}, 0)
	m.Range(func(mapKey protoreflect.MapKey, mapValue protoreflect.Value) bool {
		keyStr := fmt.Sprintf("%v", mapKey.Interface())
		decodedValue, err := decodeValue(mapValue, fieldDesc.MapValue(), field)
		if err != nil {
			return false
		}
		result = append(result, map[string]interface{}{"key": keyStr, "value": decodedValue})
		return true
	})
	return result, nil
}
