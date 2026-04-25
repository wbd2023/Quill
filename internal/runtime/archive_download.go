package runtime

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

/* -------------------------------------------- Types ------------------------------------------- */

const downloadTimeout = 30 * time.Second

type httpStatusError struct {
	url        string
	status     string
	statusCode int
}

func (err httpStatusError) Error() (message string) {
	return fmt.Sprintf("%s: unexpected HTTP status %s", err.url, err.status)
}

/* ------------------------------------------ Download ------------------------------------------ */

func downloadBytes(url string) (body []byte, err error) {
	client := &http.Client{Timeout: downloadTimeout}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := response.Body.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close download response body: %w", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, httpStatusError{
			url:        url,
			status:     response.Status,
			statusCode: response.StatusCode,
		}
	}

	return io.ReadAll(response.Body)
}

func downloadFile(url string, destination string) (err error) {
	body, err := downloadBytes(url)
	if err != nil {
		return err
	}

	return os.WriteFile(destination, body, defaultDownloadMode)
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
