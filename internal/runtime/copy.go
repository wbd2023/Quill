package runtime

import (
	"io"
	"os"
)

func copyExecutable(source string, destination string) (err error) {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := input.Close()
		if err == nil {
			err = closeErr
		}
	}()

	output, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, executableMode)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := output.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(output, input); err != nil {
		return err
	}

	return output.Chmod(executableMode)
}
