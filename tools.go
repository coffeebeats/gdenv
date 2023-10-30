//go:build tools

// Package tools lists all of the Go-related tools required to develop the
// project. To install them, run the following command:
//
// cat tools.go | grep _ | grep -v '//' | awk -F'"' '{print $2}' | xargs -tI % go install %
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "golang.org/x/tools/cmd/goimports"
)
