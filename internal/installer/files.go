package installer

import (
	"io"
	"os"
)

func copyExecutable(source string, destination string) (err error) {
	return copyFile(source, destination, executableMode)
}

func copyFile(source string, destination string, mode os.FileMode) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := sourceFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	destinationFile, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := destinationFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return destinationFile.Chmod(mode)
}
