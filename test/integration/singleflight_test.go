package integration

import (
	"context"
	"sync"
	"testing"

	"github.com/go-kod/kod"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/sysulq/graphql-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestSingleFlight(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := server.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&server.ConfigInfo{
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
	}, kod.WithFakes(kod.Fake[server.Config](mockConfig)))
}
