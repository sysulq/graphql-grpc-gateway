package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kgrpc"
	"github.com/nautilus/graphql"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestReflectionExit(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{},
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
		GraphQL: config.GraphQL{
			GenerateUnboundMethods: true,
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		querier := graphql.NewSingleRequestQueryer(gatewayUrl)

		t.Run("stop do not panic", func(t *testing.T) {
			infos.OptionServer.Stop()

			recv := map[string]interface{}{}
			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: constructsAnyQuery,
			}, &recv); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(recv, constructsAnyResponse) {
				t.Errorf("mutation failed: expected: %s got: %s", constructsAnyResponse, recv)
			}
		})
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}
