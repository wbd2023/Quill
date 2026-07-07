package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyExecutable copies source to destination as an executable file.
func copyExecutable(source string, destination string) (err error) {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := src.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", source, closeErr)
		}
	}()

	return writeExecutable(destination, src)
}

// writeExecutable writes reader to destination as an executable file, creating parent
// directories as needed.
func writeExecutable(destination string, reader io.Reader) (err error) {
	if err = os.MkdirAll(filepath.Dir(destination), standardPermissions); err != nil {
		return err
	}

	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, standardPermissions)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", destination, closeErr)
		}
	}()

	_, err = io.Copy(file, reader)
	return err
}
