//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

// This file imports tools that are used in the development process.
// The build tag ensures they're not included in regular builds.
