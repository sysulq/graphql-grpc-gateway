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

func TestGraphql2Grpc(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{
			RateLimit:      true,
			CircuitBreaker: true,
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
		GraphQL: config.GraphQL{
			Playground:             true,
			GenerateUnboundMethods: true,
			SingleFlight:           true,
			QueryCache:             true,
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		querier := graphql.NewSingleRequestQueryer(gatewayUrl)

		cases := []struct {
			name         string
			query        string
			wantResponse interface{}
		}{
			{
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
				name:         "Mutation constructs oneofs",
				query:        constructsOneofsQuery,
				wantResponse: constructsOneofsResponse,
			},
		}

		for _, tc := range cases {

			recv := map[string]interface{}{}
			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: tc.query,
			}, &recv); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(recv, tc.wantResponse) {
				t.Errorf("mutation failed: expected: %s got: %s", tc.wantResponse, recv)
			}
		}
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}
