package v2

type (
	j  = map[string]interface{}
	jj = []map[string]interface{}
	l  = []interface{}
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

var constructsMapsResponse = j{
	"boolBool":         jj{j{"key": true, "value": true}},
	"stringBar":        jj{j{"key": "test1", "value": "BAR3"}},
	"stringBytes":      jj{j{"key": "test2", "value": "dGVzdA=="}},
	"stringDouble":     jj{j{"key": "test1", "value": 12.2}},
	"stringFloat":      jj{j{"key": "test1", "value": float32(11.1)}},
	"stringFoo":        jj{j{"key": "test1", "value": j{"param1": "param1", "param2": "param2"}}},
	"stringString":     jj{j{"key": "test1", "value": "test1"}},
	"fixed32Fixed32":   jj{j{"key": uint32(7), "value": uint32(7)}},
	"fixed64Fixed64":   jj{j{"key": uint64(8), "value": uint64(8)}},
	"int32Int32":       jj{j{"key": int32(1), "value": int32(1)}},
	"int64Int64":       jj{j{"key": int64(2), "value": int64(2)}},
	"sfixed32Sfixed32": jj{j{"key": int32(9), "value": int32(9)}},
	"sfixed64Sfixed64": jj{j{"key": int64(10), "value": int64(10)}},
	"sint32Sint32":     jj{j{"key": int32(5), "value": int32(5)}},
	"sint64Sint64":     jj{j{"key": int64(6), "value": int64(5)}},
	"uint32Uint32":     jj{j{"key": uint32(3), "value": uint32(3)}},
	"uint64Uint64":     jj{j{"key": uint64(4), "value": uint64(4)}},
}

const constructsRepeatedQuery = `
mutation {
 constructsRepeated_(
    in: {
      double: [1.1]
      float: [2.2]
      int32: [3]
      int64: [4]
      uint32: [7]
      uint64: [8]
      sint32: [9]
      sint64: [10]
      fixed32: [11]
      fixed64: [12]
      sfixed32: [13]
      sfixed64: [14]
      bool: [true]
      stringX: ["test"]
      bytes: ["dGVzdA=="]
      foo: [{ param1: "param1", param2: "param2" }]
      bar: [BAR3]
    }
  ) {
    double
    float
    int32
    int64
    uint32
    uint64
    sint32
    sint64
    fixed32
    fixed64
    sfixed32
    sfixed64
    bool
    stringX
    bytes
    foo {param1 param2}
    bar
  }
}
`

var constructsRepeatedResponse = j{
	"bar":      l{"BAR3"},
	"bool":     l{true},
	"stringX":  l{"test"},
	"bytes":    l{"dGVzdA=="},
	"foo":      l{j{"param1": "param1", "param2": "param2"}},
	"double":   l{1.1},
	"float":    l{float32(2.2)},
	"fixed32":  l{uint32(11)},
	"fixed64":  l{uint64(12)},
	"int32":    l{int32(3)},
	"int64":    l{int64(4)},
	"sfixed32": l{int32(13)},
	"sfixed64": l{int64(14)},
	"sint32":   l{int32(9)},
	"sint64":   l{int64(10)},
	"uint32":   l{uint32(7)},
	"uint64":   l{uint64(8)},
}
