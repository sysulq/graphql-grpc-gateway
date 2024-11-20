package integration

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/pkg/header"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestSingleFlight(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{},
		Grpc: config.Grpc{
			Etcd: etcdv3.Config{
				Endpoints: []string{"localhost:2379"},
			},
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
			GraphQL: config.GraphQLConfig{
				GenerateUnboundMethods: true,
				SingleFlight:           true,
			},
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)

		t.Run("multiple different query", func(t *testing.T) {
			recv := map[string]interface{}{}
			querier := graphql.NewSingleRequestQueryer(gatewayUrl)
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
					querier := graphql.NewSingleRequestQueryer(gatewayUrl)
					querier.WithMiddlewares([]graphql.NetworkMiddleware{
						func(r *http.Request) error {
							r.Header.Set("Authorization", "Bearer ")
							r.Header.Set(header.MetadataHeaderPrefix+"singleflight", "true")
							return nil
						},
					})
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
