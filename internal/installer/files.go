package installer

import (
	"fmt"
	"io"
	"os"
)

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
// a missing or regular destination with the executable file. Symlinked parents and non-regular
// destinations are rejected.
func writeExecutable(root string, destination string, reader io.Reader) (err error) {
	destination, directory, _, err := prepareExecutableDestination(root, destination)
	if err != nil {
		return err
	}

	file, err := os.CreateTemp(directory, ".quill-install-*")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer func() {
		_ = os.Remove(temporary)
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

	info, statErr := os.Lstat(destination)
	if statErr == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refuse to replace non-regular destination %q", destination)
	}

	if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("inspect destination %q: %w", destination, statErr)
	}

	if err = os.Rename(temporary, destination); err != nil {
		return fmt.Errorf("replace destination %q: %w", destination, err)
	}

	return nil
}
