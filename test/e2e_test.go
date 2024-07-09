package test

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-kod/kod"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-gateway/pkg/server"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/mock/gomock"
)

//go:embed testdata/gateway-expect.graphql
var testGatewayExpectedSchema []byte

func TestE2E(t *testing.T) {
	infos := SetupDeps()

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
		Playground: true,
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.ServerComponent) {
		constructsServerCh := make(chan net.Addr)
		go func() {
			l, err := net.Listen("tcp", ":0")
			require.Nil(t, err)
			go func() {
				time.Sleep(10 * time.Millisecond)
				constructsServerCh <- l.Addr()
			}()

			handler, err := s.BuildServer(ctx)
			require.Nil(t, err)
			http.Serve(l, handler)
		}()
		gatewayServer := <-constructsServerCh

		gatewayUrl := fmt.Sprintf("http://127.0.0.1:%d/query", gatewayServer.(*net.TCPAddr).Port)
		t.Run("schema is correct", func(t *testing.T) {
			schema, err := graphql.IntrospectRemoteSchema(gatewayUrl)
			if err != nil {
				t.Fatal(err)
			}
			expectedFormattedSchema, _ := gqlparser.LoadSchema(&ast.Source{Input: string(testGatewayExpectedSchema)})
			compareGraphql(t, schema.Schema, expectedFormattedSchema)
		})
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
