package server

import (
	"context"
	"errors"
	"reflect"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/nautilus/graphql"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/sysulq/graphql-grpc-gateway/internal/config"
)

type anyMap = map[string]interface{}

type graphqlQueryer struct {
	kod.Implements[GraphqlQueryer]

	config   kod.Ref[config.Config]
	caller   kod.Ref[GraphqlCaller]
	registry kod.Ref[GraphqlCallerRegistry]
}

func (q *graphqlQueryer) Interceptors() []interceptor.Interceptor {
	if q.config.Get().Config().Engine.RateLimit {
		return []interceptor.Interceptor{
			kratelimit.Interceptor(),
		}
	}

	return nil
}

func (q *graphqlQueryer) Query(ctx context.Context, input *graphql.QueryInput, result interface{}) error {
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

func (q *graphqlQueryer) resolveMutation(ctx context.Context, selection ast.SelectionSet, res anyMap, vars map[string]interface{}) (err error) {
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

func (q *graphqlQueryer) resolveQuery(ctx context.Context, selection ast.SelectionSet, res anyMap, vars map[string]interface{}) (err error) {
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

func (q *graphqlQueryer) resolveCall(ctx context.Context, op ast.Operation, field *ast.Field, vars map[string]interface{}) (interface{}, error) {
	method := q.registry.Get().FindMethodByName(op, field.Name)
	if method == nil {
		return nil, errors.New("method not found")
	}

	inputMsg, err := q.registry.Get().Unmarshal(method.Input(), field, vars)
	if err != nil {
		return nil, err
	}

	msg, err := q.caller.Get().Call(ctx, method, inputMsg)
	if err != nil {
		return nil, err
	}

	return q.registry.Get().Marshal(msg, field)
}

func nameOrAlias(field *ast.Field) string {
	if field.Alias != "" {
		return field.Alias
	}

	return field.Name
}
