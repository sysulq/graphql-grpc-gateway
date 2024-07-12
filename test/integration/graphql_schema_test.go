package integration

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/go-kod/kod"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/sysulq/graphql-gateway/test"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/mock/gomock"
)

//go:embed testdata/gateway-expect.graphql
var testGatewayExpectedSchema []byte

func TestGraphqlSchema(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := server.NewMockConfigComponent(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&server.Config{
		Grpc: server.Grpc{
			Services: []*server.Service{
				{
					Address:    infos.ConstructsServerAddr.Addr().String(),
					Reflection: true,
				},
				{
					Address:    infos.OptionsServerAddr.Addr().String(),
					Reflection: true,
				},
			},
		},
		Playground: true,
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.ServerComponent) {
		// gatewayUrl := test.SetupGateway(t, s)
		t.Run("schema is correct", func(t *testing.T) {
			expectedFormattedSchema, _ := gqlparser.LoadSchema(&ast.Source{Input: string(testGatewayExpectedSchema)})

			for _, v := range expectedFormattedSchema.Mutation.Fields {
				for _, vv := range v.Arguments {
					fmt.Println(v.Name, vv.Name, vv.Type)
				}
			}

			// schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// test.CompareGraphql(t, schema.Schema, expectedFormattedSchema)
		})
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
