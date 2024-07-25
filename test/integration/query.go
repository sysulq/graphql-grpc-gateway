package integration

type (
	j = map[string]interface{}
	l = []interface{}
)

const contructsMultipleQuery = `
{
  queryQuery1(in: {
    stringX: "queryQuery1", 
    foo: {                         
      param1: "Foo Parameter Value"
    },
    double: [1.23, 4.56],           
    string2: "Another String Value"  
    foo2: {                         
      param1: "Foo2 Parameter Value"
    },
    double2: [7.89],                  
    bar: "BAR2",                    
    string: "Yet another string"      
  }) {
   stringX
  },
  queryQuery2(in: {
    stringX: "queryQuery2", 
    foo: {                         
      param1: "Foo Parameter Value"
    },
    double: [1.23, 4.56],           
    string2: "Another String Value"  
    foo2: {                         
      param1: "Foo2 Parameter Value"
    },
    double2: [7.89],                  
    bar: "BAR2",                    
    string: "Yet another string"        
  }) {
   stringX
  }
}
`

var constructsMultipleResponse = j{
	"queryQuery1": j{
		"stringX": "queryQuery1",
	},
	"queryQuery2": j{
		"stringX": "queryQuery2",
	},
}

const contructsMultipleSameQuery = `
{
  queryQuery1(in: {
    stringX: "queryQuery1", 
    foo: {                         
      param1: "Foo Parameter Value"
    },
    double: [1.23, 4.56],           
    string2: "Another String Value"  
    foo2: {                         
      param1: "Foo2 Parameter Value"
    },
    double2: [7.89],                  
    bar: "BAR2",                    
    string: "Yet another string"      
  }) {
   stringX
  },
  queryQuery1(in: {
    stringX: "queryQuery1", 
    foo: {                         
      param1: "Foo Parameter Value"
    },
    double: [1.23, 4.56],           
    string2: "Another String Value"  
    foo2: {                         
      param1: "Foo2 Parameter Value"
    },
    double2: [7.89],                  
    bar: "BAR2",                    
    string: "Yet another string"      
  }) {
   stringX
  },
}
`

var constructsMultipleSameResponse = j{
	"queryQuery1": j{
		"stringX": "queryQuery1",
	},
}

const constructsScalarsQuery = `
mutation {
  constructsScalars_(in: {
    double: 1.1, float: 2.2, int32: 3, int64: -4, 
    uint32: 5, uint64: 6, sint32: 7, sint64: 8, 
    fixed32: 9, fixed64: 10, sfixed32: 11, sfixed64: 12,
    bool: true, stringX: "test", bytes: "dGVzdA=="
  }) {
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
    __typename
  }
}
`

var constructsScalarsResponse = j{
	"constructsScalars_": j{
		"__typename": "Scalars",
		"bool":       true,
		"bytes":      "dGVzdA==",
		"stringX":    "test",
		"double":     1.1,
		"fixed32":    float64(9),
		"fixed64":    float64(10),
		"float":      2.2,
		"int32":      float64(3),
		"int64":      float64(-4),
		"sfixed32":   float64(11),
		"sfixed64":   float64(12),
		"sint32":     float64(7),
		"sint64":     float64(8),
		"uint32":     float64(5),
		"uint64":     float64(6),
	},
}

const constructsAnyQuery = `
mutation {
  constructsAny_(in: {__typename: "Ref", localTime2: {time: "1234"}, fileEnum: BAR2, local: {
    bar1: {param1: "param1"}, en1: A1, externalTime1: {seconds:1123, nanos: 123}
  }})
}
`

var constructsAnyResponse = j{
	"constructsAny_": j{
		"__typename": "Ref",
		"en1":        "A0",
		"en2":        "A0",
		"external": j{
			"__typename": "Timestamp",
			"nanos":      float64(0),
			"seconds":    float64(0),
		},
		"file": j{
			"__typename": "Baz",
			"param1":     "",
		},
		"fileEnum": "BAR2",
		"fileMsg": j{
			"__typename": "Foo",
			"param1":     "",
			"param2":     "",
		},
		"foreign": j{
			"__typename": "Foo2",
			"param1":     "",
		},
		"gz": j{
			"__typename": "Gz",
			"param1":     "",
		},
		"local": j{
			"__typename": "Foo",
			"bar1": j{
				"__typename": "Bar",
				"param1":     "param1",
			},
			"bar2": j{
				"__typename": "Bar",
				"param1":     "",
			},
			"en1": "A1",
			"en2": "A0",
			"externalTime1": j{
				"__typename": "Timestamp",
				"nanos":      float64(123),
				"seconds":    float64(1123),
			},
			"localTime2": j{
				"__typename": "Timestamp",
				"time":       "",
			},
		},
		"localTime": j{
			"__typename": "Timestamp",
			"time":       "",
		},
		"localTime2": j{
			"__typename": "Timestamp",
			"time":       "1234",
		},
	},
}

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
	"constructsRepeated_": j{
		"bar":      l{"BAR3"},
		"bool":     l{true},
		"stringX":  l{"test"},
		"bytes":    l{"dGVzdA=="},
		"foo":      l{j{"param1": "param1", "param2": "param2"}},
		"double":   l{1.1},
		"float":    l{2.2},
		"fixed32":  l{float64(11)},
		"fixed64":  l{float64(12)},
		"int32":    l{float64(3)},
		"int64":    l{float64(4)},
		"sfixed32": l{float64(13)},
		"sfixed64": l{float64(14)},
		"sint32":   l{float64(9)},
		"sint64":   l{float64(10)},
		"uint32":   l{float64(7)},
		"uint64":   l{float64(8)},
	},
}

var constructsOneofsResponse = j{
	"constructsOneof_": j{
		"Oneof1": j{
			"__typename": "constructs_Oneof_Oneof1",
			"param3":     "3",
		},
		"Oneof2": j{
			"__typename": "constructs_Oneof_Oneof2",
			"param5":     "5",
		},
		"Oneof3": j{
			"__typename": "constructs_Oneof_Oneof3",
			"param6":     "6",
		},
		"param1": "2",
	},
}

const constructsOneofsQuery = `
mutation {
  constructsOneof_(
    in: {param1: "2", param3:"3", param5: "5", param6: "6"}
  ){ 
    param1
    Oneof1 {
      __typename
      ... on constructs_Oneof_param2 {
        param2
      }
      ... on constructs_Oneof_param3 {
        param3
      }
    }
    Oneof2 {
      __typename
      ... on constructs_Oneof_param4 {
        param4
      }
      ... on constructs_Oneof_param5 {
        param5
      }
    }
    Oneof3 {
      __typename
      ... on constructs_Oneof_param6 {
        param6
      }
    }
  }
}
`
