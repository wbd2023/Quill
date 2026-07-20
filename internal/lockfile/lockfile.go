// Package lockfile reads and writes quill.lock, the resolved-state file for
// archive-installed tools. The Profile (quill.toml) declares intent (which
// version); the lockfile records what was verified (the per-platform hashes).
package lockfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// DefaultFilename is the lockfile filename loaded from repository roots.
const DefaultFilename = "quill.lock"

// Lockfile is the parsed content of quill.lock.
type Lockfile struct {
	// Loaded reports whether a lockfile was present on disk. False means the
	// file was absent; the caller should direct the user to run 'quill lock'.
	Loaded bool
	// Archives is the set of recorded archive-tool entries, keyed by tool ID.
	Archives map[string]Archive
}

// Archive is one tool's recorded hashes for a specific version.
type Archive struct {
	Tool    string
	Version string
	Hashes  map[string]string
}

// Load reads the lockfile from a repository root. A missing lockfile is not an
// error; the returned Lockfile has Loaded=false so the caller can distinguish
// "no lockfile" from "lockfile missing an entry".
func Load(root string) (lockfile Lockfile, err error) {
	path := filepath.Join(root, DefaultFilename)
	contents, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Lockfile{Loaded: false}, nil
		}

		return Lockfile{}, fmt.Errorf("read lockfile %q: %w", path, err)
	}

	lockfile, err = Decode(string(contents))
	if err != nil {
		return Lockfile{}, fmt.Errorf("load lockfile %q: %w", path, err)
	}

	lockfile.Loaded = true
	return lockfile, nil
}

// HashFor looks up the recorded SHA-256 hash for a tool, version, and platform.
// The distinct error cases give the caller actionable messages.
func (l Lockfile) HashFor(
	toolID string,
	wantVersion string,
	goos string,
	goarch string,
) (hash string, err error) {
	if !l.Loaded {
		return "", fmt.Errorf("quill.lock not found; run 'quill lock' to populate")
	}

	archive, ok := l.Archives[toolID]
	if !ok {
		return "", fmt.Errorf("no lockfile entry for %s; run 'quill lock'", toolID)
	}

	if archive.Version != wantVersion {
		return "", fmt.Errorf(
			"lockfile has %s %s but profile pins %s; run 'quill lock'",
			toolID,
			archive.Version,
			wantVersion,
		)
	}

	hash, ok = archive.Hashes[goos+"/"+goarch]
	if !ok {
		return "", fmt.Errorf(
			"no lockfile hash for %s on %s/%s; run 'quill lock'",
			toolID,
			goos,
			goarch,
		)
	}

	return hash, nil
}
