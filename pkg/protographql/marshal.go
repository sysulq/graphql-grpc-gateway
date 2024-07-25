package protographql

import (
	"encoding/base64"
	"fmt"
	"maps"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

// Marshal 将 gRPC 返回的 proto.Message 转换为 GraphQL 数据
func (ins *SchemaDescriptor) Marshal(msg proto.Message, field *ast.Field) (interface{}, error) {
	message := msg.ProtoReflect()
	result := make(map[string]interface{})

	if IsAny(msg.ProtoReflect().Descriptor()) {
		anyMsg, ok := msg.(*anypb.Any)
		if !ok {
			return nil, fmt.Errorf("failed to cast to Any type")
		}

		msgDesc, err := ins.ProtoTypes.FindMessageByURL(anyMsg.GetTypeUrl())
		if err != nil {
			return nil, err
		}

		msgObj := msgDesc.New().Interface()
		err = proto.Unmarshal(anyMsg.Value, msgObj)
		if err != nil {
			return nil, err
		}

		return ins.marshalMessage(msgObj.ProtoReflect(), field, true)
	}

	for _, selection := range field.SelectionSet {
		if f, ok := selection.(*ast.Field); ok {
			fieldName := f.Name
			if fieldName == "__typename" {
				result["__typename"] = string(message.Descriptor().Name())
				continue
			}

			fieldDesc := message.Descriptor().Fields().ByJSONName(fieldName)
			if fieldDesc == nil {
				oneofMap, err := ins.marshalOneof(message, f, false)
				if err != nil {
					return nil, err
				}

				result[fieldName] = oneofMap
				continue
			}
			value := message.Get(fieldDesc)
			marshaldValue, err := ins.marshalValue(value, fieldDesc, f, false)
			if err != nil {
				return nil, err
			}
			result[fieldName] = marshaldValue
		}
	}

	return result, nil
}

// marshalValue 将 protoreflect.Value 转换为适合 GraphQL 的值
func (ins *SchemaDescriptor) marshalValue(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field, allField bool) (interface{}, error) {
	if fieldDesc.IsList() {
		return ins.marshalList(value.List(), fieldDesc, field, allField)
	} else if fieldDesc.IsMap() {
		return ins.marshalMap(value.Map(), fieldDesc, field, allField)
	}

	return ins.marshalScalar(value, fieldDesc, field, allField)
}

// marshalOneof 将 oneof 类型字段转换为 map[string]interface{}
func (ins *SchemaDescriptor) marshalOneof(message protoreflect.Message, field *ast.Field, allField bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	oneofDesc := message.Descriptor().Oneofs().ByName(protoreflect.Name(field.Name))
	if oneofDesc == nil {
		return result, nil
	}

	fieldDesc := message.WhichOneof(oneofDesc)
	if fieldDesc == nil {
		return result, nil
	}

	value := message.Get(fieldDesc)

	for _, selects := range field.SelectionSet {
		if f, ok := selects.(*ast.Field); ok {
			if f.Name == "__typename" {
				result["__typename"] = field.Definition.Type.NamedType
				continue
			}

			if string(fieldDesc.Name()) != f.Name {
				continue
			}

			marshaldValue, err := ins.marshalValue(value, fieldDesc, f, allField)
			if err != nil {
				return nil, err
			}
			result[f.Name] = marshaldValue
		}
	}

	return result, nil
}

func (ins *SchemaDescriptor) marshalScalar(value protoreflect.Value, fieldDesc protoreflect.FieldDescriptor, field *ast.Field, allField bool) (interface{}, error) {
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
			return ins.marshalMap(value.Map(), fieldDesc, field, allField)
		}

		return ins.marshalMessage(value.Message(), field, allField)
	default:
		return nil, fmt.Errorf("unsupported field type: %v", fieldDesc.Kind())
	}
}

func (ins *SchemaDescriptor) getSelectionSetFromMessage(message protoreflect.Message) ast.SelectionSet {
	selectionSet := make(ast.SelectionSet, 0)

	for i := 0; i < message.Descriptor().Fields().Len(); i++ {
		field := message.Descriptor().Fields().Get(i)
		if field.Kind() == protoreflect.MessageKind && IsEmpty(field.Message()) {
			continue
		}

		selectionSet = append(selectionSet, &ast.Field{
			Name: field.JSONName(),
		})
	}

	if len(selectionSet) > 0 {
		selectionSet = append(selectionSet, &ast.Field{
			Name: "__typename",
		})
	}

	return selectionSet
}

func (ins *SchemaDescriptor) marshalMessage(message protoreflect.Message, field *ast.Field, allField bool) (interface{}, error) {
	if allField {
		field.SelectionSet = ins.getSelectionSetFromMessage(message)
	}

	result := make(map[string]interface{})
	for _, selection := range field.SelectionSet {
		if f, ok := selection.(*ast.Field); ok {
			fieldName := f.Alias
			if fieldName == "" {
				fieldName = f.Name
			}
			if fieldName == "__typename" {
				result[fieldName] = string(message.Descriptor().Name())
				continue
			}

			fieldDesc := message.Descriptor().Fields().ByJSONName(fieldName)
			if fieldDesc == nil {
				oneofField, err := ins.marshalOneof(message, f, allField)
				if err != nil {
					return nil, err
				}

				maps.Copy(result, oneofField)
				continue
			}

			value := message.Get(fieldDesc)
			marshaldValue, err := ins.marshalValue(value, fieldDesc, f, allField)
			if err != nil {
				return nil, err
			}
			result[fieldName] = marshaldValue
		}
	}
	return result, nil
}

// marshalList 将 protoreflect.List 转换为 []interface{}
func (ins *SchemaDescriptor) marshalList(list protoreflect.List, fieldDesc protoreflect.FieldDescriptor, field *ast.Field, allField bool) (interface{}, error) {
	result := make([]interface{}, list.Len())
	for i := 0; i < list.Len(); i++ {
		value := list.Get(i)
		var marshaldValue interface{}
		var err error

		if fieldDesc.Message() != nil {
			marshaldValue, err = ins.marshalMessage(value.Message(), field, allField)
		} else {
			marshaldValue, err = ins.marshalScalar(value, fieldDesc, field, allField)
		}

		if err != nil {
			return nil, err
		}
		result[i] = marshaldValue
	}
	return result, nil
}

func (ins *SchemaDescriptor) marshalMap(m protoreflect.Map, fieldDesc protoreflect.FieldDescriptor, field *ast.Field, allField bool) (interface{}, error) {
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

		marshaldValue, err := ins.marshalValue(mapValue, fieldDesc.MapValue(), valueField, allField)
		if err != nil {
			return false
		}
		result = append(result, map[string]interface{}{"key": mapKey.Interface(), "value": marshaldValue})
		return true
	})
	return result, nil
}
