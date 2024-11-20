package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

func addHeader(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md, _ := metadata.FromOutgoingContext(r.Context())
		md = metadata.Join(md, httpHeadersToGRPCMetadata(r.Header))

		ctx := metadata.NewOutgoingContext(r.Context(), md)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ins *server) jwtAuthHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}
		tokenString = tokenString[len("Bearer "):]

		token, err := ins.verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token")
			return
		}

		md, ok := metadata.FromOutgoingContext(r.Context())
		if !ok {
			md = metadata.New(nil)
		}

		forwardPayloadHeader := ins.config.Get().Config().Server.GraphQL.Jwt.ForwardPayloadHeader
		if len(forwardPayloadHeader) > 0 {
			md.Set(forwardPayloadHeader, strings.Split(token.Raw, ".")[1])
		}

		ctx := metadata.NewOutgoingContext(r.Context(), md)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ins *server) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(ins.config.Get().Config().Server.GraphQL.Jwt.LocalJwks), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
