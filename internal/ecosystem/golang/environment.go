package golang

import (
	"path/filepath"

	"ciphera/tools/internal/runtime"
)

// BuildCache returns the Go build cache directory for the given layout.
func BuildCache(layout runtime.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-build")
}

// ModuleCache returns the Go module cache directory for the given layout.
func ModuleCache(layout runtime.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-mod")
}

// GoPath returns the GOPATH directory for the given layout.
func GoPath(layout runtime.Layout) (path string) {
	return filepath.Join(layout.CacheDirectory(), "gopath")
}

// Environment builds the Go environment variables for running Go tooling within the engine's
// isolated layout. searchPath is the PATH value that makes installed tool binaries discoverable.
// Callers may add consumer-specific variables (eg GOLANGCI_LINT_CACHE for the checker, GOBIN for
// the installer) to the result.
func Environment(layout runtime.Layout, searchPath string) (environment map[string]string) {
	return map[string]string{
		"PATH":       searchPath,
		"GOCACHE":    BuildCache(layout),
		"GOMODCACHE": ModuleCache(layout),
		"GOPATH":     GoPath(layout),
	}
}
