package protographql

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (ins *SchemaDescriptor) newMessage(desc protoreflect.MessageDescriptor) protoreflect.Message {
	msgType, err := ins.ProtoTypes.FindMessageByName(desc.FullName())
	if err != nil {
		return dynamicpb.NewMessage(desc)
	}

	return msgType.New()
}

// Unmarshal 将 GraphQL 字段转换为 Protocol Buffers 消息
func (ins *SchemaDescriptor) Unmarshal(desc protoreflect.MessageDescriptor, field *ast.Field, vars map[string]interface{}) (proto.Message, error) {
	// 查找名为 "in" 的参数
	var inArg *ast.Argument
	for _, arg := range field.Arguments {
		if arg.Name == "in" {
			inArg = arg
			break
		}
	}
	if inArg == nil {
		return nil, fmt.Errorf("argument 'in' not found in the field")
	}

	// 创建一个动态消息
	dynamicMessage := ins.newMessage(desc)

	// 设置字段值
	if inArg.Value.Kind == ast.ObjectValue {
		for _, argField := range inArg.Value.Children {
			fieldDescriptor := desc.Fields().ByJSONName(argField.Name)
			if fieldDescriptor == nil {
				return nil, fmt.Errorf("field %s not found in message", argField.Name)
			}

			err := ins.unmarshalValue(dynamicMessage, fieldDescriptor, argField.Value)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("argument 'in' is not an object")
	}

	// 处理 FieldMask 字段
	for i := desc.Fields().Len() - 1; i >= 0; i-- {
		fieldDescriptor := desc.Fields().Get(i)
		if fieldDescriptor.Message() != nil && fieldDescriptor.Message().FullName() == "google.protobuf.FieldMask" {
			fieldMask := &fieldmaskpb.FieldMask{Paths: getSelectionSet(field.SelectionSet, "")}
			dynamicMessage.Set(fieldDescriptor, protoreflect.ValueOfMessage(fieldMask.ProtoReflect()))
			break
		}
	}

	return dynamicMessage.Interface(), nil
}

func getSelectionSet(selections ast.SelectionSet, prefix string) []string {
	var fields []string
	for _, selection := range selections {
		field, ok := selection.(*ast.Field)
		if !ok {
			continue
		}
		fullName := field.Name
		if prefix != "" {
			fullName = fmt.Sprintf("%s.%s", prefix, field.Name)
		}
		if len(field.SelectionSet) > 0 {
			subFields := getSelectionSet(field.SelectionSet, fullName)
			fields = append(fields, subFields...)
		} else {
			fields = append(fields, fullName)
		}
	}
	return fields
}

// unmarshalValue 根据字段类型分发到具体的解码函数
func (ins *SchemaDescriptor) unmarshalValue(msg protoreflect.Message, fd protoreflect.FieldDescriptor, value *ast.Value) error {
	if fd.IsList() {
		return ins.unmarshalList(msg, fd, value)
	} else if fd.IsMap() {
		return ins.unmarshalMap(msg, fd, value)
	} else if fd.Message() != nil {
		return ins.unmarshalMessage(msg, fd, value)
	} else {
		return ins.unmarshalScalar(msg, fd, value)
	}
}

// unmarshalList 解码 List 类型字段
func (ins *SchemaDescriptor) unmarshalList(msg protoreflect.Message, fd protoreflect.FieldDescriptor, value *ast.Value) error {
	if value.Kind != ast.ListValue {
		return fmt.Errorf("expected list value for field %s %+v", fd.Name(), value)
	}

	list := msg.Mutable(fd).List()
	for _, v := range value.Children {
		val, err := ins.convertValue(fd, v.Value)
		if err != nil {
			return err
		}

		list.Append(val)
	}

	return nil
}

// unmarshalMap 解码 GraphQL 对象到 Protocol Buffers map
func (ins *SchemaDescriptor) unmarshalMap(msg protoreflect.Message, fd protoreflect.FieldDescriptor, value *ast.Value) error {
	if value.Kind != ast.ListValue {
		return fmt.Errorf("expected list value for map field %s", fd.JSONName())
	}

	mapValue := msg.Mutable(fd).Map()
	for _, child := range value.Children {
		if child.Value.Kind != ast.ObjectValue {
			return fmt.Errorf("expected object value in list for map field %s", fd.JSONName())
		}

		var key string
		var val *ast.Value

		for _, kv := range child.Value.Children {
			if kv.Name == "key" {
				key = kv.Value.Raw
			} else if kv.Name == "value" {
				val = kv.Value
			}
		}

		if key == "" || val == nil {
			return fmt.Errorf("key or value missing in map entry for field %s", fd.JSONName())
		}

		mapKey, err := unmarshalMapKey(fd.MapKey(), key)
		if err != nil {
			return err
		}
		mapVal, err := ins.convertValue(fd.MapValue(), val)
		if err != nil {
			return err
		}
		mapValue.Set(mapKey, mapVal)
	}
	msg.Set(fd, protoreflect.ValueOfMap(mapValue))
	return nil
}

// unmarshalMessage 解码 GraphQL 对象到 Protocol Buffers message
func (ins *SchemaDescriptor) unmarshalMessage(msg protoreflect.Message, fd protoreflect.FieldDescriptor, value *ast.Value) error {
	if value.Kind != ast.ObjectValue {
		return fmt.Errorf("expected object value for message field %s", fd.JSONName())
	}

	subMsg := ins.newMessage(fd.Message())
	for _, child := range value.Children {
		subFd := fd.Message().Fields().ByJSONName(child.Name)
		if subFd == nil {
			return fmt.Errorf("field %s not found in message %s", child.Name, fd.JSONName())
		}
		err := ins.unmarshalValue(subMsg, subFd, child.Value)
		if err != nil {
			return err
		}
	}
	msg.Set(fd, protoreflect.ValueOfMessage(subMsg))
	return nil
}

// unmarshalScalar 解码 GraphQL 标量值到 Protocol Buffers 标量字段
func (ins *SchemaDescriptor) unmarshalScalar(msg protoreflect.Message, fd protoreflect.FieldDescriptor, value *ast.Value) error {
	val, err := ins.convertValue(fd, value)
	if err != nil {
		return err
	}

	msg.Set(fd, val)
	return nil
}

// unmarshalMapKey 解码 GraphQL map 键
func unmarshalMapKey(fd protoreflect.FieldDescriptor, key string) (protoreflect.MapKey, error) {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOf(key).MapKey(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return protoreflect.MapKey{}, err
		}
		return protoreflect.ValueOf(int32(v)).MapKey(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			return protoreflect.MapKey{}, err
		}
		return protoreflect.ValueOf(uint32(v)).MapKey(), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return protoreflect.MapKey{}, err
		}
		return protoreflect.ValueOf(v).MapKey(), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return protoreflect.MapKey{}, err
		}
		return protoreflect.ValueOf(v).MapKey(), nil
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(key)
		if err != nil {
			return protoreflect.MapKey{}, err
		}
		return protoreflect.ValueOf(v).MapKey(), nil
	default:
		return protoreflect.MapKey{}, fmt.Errorf("unsupported map key type: %v", fd.Kind())
	}
}

