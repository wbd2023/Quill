package node

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wbd2023/quill/internal/workspace"
)

// InstallDirectory returns the directory where npm install operates.
func InstallDirectory(layout workspace.Layout) (directory string) {
	return filepath.Join(layout.StateDirectory, "npm")
}

// BinaryDirectory returns the directory where Node-installed binaries live.
func BinaryDirectory(layout workspace.Layout) (directory string) {
	return filepath.Join(InstallDirectory(layout), "node_modules", ".bin")
}

// BinaryPath returns the platform-specific path of a managed npm tool.
func BinaryPath(layout workspace.Layout, command string) (path string) {
	if runtime.GOOS == "windows" && !strings.EqualFold(filepath.Ext(command), ".cmd") {
		command += ".cmd"
	}

	return filepath.Join(BinaryDirectory(layout), command)
}

// CacheDirectory returns the npm cache directory.
func CacheDirectory(layout workspace.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "npm")
}

// Environment returns the environment variables for executing Node tooling with isolated caches.
// npm_config_cache is set to a layout-derived path; PATH is the path argument.
func Environment(layout workspace.Layout, path string) (environment map[string]string) {
	return map[string]string{
		"PATH":             path,
		"npm_config_cache": CacheDirectory(layout),
	}
}
