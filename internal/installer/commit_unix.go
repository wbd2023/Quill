//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package installer

import (
	"os"

	"golang.org/x/sys/unix"
)

// rootedRename atomically renames oldName to newName within the directory anchored by
// directoryRoot.
//
// Go 1.24 os.Root does not yet expose Rename, so the rename is issued relative to the anchored
// directory handle via renameat(2). Because both names resolve against the pinned directory
// descriptor, a parent symlink swap performed after the directory was anchored cannot redirect the
// rename outside the repository.
func rootedRename(
	directoryRoot *os.Root,
	_ string,
	oldName string,
	newName string,
) (err error) {
	dir, err := directoryRoot.Open(".")
	if err != nil {
		return err
	}
	defer func() {
		_ = dir.Close() // closing a read-only directory cannot change renameat's completed result
	}()

	return unix.Renameat(int(dir.Fd()), oldName, int(dir.Fd()), newName)
}
