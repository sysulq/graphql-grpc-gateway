package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEncodeMaps(t *testing.T) {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: schemaGraphQL,
	})
	require.Nil(t, err)
	var msg test.Maps

	queryDoc, err := gqlparser.LoadQuery(schema, constructsMapsQuery)
	require.Nil(t, err)

	ins := &Generator{}
	generated, err := ins.GraphQL2Proto(msg.ProtoReflect().Descriptor(), queryDoc.Operations[0].SelectionSet[0].(*ast.Field), nil)

	require.Nil(t, err)

	data, err := ins.MarshalProto2GraphQL(generated, queryDoc.Operations[0].SelectionSet[0].(*ast.Field))
	require.Nil(t, err)

	require.Equal(t, constructsMapsResponse, data)
}

func TestEncodeRepeated(t *testing.T) {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: schemaGraphQL,
	})
	require.Nil(t, err)
	var msg test.Repeated

	queryDoc, err := gqlparser.LoadQuery(schema, constructsRepeatedQuery)
	require.Nil(t, err, err.Error())

	ins := &Generator{}
	generated, err := ins.GraphQL2Proto(msg.ProtoReflect().Descriptor(), queryDoc.Operations[0].SelectionSet[0].(*ast.Field), nil)

	require.Nil(t, err)

	data, err := ins.MarshalProto2GraphQL(generated, queryDoc.Operations[0].SelectionSet[0].(*ast.Field))
	require.Nil(t, err)

	require.Equal(t, constructsRepeatedResponse, data)
}
