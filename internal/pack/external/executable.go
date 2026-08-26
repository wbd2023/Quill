package external

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ResolveExecutable resolves a manifest runtime command against its Pack directory and verifies the
// result is an executable file contained beneath the Pack directory. The command must be a
// repository-relative path with no parent-traversal escape: a Pack may not reach outside its own
// directory through "..", absolute paths, drive roots, or symlinks.
//
// dir must already be canonical (symlink-resolved). The resolved executable is
// returned as an absolute path suitable for direct subprocess launch; the launch boundary
// does not re-resolve it against PATH.
func ResolveExecutable(dir string, command string) (resolved string, err error) {
	if strings.ContainsRune(command, '\x00') {
		return "", fmt.Errorf("pack runtime command %q contains a NUL byte", command)
	}

	normalised := strings.ReplaceAll(command, "\\", "/")
	if normalised == "" || strings.HasPrefix(normalised, "/") || isWindowsDrive(command) {
		return "", fmt.Errorf("pack runtime command %q must be a Pack-relative path", command)
	}

	joined := filepath.Join(dir, filepath.Clean(normalised))
	if !isWithinRoot(dir, joined) {
		return "", fmt.Errorf("pack runtime command %q escapes the Pack directory", command)
	}

	// Resolve every link in the joined path and re-verify the physical target stays inside the
	// Pack directory. The lexical check alone admits a symlink inside the Pack that points outside
	// it, so the resolved target is checked after EvalSymlinks rather than trusting the link path.
	resolved, err = filepath.EvalSymlinks(joined)
	if err != nil {
		return "", fmt.Errorf("pack runtime executable %q: %w", command, err)
	}

	if !isWithinRoot(dir, resolved) {
		return "", fmt.Errorf(
			"pack runtime command %q resolves outside the Pack directory",
			command,
		)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("pack runtime executable %q: %w", command, err)
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf(
			"pack runtime command %q must resolve to an executable regular file",
			command,
		)
	}

	if runtime.GOOS != "windows" && info.Mode()&0o111 == 0 {
		return "", fmt.Errorf("pack runtime command %q is not executable", command)
	}

	return resolved, nil
}

func isWindowsDrive(value string) (drive bool) {
	if len(value) < 2 || value[1] != ':' {
		return false
	}
	head := value[0]
	return head >= 'A' && head <= 'Z' || head >= 'a' && head <= 'z'
}
