package server

import (
	"context"
	"flag"
	"os"

	"dario.cat/mergo"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kpyroscope"
	"github.com/rs/cors"
	"gopkg.in/yaml.v3"
)

type Tls struct {
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
type Config struct {
	kod.Implements[ConfigComponent]

	Pyroscope  *kpyroscope.Config `json:"pyroscope" yaml:"pyroscope"`
	Grpc       *Grpc              `json:"grpc" yaml:"grpc"`
	Cors       *cors.Options      `json:"cors" yaml:"cors"`
	Playground *bool              `json:"playground" yaml:"playground"`
	Address    string             `json:"address" yaml:"address"`
	Tls        *Tls               `json:"tls" yaml:"tls"`
}

func DefaultConfig() *Config {
	return &Config{
		Address: ":8080",
		Pyroscope: &kpyroscope.Config{
			ServerAddress: "http://localhost:4040",
		},
		Cors:       &cors.Options{},
		Grpc:       &Grpc{},
		Playground: &[]bool{true}[0],
		Tls:        nil,
	}
}

var configFile = flag.String("config", "", "The config file (if not set will use the default configuration)")

func (c *Config) Init(ctx context.Context) error {
	flag.Parse()

	cfg := DefaultConfig()
	if *configFile != "" {
		f, err := os.Open(*configFile)
		if err != nil {
			return err
		}

		err = yaml.NewDecoder(f).Decode(cfg)
		if err != nil {
			return err
		}

		mergo.Merge(c, cfg)
	}

	return nil
}

func (c *Config) Config() *Config {
	return c
}
