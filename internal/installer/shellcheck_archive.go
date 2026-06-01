package installer

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

/* ------------------------------------------- Archive ------------------------------------------ */

func extractShellcheckBinary(
	archivePath string,
	destination string,
	version string,
) (binaryPath string, err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close shellcheck archive %q: %w", archivePath, closeErr)
		}
	}()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return "", err
	}

	expectedName := path.Join(shellcheckArchiveRoot(version), "shellcheck")
	targetPath := filepath.Join(destination, filepath.FromSlash(expectedName))
	foundBinary := false

	tarReader := tar.NewReader(xzReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if !foundBinary {
				return "", fmt.Errorf("shellcheck archive missing %s", expectedName)
			}

			return targetPath, nil
		}

		if err != nil {
			return "", err
		}

		cleanName, err := validateShellcheckArchiveEntry(header, version)
		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue

		case tar.TypeReg:
			if cleanName != expectedName {
				continue
			}

			foundBinary = true
			if err = writeShellcheckBinary(targetPath, tarReader); err != nil {
				return "", err
			}

		default:
			return "", fmt.Errorf("unsupported shellcheck archive entry %q", header.Name)
		}
	}
}

func writeShellcheckBinary(targetPath string, source io.Reader) (err error) {
	if err = os.MkdirAll(filepath.Dir(targetPath), defaultDirectoryMode); err != nil {
		return err
	}

	targetFile, err := os.OpenFile(
		targetPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		executableMode,
	)
	if err != nil {
		return err
	}

	if _, err = io.Copy(targetFile, source); err != nil {
		if closeErr := targetFile.Close(); closeErr != nil {
			return fmt.Errorf(
				"copy shellcheck file %q: %w",
				targetPath,
				errors.Join(err, closeErr),
			)
		}
		return err
	}

	return targetFile.Close()
}
