package main

import (
	"context"

	any "google.golang.org/protobuf/types/known/anypb"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/registry/etcdv3"
	"github.com/go-kod/kod-ext/server/kgrpc"
	"github.com/samber/lo"
	pb "github.com/sysulq/graphql-grpc-gateway/api/example/constructsserver"
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
		pb.RegisterConstructsServer(s, &service{})

		return s.Run(ctx)
	})
}

type service struct {
	pb.UnimplementedConstructsServer
}

func (s *service) Anyway_(ctx context.Context, a *pb.Any) (*pb.AnyInput, error) {
	return &pb.AnyInput{
		Any: a.Any,
	}, nil
}

func (s *service) Scalars_(ctx context.Context, scalars *pb.Scalars) (*pb.Scalars, error) {
	return scalars, nil
}

func (s *service) Repeated_(ctx context.Context, repeated *pb.Repeated) (*pb.Repeated, error) {
	return repeated, nil
}

func (s *service) Maps_(ctx context.Context, maps *pb.Maps) (*pb.Maps, error) {
	return maps, nil
}

func (s *service) Any_(ctx context.Context, a *any.Any) (*any.Any, error) {
	println(a.String())
	//_ = json.NewEncoder(os.Stdout).Encode(a)
	//var b pb.Bar
	//_ = ptypes.UnmarshalAny(a, &b)
	//println(b.String())
	return a, nil
}

func (s *service) Empty_(ctx context.Context, empty *empty.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *service) Empty2_(ctx context.Context, recursive *pb.EmptyRecursive) (*pb.EmptyNested, error) {
	panic("implement me")
}

func (s *service) Empty3_(ctx context.Context, empty3 *pb.Empty3) (*pb.Empty3, error) {
	return &pb.Empty3{}, nil
}

func (s *service) Ref_(ctx context.Context, ref *pb.Ref) (*pb.Ref, error) {
	return ref, nil
}

func (s *service) Oneof_(ctx context.Context, oneof *pb.Oneof) (*pb.Oneof, error) {
	return oneof, nil
}

func (s *service) CallWithId(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
