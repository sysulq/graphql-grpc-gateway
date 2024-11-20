package server

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/sysulq/graphql-grpc-gateway/internal/config"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type graphqlCaller struct {
	kod.Implements[GraphqlCaller]

	config   kod.Ref[config.Config]
	registry kod.Ref[CallerRegistry]

	singleflight singleflight.Group
}

func (c *graphqlCaller) Init(ctx context.Context) error {
	return nil
}

func (c *graphqlCaller) Call(ctx context.Context, rpc protoreflect.MethodDescriptor, message proto.Message) (proto.Message, error) {
	if c.config.Get().Config().Server.GraphQL.SingleFlight {
		if enable, ok := ctx.Value(allowSingleFlightKey).(bool); ok && enable {
			hash := Hash64.Get()
			defer Hash64.Put(hash)

			md, ok := metadata.FromOutgoingContext(ctx)
			if ok {
				hd := make([]string, 0, len(md))
				for k, v := range md {
					// skip grpc gateway prefixed metadata
					if strings.Contains(k, MetadataPrefix) {
						continue
					}
					hd = append(hd, k+strings.Join(v, ","))
				}
				slices.Sort(hd)
				for _, v := range hd {
					_, err := hash.Write([]byte(v))
					if err != nil {
						return nil, err
					}
				}
			}

			msg, err := proto.Marshal(message)
			if err != nil {
				return nil, err
			}

			// generate hash based on rpc pointer
			_, err = hash.Write([]byte(rpc.FullName()))
			if err != nil {
				return nil, err
			}
			_, err = hash.Write(msg)
			if err != nil {
				return nil, err
			}
			sum := hash.Sum64()
			key := strconv.FormatUint(sum, 10)

			res, err, _ := c.singleflight.Do(key, func() (interface{}, error) {
				return c.registry.Get().GetCallerStub(string(rpc.Parent().FullName())).InvokeRpc(ctx, rpc, message)
			})
			if err != nil {
				return nil, err
			}

			return res.(proto.Message), nil
		}
	}

	res, err := c.registry.Get().GetCallerStub(string(rpc.Parent().FullName())).InvokeRpc(ctx, rpc, message)
	return res, err
}

func (c *graphqlCaller) Interceptors() []interceptor.Interceptor {
	if c.config.Get().Config().Engine.CircuitBreaker {
		return []interceptor.Interceptor{
			kcircuitbreaker.Interceptor(),
		}
	}
	return nil
}

type allowSingleFlightType struct{}

var allowSingleFlightKey allowSingleFlightType
