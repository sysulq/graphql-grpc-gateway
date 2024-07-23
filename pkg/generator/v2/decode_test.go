package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestDecodeMaps(t *testing.T) {
	ins := &Generator{}
	msg := &test.Maps{}

	generated, err := ins.MarshalProto2GraphQL(msg, &ast.Field{
		Alias:        "constructsMaps_",
		Name:         "constructsMaps_",
		SelectionSet: ast.SelectionSet{},
	})

	require.Nil(t, err)
	require.NotNil(t, generated)
}
