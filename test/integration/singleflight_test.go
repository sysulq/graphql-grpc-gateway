package integration

import (
	"context"
	"sync"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kgrpc"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestSingleFlight(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{
			GenerateUnboundMethods: true,
			SingleFlight:           true,
		},
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
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		querier := graphql.NewSingleRequestQueryer(gatewayUrl)

		t.Run("multiple different query", func(t *testing.T) {
			recv := map[string]interface{}{}
			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv); err != nil {
				t.Fatal(err)
			}
			assert.EqualValues(t, constructsMultipleResponse, recv)
		})

		t.Run("multiple same query", func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(2)
			for i := 0; i < 2; i++ {
				go func(t *testing.T) {
					defer wg.Done()
					recv := map[string]interface{}{}
					err := querier.Query(context.Background(), &graphql.QueryInput{
						Query: contructsMultipleSameQuery,
					}, &recv)
					assert.Nil(t, err)
					assert.EqualValues(t, constructsMultipleSameResponse, recv)
				}(t)
			}
			wg.Wait()
		})
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}
