package config

import (
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kgrpc"
	"github.com/go-kod/kod/ext/client/kpyroscope"
	"github.com/go-kod/kod/ext/registry/etcdv3"
)

type config struct {
	kod.Implements[Config]
	kod.WithGlobalConfig[ConfigInfo]
}

type Pyroscope struct {
	kpyroscope.Config `mapstructure:",squash"`
	Enable            bool
}

type GraphQL struct {
	Address    string
	Disable    bool
	Playground bool
	Jwt        Jwt
}

type Jwt struct {
	Enable               bool
	LocalJwks            string
	ForwardPayloadHeader string
}

type JwtClaimToHeader struct {
	HeaderName string
	ClaimName  string
}

type EngineConfig struct {
	GenerateUnboundMethods bool
	Pyroscope              Pyroscope
	RateLimit              bool
	CircuitBreaker         bool
	QueryCache             bool
	SingleFlight           bool
}

type ConfigInfo struct {
	Engine  EngineConfig
	Grpc    Grpc
	GraphQL GraphQL
}

type Grpc struct {
	Etcd     etcdv3.Config
	Services []kgrpc.Config
}

func (ins *config) Config() *ConfigInfo {
	return ins.WithGlobalConfig.Config()
}
