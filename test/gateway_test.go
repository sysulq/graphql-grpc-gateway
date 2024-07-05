package test

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"github.com/nautilus/graphql"
)

//go:embed testdata/gateway-expect.graphql
var testGatewayExpectedSchema []byte

func Test_Gateway(t *testing.T) {
	SetupDeps()

	querier := graphql.NewSingleRequestQueryer("")
	// t.Run("schema is correct", func(t *testing.T) {
	// 	schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	expectedFormattedSchema, _ := gqlparser.LoadSchema(&ast.Source{Input: string(testGatewayExpectedSchema)})
	// 	compareGraphql(t, schema.Schema, expectedFormattedSchema)
	// })

	tests := []struct {
		name         string
		query        string
		vars         j
		wantResponse j
	}{{
		name:         "Mutation constructs scalars",
		query:        constructsScalarsQuery,
		wantResponse: constructsScalarsResponse,
	}, {
		name:         "Mutation constructs any",
		query:        constructsAnyQuery,
		wantResponse: constructsAnyResponse,
	}, {
		name:         "Mutation constructs maps",
		query:        constructsMapsQuery,
		wantResponse: constructsMapsResponse,
	}, {
		name:         "Mutation constructs repeated",
		query:        constructsRepeatedQuery,
		wantResponse: constructsRepeatedResponse,
	}, {
		name:         "Mutation constructs oneofs", // TODO check oneof input params not to allow inputting data for the same oneof input
		query:        constructsOneofsQuery,
		wantResponse: constructsOneofsResponse,
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recv := map[string]interface{}{}
			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query:     tc.query,
				Variables: tc.vars,
			}, &recv); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(recv, tc.wantResponse) {
				t.Errorf("mutation failed: expected: %s got: %s", tc.wantResponse, recv)
			}
		})
	}
}
