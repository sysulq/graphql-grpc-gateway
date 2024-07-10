package integration

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

func BenchmarkGateway(b *testing.B) {
	b.Setenv("OTEL_TRACES_SAMPLER_ARG", "0.001")
	b.Setenv("OTEL_TRACES_EXPORTER", "none")
	b.Setenv("OTEL_METRICS_EXPORTER", "none")
	b.Setenv("OTEL_LOGS_EXPORTER", "none")

	infos := test.SetupDeps(b)

	mockConfig := server.NewMockConfigComponent(gomock.NewController(b))
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

	kod.RunTest(b, func(ctx context.Context, s server.ServerComponent) {
		handler, err := s.BuildServer()
		require.Nil(b, err)
		payload, err := json.Marshal(map[string]interface{}{
			"query":         constructsScalarsQuery,
			"variables":     map[string]interface{}{},
			"operationName": "",
		})
		require.Nil(b, err)
		record := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/query", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			handler.ServeHTTP(record, req)
			record.Flush()
		}
	}, kod.WithFakes(kod.Fake[server.ConfigComponent](mockConfig)))
}
