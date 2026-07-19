//go:build tools
// +build tools

// Package tools tracks build-time dependencies for `go run` invocations,
// notably the golang-migrate CLI which needs explicit database driver imports
// to register the "postgres" and "file" drivers at build time.
package tools

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)