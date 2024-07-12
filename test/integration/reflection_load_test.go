package integration

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/sysulq/graphql-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestReflectionLoad(t *testing.T) {
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
	}).AnyTimes()

	kod.RunTest2(t, func(ctx context.Context, s server.ServerComponent, reg server.Registry) {
		_ = test.SetupGateway(t, s)

		t.Run("stop do not panic", func(t *testing.T) {
			// pp.Println(reg.SchemaDescriptorList().AsGraphql())
			t.Fail()
			// querier := graphql.NewSingleRequestQueryer(gatewayUrl)
		})
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
