package server

import (
	"context"

	"github.com/fullstorydev/grpcurl"
	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/sysulq/graphql-grpc-gateway/pkg/protojson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type invoker struct {
	kod.Implements[Invoker]
}

func (i *invoker) Invoke(ctx context.Context, c *gin.Context, upstream upstreamInfo, rpcPath string) {
	parser, err := protojson.NewRequestParser(c, upstream.resovler)
	if err != nil {
		_ = c.Error(err)

		return
	}

	handler := protojson.NewEventHandler(c.Writer, upstream.resovler)

	err = grpcurl.InvokeRPC(ctx, upstream.source, upstream.conn,
		rpcPath,
		protojson.ProcessHeaders(c.Request.Header),
		handler, parser.Next)
	if err != nil {
		st := status.Convert(err)
		c.JSON(500, gin.H{"error": st.Code(), "message": st.Message()})
		return
	}

	st := handler.Status
	if st.Code() != codes.OK {
		c.JSON(200, gin.H{"error": st.Code(), "message": st.Message()})
	}
}

func (invoker) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kaccesslog.Interceptor(),
		kmetric.Interceptor(),
		ktrace.Interceptor(),
	}
}
