package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysulq/graphql-grpc-gateway/api/test"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestDecodeMaps(t *testing.T) {
	ins := &Generator{}
	msg := &test.Maps{}

	generated, err := ins.DecodeProto2GraphQL(msg, &ast.Field{
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
	})

	require.Nil(t, err)
	require.NotNil(t, generated)
}
