//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris) && !windows

package installer

import (
	"errors"
	"os"
)

// errRootedCommitUnavailable reports that this platform cannot perform a handle-relative archive
// commit. These targets expose neither a directory-handle-relative rename (Windows) nor renameat(2)
// (the openat Unix family), so Quill declines rather than falling back to a racy absolute rename.
var errRootedCommitUnavailable = errors.New(
	"rooted archive installation commit is not supported on this platform",
)

// rootedRename is unsupported on these platforms; returning an error keeps the write from silently
// falling back to a name-based, race-vulnerable rename. See errRootedCommitUnavailable.
func rootedRename(_ *os.Root, _ string, _ string, _ string) error {
	return errRootedCommitUnavailable
}
