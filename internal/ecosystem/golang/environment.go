package golang

import (
	"path/filepath"

	"ciphera/tools/internal/runtime"
)

// BuildCacheDirectory returns the Go build cache directory.
func BuildCacheDirectory(layout runtime.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-build")
}

// ModuleCacheDirectory returns the Go module cache directory.
func ModuleCacheDirectory(layout runtime.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-mod")
}

// GoPath returns the GOPATH directory.
func GoPath(layout runtime.Layout) (path string) {
	return filepath.Join(layout.CacheDirectory(), "gopath")
}

// Environment returns the environment variables for executing Go tooling with isolated caches.
// GOCACHE, GOMODCACHE, and GOPATH are set to layout-derived paths; PATH is the path argument.
func Environment(layout runtime.Layout, path string) (environment map[string]string) {
	return map[string]string{
		"PATH":       path,
		"GOCACHE":    BuildCacheDirectory(layout),
		"GOMODCACHE": ModuleCacheDirectory(layout),
		"GOPATH":     GoPath(layout),
	}
}
