package config

import (
	"context"

	"dario.cat/mergo"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kpyroscope"
)

type Tls struct {
	Enable      bool   `json:"enable" yaml:"enable"`
	Certificate string `json:"certificate" yaml:"certificate"`
	PrivateKey  string `json:"private_key" yaml:"private_key"`
}

type Authentication struct {
	Tls *Tls `json:"tls" yaml:"tls"`
}

type Service struct {
	Address        string          `json:"address" yaml:"address"`
	Authentication *Authentication `json:"authentication" yaml:"authentication"`
	Reflection     bool            `json:"reflection" yaml:"reflection"`
	ProtoFiles     []string        `json:"proto_files" yaml:"proto_files"`
}

type config struct {
	kod.Implements[Config]
	kod.WithGlobalConfig[ConfigInfo]
}

type Pyroscope struct {
	kpyroscope.Config
	Enable bool
}

type GraphQL struct {
	Disable    bool `json:"disable" yaml:"disable"`
	Playground bool `json:"playground" yaml:"playground"`
}

type Jwt struct {
	Enable               bool   `json:"enable" yaml:"enable"`
	LocalJwks            string `json:"local_jwks" yaml:"local_jwks"`
	ForwardPayloadHeader string `json:"forward_payload_header" yaml:"forward_payload_header"`
}

type JwtClaimToHeader struct {
	HeaderName string `json:"header_name" yaml:"header_name"`
	ClaimName  string `json:"claim_name" yaml:"claim_name"`
}

type EngineConfig struct {
	GenerateUnboundMethods bool      `json:"generate_unbound_methods" yaml:"generate_unbound_methods"`
	Pyroscope              Pyroscope `json:"pyroscope" yaml:"pyroscope"`
}

type ConfigInfo struct {
	Engine  EngineConfig `json:"gateway" yaml:"gateway"`
	Grpc    Grpc         `json:"grpc" yaml:"grpc"`
	Address string       `json:"address" yaml:"address"`
	Tls     Tls          `json:"tls" yaml:"tls"`
	GraphQL GraphQL      `json:"graphql" yaml:"graphql"`
	Jwt     Jwt          `json:"jwt" yaml:"jwt"`
}

type Grpc struct {
	Services    []*Service
	ImportPaths []string
}

func defaultConfig() *ConfigInfo {
	return &ConfigInfo{
		Engine: EngineConfig{
			GenerateUnboundMethods: false,
			Pyroscope: Pyroscope{
				Enable: false,
				Config: kpyroscope.Config{
					ServerAddress: "http://localhost:4040",
				},
			},
		},
		Address: ":8080",
		Grpc:    Grpc{},
		GraphQL: GraphQL{
			Playground: true,
		},
	}
}

func (ins *config) Init(ctx context.Context) error {
	return mergo.Merge(ins.Config(), defaultConfig())
}

func (ins *config) Config() *ConfigInfo {
	return ins.WithGlobalConfig.Config()
}
