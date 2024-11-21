package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-kod/grpc-gateway/internal/config"
	"github.com/go-kod/grpc-gateway/pkg/header"
	"github.com/go-kod/grpc-gateway/pkg/protojson"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/ktrace"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpUpstreamInvoker struct {
	kod.Implements[HttpUpstreamInvoker]

	singleflight singleflight.Group
	config       kod.Ref[config.Config]
}

func (i *httpUpstreamInvoker) Invoke(ctx context.Context, rw http.ResponseWriter, r *http.Request, upstream upstreamInfo, rpcPath string, pathNames []string) {
	eh := &protojson.EventHandler{}
	err := error(nil)

	if i.config.Get().Config().Server.HTTP.SingleFlight {
		if r.Method == http.MethodGet {
			hash := Hash64.Get()
			defer Hash64.Put(hash)

			_, _ = hash.WriteString(r.URL.String())

			for k, v := range r.Header {
				_, _ = hash.WriteString(k)
				_, _ = hash.WriteString(strings.Join(v, ","))
			}

			key := strconv.FormatUint(hash.Sum64(), 10)

			handler, err1, _ := i.singleflight.Do(key, func() (interface{}, error) {
				return i.invoke(ctx, rw, r, upstream, rpcPath, pathNames)
			})

			eh = handler.(*protojson.EventHandler)
			err = err1
		}
	} else {
		eh, err = i.invoke(ctx, rw, r, upstream, rpcPath, pathNames)
	}

	if err != nil {
		http.Error(rw, err.Error(), codeFromGrpcError(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	if eh.Status != nil {
		err = eh.Marshaler.Marshal(rw, eh.Status.Proto())
	} else {
		err = eh.Marshaler.Marshal(rw, eh.Message)
	}

	if err != nil {
		i.L(ctx).Error("marshal response", "error", err)
		http.Error(rw, err.Error(), codeFromGrpcError(err))
	}
}

func (i *httpUpstreamInvoker) invoke(ctx context.Context,
	rw http.ResponseWriter, r *http.Request, upstream upstreamInfo, rpcPath string, pathNames []string,
) (*protojson.EventHandler, error) {
	parser, err := protojson.NewRequestParser(r, pathNames, upstream.resovler)
	if err != nil {
		i.L(ctx).Error("parse request", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	handler := protojson.NewEventHandler(rw, upstream.resovler)

	err = grpcurl.InvokeRPC(ctx, upstream.source, upstream.conn,
		rpcPath,
		header.ProcessHeaders(r.Header),
		handler, parser.Next)
	if err != nil {
		i.L(ctx).Error("invoke rpc", "error", err)
		if handler.Status == nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return handler, nil
}

func (*httpUpstreamInvoker) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kaccesslog.Interceptor(),
		kmetric.Interceptor(),
		ktrace.Interceptor(),
	}
}

func codeFromGrpcError(err error) int {
	code := status.Code(err)
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Internal, codes.DataLoss, codes.Unknown:
		return http.StatusInternalServerError
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	}

	return http.StatusInternalServerError
}
