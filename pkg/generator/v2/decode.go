package v2

import (
	"encoding/base64"
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MarshalProto2GraphQL 将 gRPC 返回的 proto.Message 转换为 GraphQL 数据
func (ins *Generator) MarshalProto2GraphQL(msg proto.Message, field *ast.Field) (interface{}, error) {
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
			marshaldValue, err := marshalValue(value, fieldDesc, f)
			if err != nil {
				return nil, err
			}
			result[fieldName] = marshaldValue
		}
	}

	return result, nil
}

// marshalValue 将 protoreflect.Value 转换为适合 GraphQL 的值
func marshalValue(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	switch {
	case fieldDesc.IsList():
		return marshalList(value.List(), fieldDesc, field)
	case fieldDesc.IsMap():
		return marshalMap(value.Map(), fieldDesc, field)
	}

	return marshalScalar(value, fieldDesc, field)
}

func marshalScalar(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
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
	case protoreflect.FloatKind:
		return float32(value.Float()), nil
	case protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return base64.StdEncoding.EncodeToString(value.Bytes()), nil
	case protoreflect.EnumKind:
		return string(fieldDesc.Enum().Values().ByNumber(value.Enum()).Name()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if fieldDesc.IsMap() {
			return marshalMap(value.Map(), fieldDesc, field)
		}

		return marshalMessage(value.Message(), field)
	default:
		return nil, fmt.Errorf("unsupported field type: %v", fieldDesc.Kind())
	}
}

func marshalMessage(message protoreflect.Message, field *ast.Field) (interface{}, error) {
	result := make(map[string]interface{})
	for _, selection := range field.SelectionSet {
		if f, ok := selection.(*ast.Field); ok {
			fieldName := f.Alias
			if fieldName == "" {
				fieldName = f.Name
			}

			fieldDesc := message.Descriptor().Fields().ByJSONName(fieldName)
			if fieldDesc == nil {
				continue
			}
			value := message.Get(fieldDesc)
			marshaldValue, err := marshalValue(value, fieldDesc, f)
			if err != nil {
				return nil, err
			}
			result[fieldName] = marshaldValue
		}
	}
	return result, nil
}

// marshalList 将 protoreflect.List 转换为 []interface{}
func marshalList(list protoreflect.List, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	result := make([]interface{}, list.Len())
	for i := 0; i < list.Len(); i++ {
		value := list.Get(i)
		var marshaldValue interface{}
		var err error

		if fieldDesc.Message() != nil {
			marshaldValue, err = marshalMessage(value.Message(), field)
		} else {
			marshaldValue, err = marshalScalar(value, fieldDesc, field)
		}

		if err != nil {
			return nil, err
		}
		result[i] = marshaldValue
	}
	return result, nil
}

func marshalMap(m protoreflect.Map, fieldDesc protoreflect.FieldDescriptor, field *ast.Field) (interface{}, error) {
	result := make([]map[string]interface{}, 0)
	m.Range(func(mapKey protoreflect.MapKey, mapValue protoreflect.Value) bool {
		var valueField *ast.Field
		for _, selection := range field.SelectionSet {
			if f, ok := selection.(*ast.Field); ok {
				if f.Name == "value" {
					valueField = f
				}
			}
		}

		marshaldValue, err := marshalValue(mapValue, fieldDesc.MapValue(), valueField)
		if err != nil {
			return false
		}
		result = append(result, map[string]interface{}{"key": mapKey.Interface(), "value": marshaldValue})
		return true
	})
	return result, nil
}
