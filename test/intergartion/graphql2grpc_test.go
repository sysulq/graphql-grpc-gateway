package intergartion

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/go-kod/kod"
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

		for _, c := range cases {

			handler, err := s.BuildServer(ctx)
			require.Nil(t, err)

			record := httptest.NewRecorder()

			payload, err := json.Marshal(map[string]interface{}{
				"query":         c.query,
				"variables":     map[string]interface{}{},
				"operationName": "",
			})
			require.Nil(t, err)

			req := httptest.NewRequest("POST", "/query", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(record, req)

			data := make(map[string]interface{})
			err = json.Unmarshal(record.Body.Bytes(), &data)
			require.Nil(t, err)

			require.Equal(t, c.wantResponse, data["data"])
		}
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
