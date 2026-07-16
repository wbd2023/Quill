package node

import (
	"path/filepath"

	"ciphera/tools/internal/runtime"
)

// Directory returns the directory where Node tooling state lives.
func Directory(layout runtime.Layout) (directory string) {
	return filepath.Join(layout.StateDirectory(), "npm")
}

// BinaryDirectory returns the directory where Node-installed binaries live.
func BinaryDirectory(layout runtime.Layout) (directory string) {
	return filepath.Join(Directory(layout), "node_modules", ".bin")
}

// Cache returns the NPM cache directory.
func Cache(layout runtime.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "npm")
}

// Environment builds the NPM environment variables for installing Node tooling within the engine's
// isolated layout. path is the PATH value that makes installed tool binaries discoverable.
func Environment(layout runtime.Layout, path string) (environment map[string]string) {
	return map[string]string{
		"PATH":             path,
		"npm_config_cache": Cache(layout),
	}
}
