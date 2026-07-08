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

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	limit   = 128 << 20
	timeout = 30 * time.Second
)

/* ------------------------------------------ Download ------------------------------------------ */

// downloadFile fetches a URL and writes it to destination via an atomic temp-file rename. The
// download is capped at limit bytes to prevent unbounded memory or disk usage.
func downloadFile(url string, destination string) (err error) {
	response, err := (&http.Client{Timeout: timeout}).Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := response.Body.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close download response body: %w", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: unexpected HTTP status %s", url, response.Status)
	}

	if err = os.MkdirAll(filepath.Dir(destination), standardPermissions); err != nil {
		return err
	}

	file, err := os.CreateTemp(filepath.Dir(destination), ".download-*")
	if err != nil {
		return err
	}
	temp := file.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(temp)
		}
	}()

	reader := &io.LimitedReader{R: response.Body, N: limit + 1}
	written, err := io.Copy(file, reader)
	if err != nil {
		_ = file.Close()
		return err
	}

	if written > limit {
		_ = file.Close()
		return fmt.Errorf("download exceeds maximum size")
	}

	if err = file.Close(); err != nil {
		return err
	}

	return os.Rename(temp, destination)
}

/* ------------------------------------------ Checksums ----------------------------------------- */

// hashFile returns the SHA-256 hex digest of the file at path.
func hashFile(path string) (digest string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", path, closeErr)
		}
	}()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// verifyChecksum reports whether the SHA-256 hash of the file at path matches the expected hex
// digest.
func verifyChecksum(path string, expected string) (err error) {
	actual, err := hashFile(path)
	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("checksum mismatch for %s", path)
	}

	return nil
}
