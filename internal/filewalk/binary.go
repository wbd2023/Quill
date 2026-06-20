package filewalk

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const binaryProbeLimit = 4096

// IsBinaryFile is binary file.
func IsBinaryFile(path string) (binary bool) {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	buffer := make([]byte, binaryProbeLimit)
	count, readErr := file.Read(buffer)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		if closeErr := file.Close(); closeErr != nil {
			return false
		}
		return false
	}

	if closeErr := file.Close(); closeErr != nil {
		return false
	}

	return bytes.IndexByte(buffer[:count], 0) >= 0
}
