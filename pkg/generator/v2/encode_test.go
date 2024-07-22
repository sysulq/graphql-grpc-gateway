package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
)

type (
	j = map[string]interface{}
	l = []interface{}
)

const constructsMapsQuery = `
mutation {
 constructsMaps_(
    in: {
      int32Int32: { key: 1, value: 1 }
      int64Int64: { key: 2, value: 2 }
      uint32Uint32: { key: 3, value: 3 }
      uint64Uint64: { key: 4, value: 4 }
      sint32Sint32: { key: 5, value: 5 }
      sint64Sint64: { key: 6, value: 5 }
      fixed32Fixed32: { key: 7, value: 7 }
      fixed64Fixed64: { key: 8, value: 8 }
      sfixed32Sfixed32: { key: 9, value: 9 }
      sfixed64Sfixed64: { key: 10, value: 10 }
      boolBool: { key: true, value: true }
      stringString: { key: "test1", value: "test1" }
      stringBytes: { key: "test2", value: "dGVzdA==" }
      stringFloat: { key: "test1", value: 11.1 }
      stringDouble: { key: "test1", value: 12.2 }
      stringFoo: { key: "test1", value: { param1: "param1", param2: "param2" } }
      stringBar: { key: "test1", value: BAR3 }
    }
  ) {
    int32Int32 { key value }
    int64Int64 { key value}
    uint32Uint32 { key value }
	uint64Uint64 { key value}
    sint32Sint32 { key value}
    sint64Sint64 { key value}
    fixed32Fixed32 { key value}
    fixed64Fixed64 { key value}
    sfixed32Sfixed32 { key value}
    sfixed64Sfixed64 { key value}
    boolBool { key value}
    stringString { key value}
    stringBytes { key value}
    stringFloat { key value}
    stringDouble { key value}
    stringFoo { key value { param1 param2 }}
    stringBar { key value}
  }
}
`

var constructsMapsResponse = j{
	"constructsMaps_": j{
		"boolBool":         l{j{"key": true, "value": true}},
		"stringBar":        l{j{"key": "test1", "value": "BAR3"}},
		"stringBytes":      l{j{"key": "test2", "value": "dGVzdA=="}},
		"stringDouble":     l{j{"key": "test1", "value": 12.2}},
		"stringFloat":      l{j{"key": "test1", "value": 11.1}},
		"stringFoo":        l{j{"key": "test1", "value": j{"param1": "param1", "param2": "param2"}}},
		"stringString":     l{j{"key": "test1", "value": "test1"}},
		"fixed32Fixed32":   l{j{"key": float64(7), "value": float64(7)}},
		"fixed64Fixed64":   l{j{"key": float64(8), "value": float64(8)}},
		"int32Int32":       l{j{"key": float64(1), "value": float64(1)}},
		"int64Int64":       l{j{"key": float64(2), "value": float64(2)}},
		"sfixed32Sfixed32": l{j{"key": float64(9), "value": float64(9)}},
		"sfixed64Sfixed64": l{j{"key": float64(10), "value": float64(10)}},
		"sint32Sint32":     l{j{"key": float64(5), "value": float64(5)}},
		"sint64Sint64":     l{j{"key": float64(6), "value": float64(5)}},
		"uint32Uint32":     l{j{"key": float64(3), "value": float64(3)}},
		"uint64Uint64":     l{j{"key": float64(4), "value": float64(4)}},
	},
}

func TestEncodeMaps(t *testing.T) {
	ins := &Generator{}

	var msg test.Maps

	generated, err := ins.GraphQL2Proto(msg.ProtoReflect().Type().Descriptor(), &ast.Field{
		Alias: "constructsMaps_",
		Name:  "constructsMaps_",
		Arguments: ast.ArgumentList{
			&ast.Argument{
				Name: "in",
				Value: &ast.Value{
					Children: ast.ChildValueList{
						&ast.ChildValue{
							Name: "int32Int32",
							Value: &ast.Value{
								Children: ast.ChildValueList{
									&ast.ChildValue{
										Name: "",
										Value: &ast.Value{
											Children: ast.ChildValueList{
												&ast.ChildValue{
													Name: "key",
													Value: &ast.Value{
														Raw:  "1",
														Kind: ast.IntValue,
													},
												},
												&ast.ChildValue{
													Name: "value",
													Value: &ast.Value{
														Raw:  "2",
														Kind: ast.IntValue,
													},
												},
											},
											Kind: ast.ObjectValue,
										},
									},
								},
								Kind: ast.ListValue,
							},
						},
						&ast.ChildValue{
							Name: "int64Int64",
							Value: &ast.Value{
								Children: ast.ChildValueList{
									&ast.ChildValue{
										Name: "",
										Value: &ast.Value{
											Children: ast.ChildValueList{
												&ast.ChildValue{
													Name: "key",
													Value: &ast.Value{
														Raw:  "1",
														Kind: ast.IntValue,
													},
												},
												&ast.ChildValue{
													Name: "value",
													Value: &ast.Value{
														Raw:  "2",
														Kind: ast.IntValue,
													},
												},
											},
											Kind: ast.ObjectValue,
										},
									},
									&ast.ChildValue{
										Name: "",
										Value: &ast.Value{
											Children: ast.ChildValueList{
												&ast.ChildValue{
													Name: "key",
													Value: &ast.Value{
														Raw:  "2",
														Kind: ast.IntValue,
													},
												},
												&ast.ChildValue{
													Name: "value",
													Value: &ast.Value{
														Raw:  "3",
														Kind: ast.IntValue,
													},
												},
											},
											Kind: ast.ObjectValue,
										},
									},
								},
								Kind: ast.ListValue,
							},
						},
					},
					Kind: ast.ObjectValue,
				},
			},
		},
	}, map[string]interface{}{})

	require.Nil(t, err)
	d1, _ := proto.Marshal(generated)
	d2, _ := proto.Marshal(&test.Maps{
		Int32Int32: map[int32]int32{1: 2},
		Int64Int64: map[int64]int64{1: 2, 2: 3},
	})
	require.Equal(t, d1, d2)
}
