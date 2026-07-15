package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// anyExecutePerm is the bitmask for any execute permission bit (user, group, or other).
const anyExecutePerm os.FileMode = 0o111

// ResolveCommandPath resolves command to an executable path. It honours the provided environment's
// PATH when set, otherwise falls back to exec.LookPath. Absolute paths and paths containing a
// separator are returned as-is.
func ResolveCommandPath(command string, environment map[string]string) (path string, err error) {
	paths, ok := environment["PATH"]
	if !ok {
		paths = os.Getenv("PATH")
	}
	if paths == "" {
		return exec.LookPath(command)
	}

	if filepath.IsAbs(command) || strings.ContainsRune(command, os.PathSeparator) {
		return command, nil
	}

	for _, dir := range filepath.SplitList(paths) {
		candidate := filepath.Join(dir, command)
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() || !isExecutable(info.Mode()) {
			continue
		}

		return candidate, nil
	}

	return "", exec.ErrNotFound
}

func isExecutable(mode os.FileMode) (executable bool) {
	return mode&anyExecutePerm != 0
}