// convertValue 将 GraphQL 值转换为 Protocol Buffers 值
func (ins *SchemaDescriptor) convertValue(fd protoreflect.FieldDescriptor, gqlValue *ast.Value) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		boolValue, err := strconv.ParseBool(gqlValue.Raw)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(boolValue), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		intValue, err := strconv.ParseInt(gqlValue.Raw, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(int32(intValue)), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		uintValue, err := strconv.ParseUint(gqlValue.Raw, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(uint32(uintValue)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		intValue, err := strconv.ParseInt(gqlValue.Raw, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(intValue), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		uintValue, err := strconv.ParseUint(gqlValue.Raw, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(uintValue), nil
	case protoreflect.FloatKind:
		floatValue, err := strconv.ParseFloat(gqlValue.Raw, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(float32(floatValue)), nil
	case protoreflect.DoubleKind:
		floatValue, err := strconv.ParseFloat(gqlValue.Raw, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOf(floatValue), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOf(gqlValue.Raw), nil
	case protoreflect.BytesKind:
		// base64 decode
		data, err := base64.StdEncoding.DecodeString(gqlValue.Raw)
		if err != nil {
			return protoreflect.Value{}, err
		}

		return protoreflect.ValueOf(data), nil
	case protoreflect.EnumKind:
		enumDesc := fd.Enum().Values().ByName(protoreflect.Name(gqlValue.Raw))
		if enumDesc == nil {
			return protoreflect.Value{}, fmt.Errorf("enum type not found")
		}
		return protoreflect.ValueOf(enumDesc.Number()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		dynamicMessage := ins.newMessage(fd.Message())
		for _, child := range gqlValue.Children {
			subFd := fd.Message().Fields().ByJSONName(child.Name)
			if subFd == nil {
				return protoreflect.Value{}, fmt.Errorf("field %s not found in message %s", child.Name, fd.JSONName())
			}
			err := ins.unmarshalValue(dynamicMessage, subFd, child.Value)
			if err != nil {
				return protoreflect.Value{}, err
			}
		}
		return protoreflect.ValueOfMessage(dynamicMessage), nil

	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %v", fd.Kind())
	}
}
