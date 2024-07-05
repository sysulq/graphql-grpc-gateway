package test

import (
	"bytes"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

func compareGraphql(t *testing.T, got, expect *ast.Schema) {
	t.Helper()

	expectedGraphql := &bytes.Buffer{}
	actualGraphql := &bytes.Buffer{}
	formatter.NewFormatter(actualGraphql).FormatSchema(got)
	formatter.NewFormatter(expectedGraphql).FormatSchema(expect)

	if actualGraphql.String() != expectedGraphql.String() {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expectedGraphql.String()),
			B:        difflib.SplitLines(actualGraphql.String()),
			FromFile: "expect",
			ToFile:   "got",
			Context:  3,
		}
		t.Errorf("Generated graphql file does not match expectations")
		str, _ := difflib.GetUnifiedDiffString(diff)
		t.Errorf("%s", str)
	}
}
