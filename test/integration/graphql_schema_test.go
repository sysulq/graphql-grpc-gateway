package integration

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/validator"
	"go.uber.org/mock/gomock"
)

//go:embed testdata/gateway-expect.graphql
var testGatewayExpectedSchema []byte

//go:embed testdata/gateway-expect-without-unbound-method.graphql
var testGatewayExpectedSchemaWithoutUnboundMethod []byte

func TestGraphqlSchema(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Grpc: config.Grpc{
			Services: []kgrpc.Config{
				{
					Target: infos.ConstructsServerAddr.Addr().String(),
				},
				{
					Target: infos.OptionsServerAddr.Addr().String(),
				},
			},
		},
		Server: config.ServerConfig{
			GraphQL: config.GraphQL{
				GenerateUnboundMethods: true,
				Playground:             true,
			},
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		t.Run("schema is correct", func(t *testing.T) {
			schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
			if err != nil {
				t.Fatal(err)
			}

			file, err := os.Create("testdata/gateway-generate.graphql")
			require.Nil(t, err)
			formatter.NewFormatter(file).FormatSchema(schema.Schema)

			generated, err := os.ReadFile("testdata/gateway-generate.graphql")
			require.Nil(t, err)

			require.Equal(t, string(testGatewayExpectedSchema), string(generated))
		})
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}

func TestGraphqlSchemaWithoutUnboundMethod(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Grpc: config.Grpc{
			Services: []kgrpc.Config{
				{
					Target: infos.ConstructsServerAddr.Addr().String(),
				},
				{
					Target: infos.OptionsServerAddr.Addr().String(),
				},
			},
		},
		Server: config.ServerConfig{
			GraphQL: config.GraphQL{
				Playground:             true,
				GenerateUnboundMethods: false,
			},
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		t.Run("schema is correct", func(t *testing.T) {
			gatewayUrl := test.SetupGateway(t, s)
			schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
			if err != nil {
				t.Fatal(err)
			}

			file, err := os.Create("testdata/gateway-generate-without-unbound-method.graphql")
			require.Nil(t, err)
			formatter.NewFormatter(file).FormatSchema(schema.Schema)

			generated, err := os.ReadFile("testdata/gateway-generate-without-unbound-method.graphql")
			require.Nil(t, err)

			require.Equal(t, string(testGatewayExpectedSchemaWithoutUnboundMethod), string(generated))
		})

		t.Run("schema request", func(t *testing.T) {
			gatewayUrl := test.SetupGateway(t, s)
			schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
			require.Nil(t, err)

			schemaDefine := new(bytes.Buffer)
			formatter.NewFormatter(schemaDefine).FormatSchema(schema.Schema)
			_, err = validator.LoadSchema(validator.Prelude, &ast.Source{
				Input: schemaDefine.String(),
			})
			require.Nil(t, err)
		})
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}
