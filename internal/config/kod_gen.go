// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package config

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/sysulq/graphql-grpc-gateway/internal/config/Config",
		Interface: reflect.TypeOf((*Config)(nil)).Elem(),
		Impl:      reflect.TypeOf(config{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return config_local_stub{
				impl:        info.Impl.(Config),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[Config] = (*config)(nil)

// Local stub implementations.

type config_local_stub struct {
	impl        Config
	name        string
	interceptor interceptor.Interceptor
}

// Check that config_local_stub implements the Config interface.
var _ Config = (*config_local_stub)(nil)

func (s config_local_stub) Config() (r0 *ConfigInfo) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	r0 = s.impl.Config()
	return
}
