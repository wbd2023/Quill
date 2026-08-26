package installer

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const temporaryFilePermissions os.FileMode = 0o600

/* ---------------------------------- Staged Executable Writes ---------------------------------- */

// copyExecutable copies source to destination as an executable file.
func copyExecutable(root string, source string, destination string) (err error) {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := src.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", source, closeErr)
		}
	}()

	return writeExecutable(root, destination, src)
}

// writeExecutable stages reader beneath root in the destination directory, then atomically replaces
// a missing or regular destination with the executable file.
//
// Staging, leaf validation, and the atomic rename all operate relative to a directory handle rooted
// at root, so a parent symlink swap performed between validation and commit cannot redirect the
// write outside the repository. Symlinked parents and non-regular destinations are still rejected
// up front by prepareExecutableDestination.
func writeExecutable(root string, destination string, reader io.Reader) (err error) {
	destination, directory, _, err := prepareExecutableDestination(root, destination)
	if err != nil {
		return err
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve installation root %q: %w", root, err)
	}

	relative, err := filepath.Rel(rootAbs, destination)
	if err != nil {
		return fmt.Errorf("resolve destination beneath root: %w", err)
	}

	rootDir, err := os.OpenRoot(rootAbs)
	if err != nil {
		return fmt.Errorf("open installation root %q: %w", rootAbs, err)
	}
	defer func() {
		if closeErr := rootDir.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close installation root %q: %w", rootAbs, closeErr)
		}
	}()

	// Anchor the destination's parent directory relative to the rooted handle. The handle pins the
	// directory inode, so a later swap of the parent path cannot move the staging or commit target.
	directoryRoot, err := rootDir.OpenRoot(filepath.Dir(relative))
	if err != nil {
		return fmt.Errorf("anchor installation directory %q: %w", directory, err)
	}
	defer func() {
		if closeErr := directoryRoot.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close installation directory %q: %w", directory, closeErr)
		}
	}()

	leaf := filepath.Base(relative)
	file, temporary, err := createTempInRoot(directoryRoot, ".quill-install-*")
	if err != nil {
		return fmt.Errorf("stage executable in %q: %w", directory, err)
	}
	defer func() {
		_ = directoryRoot.Remove(temporary)
	}()
	defer func() {
		if file == nil {
			return
		}
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", temporary, closeErr)
		}
	}()

	if err = file.Chmod(standardPermissions); err != nil {
		return err
	}

	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("close %q: %w", temporary, err)
	}
	file = nil

	info, statErr := directoryRoot.Lstat(leaf)
	if statErr == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refuse to replace non-regular destination %q", destination)
	}

	if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("inspect destination %q: %w", destination, statErr)
	}

	if err = rootedRename(directoryRoot, directory, temporary, leaf); err != nil {
		return fmt.Errorf("replace destination %q: %w", destination, err)
	}

	return nil
}

/* ----------------------------------- Rooted Temporary Files ----------------------------------- */

// createTempInRoot creates a unique temporary regular file inside dir, mirroring os.CreateTemp
// against a rooted directory handle (Go 1.24 os.Root does not yet provide CreateTemp).
func createTempInRoot(
	dir *os.Root,
	pattern string,
) (file *os.File, name string, err error) {
	prefix, suffix := pattern, ""
	if star := strings.LastIndexByte(pattern, '*'); star >= 0 {
		prefix, suffix = pattern[:star], pattern[star+1:]
	}

	for attempt := 0; ; attempt++ {
		name = prefix + strconv.FormatUint(uint64(rand.Uint32()), 10) + suffix
		file, err = dir.OpenFile(
			name,
			os.O_RDWR|os.O_CREATE|os.O_EXCL,
			temporaryFilePermissions,
		)
		if err == nil {
			return file, name, nil
		}

		if !os.IsExist(err) || attempt >= 10000 {
			return nil, "", err
		}
	}
}
