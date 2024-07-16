package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-kod/kod"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nautilus/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/pkg/server"
	"github.com/sysulq/graphql-grpc-gateway/test"
	"go.uber.org/mock/gomock"
)

func TestJwt(t *testing.T) {
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
		Jwt: server.Jwt{
			Enable:               true,
			LocalJwks:            "key",
			ForwardPayloadHeader: "x-jwt-payload",
		},
	}).AnyTimes()

	kod.RunTest(t, func(ctx context.Context, s server.Gateway) {
		gatewayUrl := test.SetupGateway(t, s)
		querier := graphql.NewSingleRequestQueryer(gatewayUrl)

		t.Run("jwt auth", func(t *testing.T) {
			recv := map[string]interface{}{}
			querier.WithMiddlewares([]graphql.NetworkMiddleware{
				func(r *http.Request) error {
					token, err := createToken("bob", mockConfig.Config().Jwt.LocalJwks)
					require.Nil(t, err)

					r.Header.Set("Authorization", "Bearer "+token)
					return nil
				},
			})

			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv); err != nil {
				t.Fatal(err)
			}
			assert.EqualValues(t, constructsMultipleResponse, recv)
		})

		t.Run("jwt auth failed", func(t *testing.T) {
			recv := map[string]interface{}{}
			querier.WithMiddlewares([]graphql.NetworkMiddleware{
				func(r *http.Request) error {
					token, err := createToken("bob", "wrong key")
					require.Nil(t, err)

					r.Header.Set("Authorization", "Bearer "+token)
					return nil
				},
			})

			if err := querier.Query(context.Background(), &graphql.QueryInput{
				Query: contructsMultipleQuery,
			}, &recv); err != nil {
				require.NotNil(t, err)
			}
		})
	}, kod.WithFakes(kod.Fake[server.Config](mockConfig)))
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
