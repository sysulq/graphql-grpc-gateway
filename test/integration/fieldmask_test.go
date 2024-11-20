package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/nautilus/graphql"
	apitest "github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestFieldMask(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{
			RateLimit:      true,
			CircuitBreaker: true,
		},
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
			GraphQL: config.GraphQL{
				Playground:             true,
				GenerateUnboundMethods: true,
				SingleFlight:           true,
				QueryCache:             true,
			},
		},
	}).AnyTimes()

	scalars := &apitest.Scalars{
		Double: 1.1, Float: 2.2, Int32: 3, Int64: -4, Uint32: 5, Uint64: 6, Sint32: 7, Sint64: 8, Fixed32: 9, Fixed64: 10, Sfixed32: 11,
		Sfixed64: 12, Bool: true, StringX: "test", Bytes: []byte("test"), Paths: &fieldmaskpb.FieldMask{Paths: []string{"double", "float"}},
	}
	mockCaller := server.NewMockCaller(gomock.NewController(t))
	mockCaller.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Cond(func(x any) bool {
		d1, _ := protojson.Marshal(x.(*dynamicpb.Message))
		d2, _ := protojson.Marshal(scalars)
		return string(d1) == string(d2)
	})).Return(scalars, nil).Times(1)

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		querier := graphql.NewSingleRequestQueryer(gatewayUrl)

		cases := []struct {
			name         string
			query        string
			wantResponse interface{}
		}{
			{
				name: "Mutation constructs scalars",
				query: `mutation {
  constructsScalars_(in: {
    double: 1.1, float: 2.2, int32: 3, int64: -4, 
    uint32: 5, uint64: 6, sint32: 7, sint64: 8, 
    fixed32: 9, fixed64: 10, sfixed32: 11, sfixed64: 12,
    bool: true, stringX: "test", bytes: "dGVzdA=="
  }) {
    double
    float
  }
}`,
				wantResponse: j{
					"constructsScalars_": j{
						"double": 1.1,
						"float":  2.2,
					},
				},
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
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig), kod.Fake[server.Caller](mockCaller)))
}
