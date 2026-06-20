package installer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/* -------------------------------------------- Types ------------------------------------------- */

// download constants.
const (
	downloadTimeout = 30 * time.Second
	maxDownloadSize = 128 << 20
)

type httpStatusError struct {
	url        string
	status     string
	statusCode int
}

func (err httpStatusError) Error() (message string) {
	return fmt.Sprintf("%s: unexpected HTTP status %s", err.url, err.status)
}

/* ------------------------------------------ Download ------------------------------------------ */

func downloadFile(url string, destination string) (err error) {
	client := &http.Client{Timeout: downloadTimeout}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := response.Body.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close download response body: %w", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return httpStatusError{
			url:        url,
			status:     response.Status,
			statusCode: response.StatusCode,
		}
	}

	if err = os.MkdirAll(filepath.Dir(destination), defaultDirectoryMode); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(destination), ".download-*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(tempPath)
		}
	}()

	if err = writeDownloadFile(tempFile, response.Body, maxDownloadSize); err != nil {
		_ = tempFile.Close() // returning the write error is more useful
		return err
	}

	if err = tempFile.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, destination)
}

func writeDownload(destination string, body io.Reader) (err error) {
	return writeDownloadWithLimit(destination, body, maxDownloadSize)
}

func writeDownloadWithLimit(destination string, body io.Reader, limit int64) (err error) {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, downloadMode)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close download file %q: %w", destination, closeErr)
		}
	}()

	return writeDownloadFile(file, body, limit)
}

func writeDownloadFile(file *os.File, body io.Reader, limit int64) (err error) {
	reader := &io.LimitedReader{R: body, N: limit + 1}
	written, err := io.Copy(file, reader)
	if err != nil {
		return err
	}

	if written > limit {
		return fmt.Errorf("download exceeds maximum size")
	}

	return file.Chmod(downloadMode)
}

/* ------------------------------------------ Checksums ----------------------------------------- */

func verifyFileChecksum(
	archivePath string,
	archiveName string,
	expected string,
) (err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close archive %q: %w", archivePath, closeErr)
		}
	}()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return err
	}

	actual := fmt.Sprintf("%x", hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf("checksum mismatch for %s", archiveName)
	}

	return nil
}
