package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	j = map[string]interface{}
	l = []interface{}
)

const constructsMapsQuery = `
mutation {
 constructsMaps_(
    in: {
      int32Int32: [{ key: 1, value: 1 }]
      int64Int64: [{ key: 2, value: 2 }]
      uint32Uint32: [{ key: 3, value: 3 }]
      uint64Uint64: [{ key: 4, value: 4 }]
      sint32Sint32: [{ key: 5, value: 5 }]
      sint64Sint64: [{ key: 6, value: 5 }]
      fixed32Fixed32: [{ key: 7, value: 7 }]
      fixed64Fixed64: [{ key: 8, value: 8 }]
      sfixed32Sfixed32: [{ key: 9, value: 9 }]
      sfixed64Sfixed64: [{ key: 10, value: 10 }]
      boolBool: [{ key: true, value: true }]
      stringString: [{ key: "test1", value: "test1" }]
      stringBytes: [{ key: "test2", value: "dGVzdA==" }]
      stringFloat: [{ key: "test1", value: 11.1 }]
      stringDouble: [{ key: "test1", value: 12.2 }]
      stringFoo: [{ key: "test1", value: { param1: "param1", param2: "param2" } }]
      stringBar: [{ key: "test1", value: BAR3 }]
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

var constructsMaps = &test.Maps{
	Int32Int32: map[int32]int32{1: 1},
	Int64Int64: map[int64]int64{2: 2},
	Uint32Uint32: map[uint32]uint32{
		3: 3,
	},
	Uint64Uint64: map[uint64]uint64{
		4: 4,
	},
	Sint32Sint32: map[int32]int32{
		5: 5,
	},
	Sint64Sint64: map[int64]int64{
		6: 5,
	},
	Fixed32Fixed32: map[uint32]uint32{
		7: 7,
	},
	Fixed64Fixed64: map[uint64]uint64{
		8: 8,
	},
	Sfixed32Sfixed32: map[int32]int32{
		9: 9,
	},
	Sfixed64Sfixed64: map[int64]int64{
		10: 10,
	},
	BoolBool: map[bool]bool{
		true: true,
	},
	StringString: map[string]string{
		"test1": "test1",
	},
	StringBytes: map[string][]byte{
		"test2": []byte("test"),
	},
	StringFloat: map[string]float32{
		"test1": 11.1,
	},
	StringDouble: map[string]float64{
		"test1": 12.2,
	},
	StringFoo: map[string]*test.Foo{
		"test1": {Param1: "param1", Param2: "param2"},
	},
	StringBar: map[string]test.Bar{
		"test1": test.Bar_BAR3,
	},
}

func TestEncodeMaps(t *testing.T) {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: schemaGraphQL,
	})
	require.Nil(t, err)
	var msg test.Maps

	queryDoc, err := gqlparser.LoadQuery(schema, constructsMapsQuery)
	require.Nil(t, err)

	ins := &Generator{}
	generated, err := ins.GraphQL2Proto(msg.ProtoReflect().Descriptor(), queryDoc.Operations[0].SelectionSet[0].(*ast.Field), nil)

	require.Nil(t, err)

	data, err := protojson.Marshal(generated)
	require.Nil(t, err)
	expected, err := protojson.Marshal(constructsMaps)
	require.Nil(t, err)
	require.Equal(t, string(expected), string(data))
}
