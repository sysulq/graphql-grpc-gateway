package main

import (
	"context"

	pb "github.com/go-kod/grpc-gateway/api/example/optionsserver"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/go-kod/kod-ext/server/kgrpc"
	"github.com/samber/lo"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.MustRun(context.Background(), func(ctx context.Context, app *app) error {
		etcd := lo.Must(etcdv3.Config{Endpoints: []string{"localhost:2379"}}.Build(ctx))

		s := kgrpc.Config{
			Address: ":8082",
		}.Build().WithRegistry(etcd)
		pb.RegisterServiceServer(s, &service{})

		return s.Run(ctx)
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
