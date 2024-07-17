package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-kod/kod"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestJwt(t *testing.T) {
	infos := test.SetupDeps(t)

	mockConfig := config.NewMockConfig(gomock.NewController(t))
	mockConfig.EXPECT().Config().Return(&config.ConfigInfo{
		Engine: config.EngineConfig{
			GenerateUnboundMethods: true,
		},
		Grpc: config.Grpc{
			Services: []*config.Service{
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
		Jwt: config.Jwt{
			Enable:               true,
			LocalJwks:            "key",
			ForwardPayloadHeader: "x-jwt-payload",
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)

		t.Run("jwt auth", func(t *testing.T) {
			querier := graphql.NewSingleRequestQueryer(gatewayUrl)

			recv := map[string]interface{}{}
			querier.WithMiddlewares([]graphql.NetworkMiddleware{
				func(r *http.Request) error {
					token, err := createToken("bob", mockConfig.Config().Jwt.LocalJwks)
					require.Nil(t, err)

					r.Header.Set("Authorization", "Bearer "+token)
					return nil
				},
			})

			err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv)
			require.Nil(t, err)
			require.EqualValues(t, constructsMultipleResponse, recv)
		})

		t.Run("jwt auth failed", func(t *testing.T) {
			querier := graphql.NewSingleRequestQueryer(gatewayUrl)
			recv := map[string]interface{}{}
			querier.WithMiddlewares([]graphql.NetworkMiddleware{
				func(r *http.Request) error {
					token, err := createToken("bob", "wrong key")
					require.Nil(t, err)

					r.Header.Set("Authorization", "Bearer "+token)
					return nil
				},
			})

			err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv)
			require.NotNil(t, err)
		})

		t.Run("not authorization", func(t *testing.T) {
			querier := graphql.NewSingleRequestQueryer(gatewayUrl)
			recv := map[string]interface{}{}

			err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv)
			require.NotNil(t, err)
		})

		t.Run("expired authorization", func(t *testing.T) {
			querier := graphql.NewSingleRequestQueryer(gatewayUrl)
			recv := map[string]interface{}{}
			querier.WithMiddlewares([]graphql.NetworkMiddleware{
				func(r *http.Request) error {
					token, err := createExpiredToken("bob", mockConfig.Config().Jwt.LocalJwks)
					require.Nil(t, err)

					r.Header.Set("Authorization", "Bearer "+token)
					return nil
				},
			})

			err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv)
			require.NotNil(t, err)
			require.ErrorContains(t, err, "response was not successful with status code: 401")
		})
	}, kod.WithFakes(kod.Fake[config.Config](mockConfig)))
}

func createToken(username string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createExpiredToken(username string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * -24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
