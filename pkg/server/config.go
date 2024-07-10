package server

import (
	"context"

	"dario.cat/mergo"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kpyroscope"
	"github.com/rs/cors"
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
	kod.Implements[ConfigComponent]
	kod.WithGlobalConfig[Config]
}

type Config struct {
	Pyroscope  kpyroscope.Config `json:"pyroscope" yaml:"pyroscope"`
	Grpc       Grpc              `json:"grpc" yaml:"grpc"`
	Cors       cors.Options      `json:"cors" yaml:"cors"`
	Playground bool              `json:"playground" yaml:"playground"`
	Address    string            `json:"address" yaml:"address"`
	Tls        Tls               `json:"tls" yaml:"tls"`
}

func defaultConfig() *Config {
	return &Config{
		Address: ":8080",
		Pyroscope: kpyroscope.Config{
			ServerAddress: "http://localhost:4040",
		},
		Cors:       cors.Options{},
		Grpc:       Grpc{},
		Playground: true,
	}
}

func (ins *config) Init(ctx context.Context) error {
	return mergo.Merge(ins.Config(), defaultConfig())
}

func (ins *config) Config() *Config {
	return ins.WithGlobalConfig.Config()
}
