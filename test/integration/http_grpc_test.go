package integration

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/client/kgrpc"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/stretchr/testify/assert"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestHTTP2Grpc(t *testing.T) {
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
				{
					Target: infos.HelloworldServerAddr.Addr().String(),
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

	t.Run("http to grpc", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, up server.HttpUpstream) {
			router := http.NewServeMux()
			up.Register(ctx, router)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/say/bob", nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{\"message\":\"Hello bob\"}", rec.Body.String())
		}, kod.WithFakes(kod.Fake[config.Config](mockConfig)), kod.WithOpenTelemetryDisabled())
	})

	t.Run("not found", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, up server.HttpUpstream) {
			router := http.NewServeMux()
			up.Register(ctx, router)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/say-notfound", nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, "404 page not found\n", rec.Body.String())
		}, kod.WithFakes(kod.Fake[config.Config](mockConfig)), kod.WithOpenTelemetryDisabled())
	})

	t.Run("error", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, up server.HttpUpstream) {
			router := http.NewServeMux()
			up.Register(ctx, router)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/say/error", nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, `{"code":2,"message":"error","details":[]}`, rec.Body.String())
		}, kod.WithFakes(kod.Fake[config.Config](mockConfig)), kod.WithOpenTelemetryDisabled())
	})

	t.Run("http body", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, up server.HttpUpstream) {
			router := http.NewServeMux()
			up.Register(ctx, router)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/say/sam", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, `{"message":"Hello sam"}`, rec.Body.String())
		}, kod.WithFakes(kod.Fake[config.Config](mockConfig)), kod.WithOpenTelemetryDisabled())
	})

	t.Run("invalid http body", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, up server.HttpUpstream) {
			router := http.NewServeMux()
			up.Register(ctx, router)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/say/sam", bytes.NewBufferString("{invalid data}"))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `invalid character 'i' looking for beginning of object key string`, rec.Body.String())
		}, kod.WithFakes(kod.Fake[config.Config](mockConfig)), kod.WithOpenTelemetryDisabled())
	})
}
