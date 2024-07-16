//go:build tools
// +build tools

package main

import (
	_ "github.com/go-kod/kod/cmd/kod"
	_ "go.uber.org/mock/mockgen"
)
