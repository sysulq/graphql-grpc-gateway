package protographql

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/validator"
)

//go:embed schema.graphql
var schemaGraphQL string

func TestSchema(t *testing.T) {
	ins := New()
	var err error
	err = ins.GenerateFile(true, test.File_test_constructs_input_proto)
	require.Nil(t, err)
	err = ins.GenerateFile(true, test.File_test_options_input_proto)
	require.Nil(t, err)

	schema := ins.AsGraphQL()
	file, err := os.Create("schema.graphql")
	require.Nil(t, err)
	formatter.NewFormatter(file).FormatSchema(schema)

	_, err = validator.LoadSchema(validator.Prelude, &ast.Source{Input: schemaGraphQL})
	require.Nil(t, err)
}
