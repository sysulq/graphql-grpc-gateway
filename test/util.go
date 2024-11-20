package test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/example/helloworld"
	pb "github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/sysulq/graphql-grpc-gateway/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type DepsInfo struct {
	OptionsServerAddr    net.Listener
	ConstructsServerAddr net.Listener
	HelloworldServerAddr net.Listener
	OptionServer         *grpc.Server
	ConstructServer      *grpc.Server
	HelloworldServer     *grpc.Server
}

func SetupGateway(t testing.TB, s server.Gateway) string {
	serverCh := make(chan net.Addr)
	go func() {
		l, err := net.Listen("tcp", ":0")
		require.Nil(t, err)
		go func() {
			time.Sleep(10 * time.Millisecond)
			serverCh <- l.Addr()
		}()

		_, _ = s.BuildHTTPServer()
		handler, err := s.BuildServer()
		require.Nil(t, err)
		_ = http.Serve(l, handler)
	}()
	gatewayServer := <-serverCh

	gatewayUrl := fmt.Sprintf("http://127.0.0.1:%d/query", gatewayServer.(*net.TCPAddr).Port)
	return gatewayUrl
}

func SetupDeps(t testing.TB) DepsInfo {
	var (
		optionsServerCh    = make(chan net.Listener)
		constructsServerCh = make(chan net.Listener)
		HelloworldServerCh = make(chan net.Listener)

		optionsServer    atomic.Pointer[*grpc.Server]
		constructsServer atomic.Pointer[*grpc.Server]
		helloworldServer atomic.Pointer[*grpc.Server]
	)
	go func() {
		l, err := net.Listen("tcp", "localhost:0")
		require.Nil(t, err)
		go func() {
			time.Sleep(10 * time.Millisecond)
			optionsServerCh <- l
		}()
		s := grpc.NewServer()
		pb.RegisterServiceServer(s, &optionsServiceMock{})
		pb.RegisterQueryServer(s, &optionsQueryMock{})
		reflection.Register(s)
		optionsServer.Store(&s)
		_ = s.Serve(l)
	}()
	go func() {
		l, err := net.Listen("tcp", "localhost:0")
		require.Nil(t, err)
		go func() {
			time.Sleep(10 * time.Millisecond)
			constructsServerCh <- l
		}()
		s := grpc.NewServer()
		pb.RegisterConstructsServer(s, constructsServiceMock{})
		reflection.Register(s)
		constructsServer.Store(&s)
		_ = s.Serve(l)
	}()
	go func() {
		l, err := net.Listen("tcp", "localhost:0")
		require.Nil(t, err)
		go func() {
			time.Sleep(10 * time.Millisecond)
			HelloworldServerCh <- l
		}()
		s := grpc.NewServer()
		helloworld.RegisterGreeterServer(s, &helloworldService{})
		reflection.Register(s)
		helloworldServer.Store(&s)
		_ = s.Serve(l)
	}()

	return DepsInfo{
		OptionsServerAddr:    <-optionsServerCh,
		ConstructsServerAddr: <-constructsServerCh,
		HelloworldServerAddr: <-HelloworldServerCh,
		OptionServer:         *optionsServer.Load(),
		ConstructServer:      *constructsServer.Load(),
		HelloworldServer:     *helloworldServer.Load(),
	}
}

type constructsServiceMock struct {
	pb.ConstructsServer
}

// FilterMessageFields 根据 FieldMask 过滤并返回指定字段的消息
func FilterMessageFields(original protoiface.MessageV1, mask *fieldmaskpb.FieldMask) protoiface.MessageV1 {
	if mask == nil {
		return original
	}

	// 创建一个新的空消息，用于存储筛选后的字段
	filtered := proto.Clone(protoadapt.MessageV2Of(original))
	filteredReflect := filtered.ProtoReflect()

	maskMap := make(map[string]struct{})
	for _, path := range mask.GetPaths() {
		maskMap[path] = struct{}{}
	}

	// 清空不在mask的字段
	filteredReflect.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		if _, ok := maskMap[fd.JSONName()]; !ok {
			filteredReflect.Clear(fd)
		}
		return true
	})

	return protoadapt.MessageV1Of(filtered)
}

func (c constructsServiceMock) Scalars_(ctx context.Context, scalars *pb.Scalars) (*pb.Scalars, error) {
	newScalers := FilterMessageFields(scalars, scalars.GetPaths())
	return newScalers.(*pb.Scalars), nil
}

func (c constructsServiceMock) Repeated_(ctx context.Context, repeated *pb.Repeated) (*pb.Repeated, error) {
	return repeated, nil
}

func (c constructsServiceMock) Maps_(ctx context.Context, maps *pb.Maps) (*pb.Maps, error) {
	return maps, nil
}

func (c constructsServiceMock) Any_(ctx context.Context, any *anypb.Any) (*anypb.Any, error) {
	return any, nil
}

func (c constructsServiceMock) Empty_(ctx context.Context, empty *emptypb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (c constructsServiceMock) Empty2_(ctx context.Context, recursive *pb.EmptyRecursive) (*pb.EmptyNested, error) {
	return &pb.EmptyNested{}, nil
}

func (c constructsServiceMock) Empty3_(ctx context.Context, empty3 *pb.Empty3) (*pb.Empty3, error) {
	return &pb.Empty3{}, nil
}

func (c constructsServiceMock) Ref_(ctx context.Context, ref *pb.Ref) (*pb.Ref, error) {
	return &pb.Ref{}, nil
}

func (c constructsServiceMock) Oneof_(ctx context.Context, oneof *pb.Oneof) (*pb.Oneof, error) {
	return oneof, nil
}

func (c constructsServiceMock) CallWithId(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

type optionsServiceMock struct {
	pb.ServiceServer
}

func (o *optionsServiceMock) Mutate1(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (o *optionsServiceMock) Mutate2(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (o *optionsServiceMock) Query1(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (o *optionsServiceMock) Ignore(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

func (o *optionsServiceMock) Name(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	return data, nil
}

type optionsQueryMock struct {
	pb.UnimplementedQueryServer
}

func (o *optionsQueryMock) Query1(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	fmt.Println(metadata.FromIncomingContext(ctx))
	return data, nil
}

func (o *optionsQueryMock) Query2(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	fmt.Println(metadata.FromIncomingContext(ctx))
	return data, nil
}

type helloworldService struct {
	helloworld.UnimplementedGreeterServer
}

func (s *helloworldService) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if req.Name == "error" {
		return nil, fmt.Errorf("error")
	}

	return &helloworld.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}
