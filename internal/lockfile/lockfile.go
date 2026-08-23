// Package lockfile reads and writes quill.lock, the resolved-state file for
// archive-installed tools. The Profile (quill.toml) declares intent (which
// version); the lockfile records what was verified (the per-platform hashes).
package lockfile

import (
	"errors"
	"fmt"
	"io"
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

// maxLockfileBytes caps how much of quill.lock Load reads into memory. A lockfile is a small
// TOML manifest of resolved archive hashes (a few KiB at most); the cap bounds memory when a
// hostile or corrupt root presents an outsized regular file and pairs with the rooted open and
// regular-file check to make Load refuse unbounded input.
const maxLockfileBytes int64 = 1 << 20 // 1 MiB

// Load reads the lockfile from a repository root. The root is opened as a Go 1.24 rooted
// filesystem so a symlinked or escaping quill.lock cannot redirect the read outside the
// repository. Only a regular, non-symlink file at most maxLockfileBytes in size is decoded;
// symlinks, FIFOs, sockets, devices, and oversized entries fail closed. A missing lockfile is
// not an error; the returned Lockfile has Loaded=false so the caller can distinguish "no
// lockfile" from "lockfile missing an entry".
func Load(root string) (lockfile Lockfile, err error) {
	path := filepath.Join(root, DefaultFilename)

	rooted, err := os.OpenRoot(root)
	if err != nil {
		// A missing root means there is no lockfile to load; preserve the
		// historical "not loaded, not an error" outcome for callers that probe
		// before locking.
		if errors.Is(err, os.ErrNotExist) {
			return Lockfile{Loaded: false}, nil
		}

		return Lockfile{}, fmt.Errorf("open repository root %q: %w", root, err)
	}
	defer func() { _ = rooted.Close() }()

	info, err := rooted.Lstat(DefaultFilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Lockfile{Loaded: false}, nil
		}

		return Lockfile{}, fmt.Errorf("stat lockfile %q: %w", path, err)
	}

	if !info.Mode().IsRegular() {
		return Lockfile{}, fmt.Errorf(
			"load lockfile %q: not a regular file (%s)",
			path,
			info.Mode().Type(),
		)
	}

	if info.Size() > maxLockfileBytes {
		return Lockfile{}, fmt.Errorf(
			"load lockfile %q: size %d exceeds %d-byte limit",
			path,
			info.Size(),
			maxLockfileBytes,
		)
	}

	contents, err := readBounded(rooted, DefaultFilename)
	if err != nil {
		return Lockfile{}, fmt.Errorf("read lockfile %q: %w", path, err)
	}

	lockfile, err = Decode(string(contents))
	if err != nil {
		return Lockfile{}, fmt.Errorf("load lockfile %q: %w", path, err)
	}

	lockfile.Loaded = true
	return lockfile, nil
}

// readBounded reads name relative to root, capping the result at maxLockfileBytes. A file that
// grew past the stat-checked size between the Lstat and the open (or any other unbounded source)
// surfaces as an error rather than unbounded memory use.
func readBounded(root *os.Root, name string) (contents []byte, err error) {
	file, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	contents, err = io.ReadAll(io.LimitReader(file, maxLockfileBytes+1))
	if err != nil {
		return nil, err
	}

	if int64(len(contents)) > maxLockfileBytes {
		return nil, fmt.Errorf("size exceeds %d-byte limit", maxLockfileBytes)
	}

	return contents, nil
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
