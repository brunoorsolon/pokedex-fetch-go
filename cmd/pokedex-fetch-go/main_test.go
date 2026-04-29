package main

import "testing"

func TestMainPackageSmoke(t *testing.T) {
	// Keep this package visible to `go test ./...`. The real behavior is covered
	// through internal/cli tests; calling main would execute the CLI and exit.
}
