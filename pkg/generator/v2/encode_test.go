package v2

import (
	"testing"

	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
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

func TestEncodeDecode(t *testing.T) {
	cases := []struct {
		name         string
		query        string
		wantResponse interface{}
		msg          proto.Message
	}{
		{
			name:         "encode decode scalar",
			query:        constructsScalarsQuery,
			wantResponse: constructsScalarsResponse,
			msg:          &test.Scalars{},
		},
		{
			name:         "encode decode maps",
			query:        constructsMapsQuery,
			wantResponse: constructsMapsResponse,
			msg:          &test.Maps{},
		},
		{
			name:         "encode decode repeated",
			query:        constructsRepeatedQuery,
			wantResponse: constructsRepeatedResponse,
			msg:          &test.Repeated{},
		},
		{
			name:         "encode decode oneof",
			query:        constructsOneofsQuery,
			wantResponse: constructsOneofsResponse,
			msg:          &test.Oneof{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			schema, err := gqlparser.LoadSchema(&ast.Source{
				Input: schemaGraphQL,
			})
			require.Nil(t, err, err)

			queryDoc, err := gqlparser.LoadQuery(schema, c.query)
			require.Nil(t, err, err)

			ins := &Generator{}

			selection, err := graphql.ApplyFragments(queryDoc.Operations[0].SelectionSet, queryDoc.Fragments)
			require.Nil(t, err, err)

			generated, err := ins.GraphQL2Proto(c.msg.ProtoReflect().Descriptor(), selection[0].(*ast.Field), nil)

			require.Nil(t, err)

			data, err := ins.MarshalProto2GraphQL(generated, selection[0].(*ast.Field))
			require.Nil(t, err)

			require.Equal(t, c.wantResponse, data)
		})
	}
}
