package bashstyle

import (
	"errors"
	"os"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var errViolationsFound = errors.New("violations found")

/* ---------------------------------------- File Closing ---------------------------------------- */

func closeFile(file *os.File, existingErr error) (err error) {
	if closeErr := file.Close(); closeErr != nil {
		return errors.Join(existingErr, closeErr)
	}

	return existingErr
}
