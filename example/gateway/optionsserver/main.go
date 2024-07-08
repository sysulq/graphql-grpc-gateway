package main

import (
	"context"
	"log"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-kod/kod"
	pb "github.com/sysulq/graphql-gateway/example/gateway/optionsserver/pb"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		l, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatal(err)
		}

		s := grpc.NewServer(
			grpc.ServerOption(
				grpc.StatsHandler(otelgrpc.NewServerHandler()),
			),
		)
		pb.RegisterServiceServer(s, &service{})
		reflection.Register(s)
		log.Fatal(s.Serve(l))

		return nil
	})
}

type service struct{ pb.UnimplementedServiceServer }

func (s service) Mutate1(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (s service) Query2(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (s service) Mutate2(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (s service) Query1(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (s service) Publish(server pb.Service_PublishServer) error {
	panic("implement me")
}

func (s service) Subscribe(data *pb.Data, server pb.Service_SubscribeServer) error {
	panic("implement me")
}

func (s service) PubSub1(server pb.Service_PubSub1Server) error {
	panic("implement me")
}

func (s service) InvalidSubscribe1(server pb.Service_InvalidSubscribe1Server) error {
	panic("implement me")
}

func (s service) InvalidSubscribe2(data *pb.Data, server pb.Service_InvalidSubscribe2Server) error {
	panic("implement me")
}

func (s service) InvalidSubscribe3(server pb.Service_InvalidSubscribe3Server) error {
	panic("implement me")
}

func (s service) PubSub2(server pb.Service_PubSub2Server) error {
	panic("implement me")
}
