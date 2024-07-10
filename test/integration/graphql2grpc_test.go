package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-kod/kod"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/sysulq/graphql-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestGraphql2Grpc(t *testing.T) {
	infos := test.SetupDeps()

	mockConfig := server.NewMockConfigComponent(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&server.Config{
		Grpc: server.Grpc{
			Services: []*server.Service{
				{
					Address:    infos.ConstructsServerAddr,
					Reflection: true,
				},
				{
					Address:    infos.OptionsServerAddr,
					Reflection: true,
				},
			},
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.ServerComponent) {
		serverCh := make(chan net.Addr)
		go func() {
			l, err := net.Listen("tcp", ":0")
			require.Nil(t, err)
			go func() {
				time.Sleep(10 * time.Millisecond)
				serverCh <- l.Addr()
			}()

			handler, err := s.BuildServer()
			require.Nil(t, err)
			http.Serve(l, handler)
		}()
		gatewayServer := <-serverCh

		gatewayUrl := fmt.Sprintf("http://127.0.0.1:%d/query", gatewayServer.(*net.TCPAddr).Port)
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
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
