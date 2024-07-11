package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/nautilus/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/sysulq/graphql-gateway/pkg/generator"
)

type anyMap = map[string]interface{}

type queryer struct {
	kod.Implements[Queryer]

	registry kod.Ref[Registry]
	caller   kod.Ref[Caller]
}

func (q *queryer) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kratelimit.Interceptor(),
	}
}

func (q *queryer) Query(ctx context.Context, input *graphql.QueryInput, result interface{}) error {
	res := map[string]interface{}{}
	var err error
	var selection ast.SelectionSet
	for _, op := range input.QueryDocument.Operations {
		selection, err = graphql.ApplyFragments(op.SelectionSet, input.QueryDocument.Fragments)
		if err != nil {
			return err
		}
		switch op.Operation {
		case ast.Query:
			// we allow single flight for queries right now
			ctx = context.WithValue(ctx, allowSingleFlightKey, true)
			err = q.resolveQuery(ctx, selection, res, input.Variables)

		case ast.Mutation:
			err = q.resolveMutation(ctx, selection, res, input.Variables)

		case ast.Subscription:
			return &graphql.Error{
				Extensions: map[string]interface{}{"code": "UNIMPLEMENTED"},
				Message:    "subscription is not supported",
			}
		}
	}
	if err != nil {
		return &graphql.Error{
			Extensions: map[string]interface{}{"code": "UNKNOWN_ERROR"},
			Message:    err.Error(),
		}
	}

	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}
	val = val.Elem()
	if !val.CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}
	val.Set(reflect.ValueOf(res))
	return nil
}

