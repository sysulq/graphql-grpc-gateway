package protographql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestDecodeMaps(t *testing.T) {
	ins := New()
	err := ins.RegisterFileDescriptor(true, test.File_test_constructs_input_proto)
	require.Nil(t, err)

	msg := &test.Maps{}

	generated, err := ins.Marshal(msg, &ast.Field{
		Alias:        "constructsMaps_",
		Name:         "constructsMaps_",
		SelectionSet: ast.SelectionSet{},
	})

	require.Nil(t, err)
	require.NotNil(t, generated)
}
