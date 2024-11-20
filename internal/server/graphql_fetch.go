package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kod/kod"
	"github.com/jhump/protoreflect/v2/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func isReflectionServiceName(name string) bool {
	return name == grpc_reflection_v1.ServerReflection_ServiceDesc.ServiceName ||
		name == grpc_reflection_v1alpha.ServerReflection_ServiceDesc.ServiceName
}

var ErrTLSHandshakeFailed = errors.New("TLS handshake failed")

type reflection struct {
	kod.Implements[Reflection]
}

func (ins *reflection) ListPackages(ctx context.Context, cc grpc.ClientConnInterface) ([]protoreflect.FileDescriptor, error) {
	client := grpcreflect.NewClientAuto(ctx, cc)
	ssvcs, err := client.ListServices()
	if err != nil {
		msg := status.Convert(err).Message()
		// Check whether the error message contains TLS related error.
		// If the server didn't enable TLS, the error message contains the first string.
		// If Evans didn't enable TLS against to the TLS enabled server, the error message contains
		// the second string.
		if strings.Contains(msg, "tls: first record does not look like a TLS handshake") ||
			strings.Contains(msg, "latest connection error: <nil>") {
			return nil, ErrTLSHandshakeFailed
		}
		return nil, fmt.Errorf("failed to list services from reflecton enabled gRPC server: %w", err)
	}

	var fds []protoreflect.FileDescriptor
	for _, s := range ssvcs {
		if isReflectionServiceName(string(s)) {
			continue
		}
		svc, err := client.AsResolver().FindDescriptorByName(s)
		if err != nil {
			return nil, err
		}

		fd := svc.ParentFile() //.AsFileDescriptorProto()
		fds = append(fds, fd)
	}
	return fds, nil
}
