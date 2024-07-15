package integration

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/go-kod/kod"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/sysulq/graphql-gateway/test"
	"github.com/vektah/gqlparser/v2/formatter"
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
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
