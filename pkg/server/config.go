package server

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

type ConfigInfo struct {
	Pyroscope Pyroscope `json:"pyroscope" yaml:"pyroscope"`
	Grpc      Grpc      `json:"grpc" yaml:"grpc"`
	Address   string    `json:"address" yaml:"address"`
	Tls       Tls       `json:"tls" yaml:"tls"`
	GraphQL   GraphQL   `json:"graphql" yaml:"graphql"`
}

func defaultConfig() *ConfigInfo {
	return &ConfigInfo{
		Address: ":8080",
		Pyroscope: Pyroscope{
			Enable: false,
			Config: kpyroscope.Config{
				ServerAddress: "http://localhost:4040",
			},
		},
		Grpc: Grpc{},
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
