package test

import (
	"context"
	"testing"

	"github.com/nautilus/graphql"
)

func BenchmarkGateway(b *testing.B) {
	gatewayUrl := setupGateway()

	querier := graphql.NewSingleRequestQueryer(gatewayUrl)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		querier.Query(context.Background(), &graphql.QueryInput{
			Query: `{\n  serviceQuery1(in: {\n    string: \"sdfs\",\n    foo: {\n      param1: \"para\"\n    },\n    string2: \"ssdfsd\",\n    float: 111\n  }) {\n    foo {\n      param1\n    }\n  }`,
		}, nil)
	}
}
