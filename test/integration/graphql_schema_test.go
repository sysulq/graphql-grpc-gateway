package integration

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/go-kod/kod"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"github.com/vektah/gqlparser/v2/formatter"
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
		Engine: config.EngineConfig{
			GenerateUnboundMethods: true,
		},
		Grpc: config.Grpc{
			Services: []*config.Service{
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
		GraphQL: config.GraphQL{
			Playground: true,
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
		Engine: config.EngineConfig{
			GenerateUnboundMethods: false,
		},
		Grpc: config.Grpc{
			Services: []*config.Service{
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
		GraphQL: config.GraphQL{
			Playground: true,
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		t.Run("schema is correct", func(t *testing.T) {
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
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}
