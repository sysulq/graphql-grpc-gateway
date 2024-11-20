package main

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/go-kod/kod-ext/server/kgrpc"
	"github.com/samber/lo"
	pb "github.com/sysulq/graphql-grpc-gateway/api/example/helloworld"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.MustRun(context.Background(), func(ctx context.Context, app *app) error {
		etcd := lo.Must(etcdv3.Config{Endpoints: []string{"localhost:2379"}}.Build(ctx))

		s := kgrpc.Config{
			Address: ":8081",
		}.Build().WithRegistry(etcd)
		pb.RegisterGreeterServer(s, &service{})

		return s.Run(ctx)
	})
}

type service struct {
	pb.UnimplementedGreeterServer
}

func (s *service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}