func (q *queryer) resolveMutation(ctx context.Context, selection ast.SelectionSet, res anyMap, vars map[string]interface{}) (err error) {
	for _, ss := range selection {
		field, ok := ss.(*ast.Field)
		if !ok {
			continue
		}
		if field.Name == "__typename" {
			res[nameOrAlias(field)] = field.ObjectDefinition.Name
			continue
		}
		res[nameOrAlias(field)], err = q.resolveCall(ctx, ast.Mutation, field, vars)
		if err != nil {
			return err
		}
	}
	return
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

func (q *queryer) resolveQuery(ctx context.Context, selection ast.SelectionSet, res anyMap, vars map[string]interface{}) (err error) {
	type mapEntry struct {
		key string
		val interface{}
	}
	errCh := make(chan error, 4)
	resCh := make(chan mapEntry, 4)
	for _, ss := range selection {
		field, ok := ss.(*ast.Field)
		if !ok {
			continue
		}
		go func(field *ast.Field) {
			if field.Name == "__typename" {
				resCh <- mapEntry{
					key: nameOrAlias(field),
					val: field.ObjectDefinition.Name,
				}
				return
			}
			resolvedValue, err := q.resolveCall(ctx, ast.Query, field, vars)
			if err != nil {
				errCh <- err
				return
			}
			resCh <- mapEntry{
				key: nameOrAlias(field),
				val: resolvedValue,
			}
		}(field)
	}
	var errs graphql.ErrorList
	for i := 0; i < len(selection); i++ {
		select {
		case r := <-resCh:
			res[r.key] = r.val
		case err := <-errCh:
			errs = append(errs, err)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return
}

func (q *queryer) resolveCall(ctx context.Context, op ast.Operation, field *ast.Field, vars map[string]interface{}) (interface{}, error) {
	method := q.registry.Get().FindMethodByName(op, field.Name)
	if method == nil {
		return nil, errors.New("method not found")
	}

	inputMsg, err := q.pbEncode(method.Input(), field, vars)
	if err != nil {
		return nil, err
	}

	msg, err := q.caller.Get().Call(ctx, method, inputMsg)
	if err != nil {
		return nil, err
	}

	return q.pbDecode(field, msg)
}

func (q *queryer) pbEncode(in protoreflect.MessageDescriptor, field *ast.Field, vars map[string]interface{}) (proto.Message, error) {
	inputMsg := dynamicpb.NewMessage(in)
	inArg := field.Arguments.ForName("in")
	if inArg == nil {
		return inputMsg, nil
	}

	var anyObj protoreflect.MessageDescriptor
	if generator.IsAny(in) {
		if len(inArg.Value.Children) == 0 {
			return nil, errors.New("no '__typename' provided")
		}
		typename := inArg.Value.Children.ForName("__typename")
		if typename == nil {
			return nil, errors.New("no '__typename' provided")
		}

		vv, err := typename.Value(vars)
		if err != nil {
			return nil, errors.New("no '__typename' provided")
		}
		vvv, ok := vv.(string)
		if !ok {
			return nil, errors.New("no '__typename' provided")
		}

		obj := q.registry.Get().FindObjectByName(vvv)
		if obj == nil {
			return nil, errors.New("__typename should be a valid typename")
		}
		anyObj = obj
		inputMsg = dynamicpb.NewMessage(anyObj)
	}

	for _, arg := range inArg.Value.Children {
		val, err := arg.Value.Value(vars)
		if err != nil {
			return nil, err
		}
		if arg.Name == "__typename" {
			continue
		}

		var reqDesc protoreflect.FieldDescriptor
		if anyObj != nil {
			reqDesc = q.registry.Get().FindFieldByName(anyObj, arg.Name)
		} else {
			reqDesc = q.registry.Get().FindFieldByName(in, arg.Name)
		}

		if val, err = q.pbValue(val, reqDesc); err != nil {
			return nil, err
		}

		if reqDesc.Cardinality() == protoreflect.Repeated && reflect.TypeOf(val).Kind() != reflect.Slice {
			inputMsg.Mutable(reqDesc).List().Append(protoreflect.ValueOf(val))
		} else {
			inputMsg.Set(reqDesc, protoreflect.ValueOf(val))
		}
	}

	// set fieldmask based on the selection set
	for i := 0; i < inputMsg.Descriptor().Fields().Len(); i++ {
		v := inputMsg.Descriptor().Fields().Get(i)
		if v.Message() != nil && v.Message().FullName() == "google.protobuf.FieldMask" {

			fm := &fieldmaskpb.FieldMask{
				Paths: getSelectionSet(field.SelectionSet, ""),
			}
			inputMsg.Set(v, protoreflect.ValueOf(fm.ProtoReflect()))
		}
	}

	if anyObj != nil {
		// anyMsgDesc := q.pm.messages[in.DescriptorProto]
		// anyMsg := dynamicpb.NewMessage(q.pm.messages[in.DescriptorProto])
		// typeUrl := anyMsgDesc.FindFieldByName("type_url")
		// value := anyMsgDesc.FindFieldByName("value")
		// anyMsg.SetField(typeUrl, anyObj.GetFullyQualifiedName())
		return anypb.New(inputMsg)
	}
	return inputMsg, nil
}

func (q *queryer) pbValue(val interface{}, reqDesc protoreflect.FieldDescriptor) (_ interface{}, err error) {
	var msgDesc protoreflect.MessageDescriptor

	switch v := val.(type) {
	case float64:
		if reqDesc.Kind() == protoreflect.FloatKind {
			return float32(v), nil
		}
	case int64:
		switch reqDesc.Kind() {
		case protoreflect.Int32Kind,
			protoreflect.Sint32Kind,
			protoreflect.Sfixed32Kind:
			return int32(v), nil

		case protoreflect.Uint32Kind,
			protoreflect.Fixed32Kind:
			return uint32(v), nil

		case protoreflect.Uint64Kind,
			protoreflect.Fixed64Kind:
			return uint64(v), nil
		case protoreflect.FloatKind:
			return float32(v), nil
		case protoreflect.DoubleKind:
			return float64(v), nil
		}
	case string:
		switch reqDesc.Kind() {
		case protoreflect.EnumKind:
			// TODO predefine this
			enumDesc := reqDesc.Enum()
			values := map[string]int32{}
			for i := 0; i < enumDesc.Values().Len(); i++ {
				v := enumDesc.Values().Get(i)
				values[string(v.Name())] = int32(v.Number())
			}
			return values[v], nil
		case protoreflect.BytesKind:
			bytes, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("bytes should be a base64 encoded string")
			}
			return bytes, nil
		}
	case []interface{}:
		v2 := make([]interface{}, len(v))
		for i, vv := range v {
			v2[i], err = q.pbValue(vv, reqDesc)
			if err != nil {
				return nil, err
			}
		}
		return v2, nil
	case map[string]interface{}:
		var anyTypeDescriptor protoreflect.MessageDescriptor
		var vvv string
		var ok bool
		for kk, vv := range v {
			if kk == "__typename" {
				vvv, ok = vv.(string)
				if !ok {
					return nil, errors.New("'__typename' must be a string")
				}
				delete(v, "__typename")
				break
			}
		}
		var msg *dynamicpb.Message
		if reqDesc != nil {
			msgDesc = reqDesc.Message()
		} else {
			msgDesc = q.registry.Get().FindObjectByName(vvv)
		}

		protoDesc := msgDesc

		if generator.IsAny(protoDesc) {
			anyTypeDescriptor = q.registry.Get().FindObjectByName(vvv)
			if anyTypeDescriptor == nil {
				return nil, errors.New("'__typename' must be a valid INPUT_OBJECT")
			}
			msg = dynamicpb.NewMessage(anyTypeDescriptor)
		} else {
			msg = dynamicpb.NewMessage(protoDesc)
		}
		oneofValidate := map[protoreflect.OneofDescriptor]struct{}{}
		for kk, vv := range v {

			if anyTypeDescriptor != nil {
				msgDesc = anyTypeDescriptor
			}
			// plugType := q.pm.inputs[msgDesc]
			// fieldDesc := q.pm.fields[q.p.FieldBack(plugType.DescriptorProto, kk)]
			fieldDesc := q.registry.Get().FindFieldByName(msgDesc, kk)
			oneof := fieldDesc.ContainingOneof()
			if oneof != nil {
				_, ok := oneofValidate[oneof]
				if ok {
					return nil, fmt.Errorf("field with name %q on Object %q can't be set", kk, protoDesc.Name())
				}
				oneofValidate[oneof] = struct{}{}
			}

			vv2, err := q.pbValue(vv, fieldDesc)
			if err != nil {
				return nil, err
			}
			msg.Set(fieldDesc, protoreflect.ValueOf(vv2))
		}
		if anyTypeDescriptor != nil {
			return anypb.New(msg)
		}
		return msg, nil
	}

	return val, nil
}

func (q *queryer) pbDecodeOneofField(desc protoreflect.MessageDescriptor, dynamicMsg proto.Message, selection ast.SelectionSet) (oneof anyMap, err error) {
	oneof = anyMap{}
	for _, f := range selection {
		out, ok := f.(*ast.Field)
		if !ok {
			continue
		}
		if out.Name == "__typename" {
			oneof[nameOrAlias(out)] = out.ObjectDefinition.Name
			continue
		}

		fieldDesc := q.registry.Get().FindUnionFieldByMessageFQNAndName(string(desc.FullName()), out.Name)
		protoVal := dynamicMsg.ProtoReflect().Get(fieldDesc)
		oneof[nameOrAlias(out)], err = q.gqlValue(protoVal, fieldDesc.Message(), fieldDesc.Enum(), out)
		if err != nil {
			return nil, err
		}
	}
	return oneof, nil
}

func (q *queryer) pbDecode(field *ast.Field, msg proto.Message) (res interface{}, err error) {
	switch dynamicMsg := msg.(type) {
	case *dynamicpb.Message:
		return q.gqlValue(dynamicMsg, dynamicMsg.Descriptor(), nil, field)
	case *anypb.Any:
		return q.gqlValue(dynamicMsg, nil, nil, field)
	default:
		return nil, fmt.Errorf("expected proto message of type *dynamic.Message or *anypb.Any but received: %T", msg)
	}
}

// FIXME take care of recursive calls
func (q *queryer) gqlValue(val interface{}, msgDesc protoreflect.MessageDescriptor, enumDesc protoreflect.EnumDescriptor, field *ast.Field) (_ interface{}, err error) {
	switch v := val.(type) {
	case int32:
		// int32 enum
		if enumDesc != nil {
			values := map[int32]string{}
			for i := 0; i < enumDesc.Values().Len(); i++ {
				v := enumDesc.Values().Get(i)
				values[int32(v.Number())] = string(v.Name())
			}

			return values[v], nil
		}

	case map[interface{}]interface{}:
		res := make([]interface{}, len(v))
		i := 0
		for kk, vv := range v {
			vals := anyMap{}
			for _, f := range field.SelectionSet {
				out, ok := f.(*ast.Field)
				if !ok {
					continue
				}
				switch out.Name {
				case "value":
					valueField := msgDesc.Fields().ByName("value")
					if vals[nameOrAlias(out)], err = q.gqlValue(vv, valueField.Message(), valueField.Enum(), out); err != nil {
						return nil, err
					}
				case "key":
					vals[nameOrAlias(out)] = kk
				case "__typename":
					vals[nameOrAlias(out)] = out.ObjectDefinition.Name
				}
			}

			res[i] = vals
			i++
		}
		return res, nil

	case []interface{}:
		v2 := make([]interface{}, len(v))
		for i, vv := range v {
			v2[i], err = q.gqlValue(vv, msgDesc, enumDesc, field)
			if err != nil {
				return nil, err
			}
		}
		return v2, nil

	case *dynamicpb.Message:
		if v == nil {
			return nil, nil
		}
		fields := msgDesc.Fields()
		vals := make(map[string]interface{}, fields.Len())
		// gqlFields := map[string]string{}
		for _, s := range field.SelectionSet {
			out, ok := s.(*ast.Field)
			if !ok {
				continue
			}

			if out.Name == "__typename" {
				vals[nameOrAlias(out)] = out.ObjectDefinition.Name
				continue
			}

			fieldDesc := q.registry.Get().FindFieldByName(msgDesc, out.Name)
			if fieldDesc == nil {
				vals[nameOrAlias(out)], err = q.pbDecodeOneofField(msgDesc, v, out.SelectionSet)
				if err != nil {
					return nil, err
				}
				continue
			}

			vals[nameOrAlias(out)], err = q.gqlValue(v.Get(fieldDesc).Interface(), fieldDesc.Message(), fieldDesc.Enum(), out)
			if err != nil {
				return nil, err
			}
		}
		return vals, nil
	case *anypb.Any:
		vals, err := q.anyMessageToMap(v)
		if err != nil {
			return nil, err
		}
		return vals, nil

	case []byte:
		return base64.StdEncoding.EncodeToString(v), nil
	}

	return val, nil
}

func (q *queryer) anyMessageToMap(v *anypb.Any) (map[string]interface{}, error) {
	fqn := string(v.MessageName())

	grpcType, definition := q.registry.Get().FindObjectByFullyQualifiedName(fqn)
	outputMsg := dynamicpb.NewMessage(grpcType)
	if err := anypb.UnmarshalTo(v, outputMsg, proto.UnmarshalOptions{}); err != nil {
		return nil, err
	}
	return q.protoMessageToMap(outputMsg, definition)
}

func (q *queryer) protoMessageToMap(outputMsg *dynamicpb.Message, definition *ast.Definition) (map[string]interface{}, error) {
	fields := outputMsg.Descriptor().Fields()
	vals := make(map[string]interface{}, fields.Len())
	vals["__typename"] = definition.Name
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldDef := q.registry.Get().FindGraphqlFieldByProtoField(definition, string(field.Name()))
		// the field is probably invalid or ignored
		if fieldDef == nil {
			continue
			// return nil, fmt.Errorf("proto field %q doesn't have a graphql counterpart on type %q", field.GetName(), definition.Name)
		}
		val := outputMsg.Get(field)
		switch vv := val.Interface().(type) {
		case int32:
			if field.Enum() != nil {
				values := map[int32]string{}
				for i := 0; i < field.Enum().Values().Len(); i++ {
					v := field.Enum().Values().Get(i)
					values[int32(v.Number())] = string(v.Name())
				}

				vals[fieldDef.Name] = values[vv]
			}

		case *dynamicpb.Message:
			_, definition := q.registry.Get().FindObjectByFullyQualifiedName(string(vv.Descriptor().FullName()))
			val, err := q.protoMessageToMap(vv, definition)
			if err != nil {
				return nil, err
			}
			vals[fieldDef.Name] = val
		case *anypb.Any:
			val, err := q.anyMessageToMap(vv)
			if err != nil {
				return nil, err
			}
			vals[fieldDef.Name] = val
		case []interface{}:
			var arrayVals []interface{}
			for _, val := range vv {
				switch vv := val.(type) {
				case int32:
					if field.Enum() != nil {
						values := map[int32]string{}
						for i := 0; i < field.Enum().Values().Len(); i++ {
							v := field.Enum().Values().Get(i)
							values[int32(v.Number())] = string(v.Name())
						}

						arrayVals = append(arrayVals, values[vv])
					}

				case *dynamicpb.Message:
					_, definition := q.registry.Get().FindObjectByFullyQualifiedName(string(vv.Descriptor().FullName()))
					val, err := q.protoMessageToMap(vv, definition)
					if err != nil {
						return nil, err
					}
					arrayVals = append(arrayVals, val)
				case *anypb.Any:
					val, err := q.anyMessageToMap(vv)
					if err != nil {
						return nil, err
					}
					arrayVals = append(arrayVals, val)
				default:
					arrayVals = append(arrayVals, vv)
				}
			}
		default:
			vals[fieldDef.Name] = vv
		}
	}
	return vals, nil
}

func nameOrAlias(field *ast.Field) string {
	if field.Alias != "" {
		return field.Alias
	}

	return field.Name
}
