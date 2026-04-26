package runtime

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

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
			if err = os.MkdirAll(filepath.Dir(targetPath), defaultDirectoryMode); err != nil {
				return "", err
			}
			targetFile, err := os.OpenFile(
				targetPath,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				executableMode,
			)
			if err != nil {
				return "", err
			}

			if _, err = io.Copy(targetFile, tarReader); err != nil {
				if closeErr := targetFile.Close(); closeErr != nil {
					return "", fmt.Errorf(
						"copy shellcheck file %q: %w",
						targetPath,
						errors.Join(err, closeErr),
					)
				}
				return "", err
			}

			if err = targetFile.Close(); err != nil {
				return "", err
			}

		default:
			return "", fmt.Errorf("unsupported shellcheck archive entry %q", header.Name)
		}
	}
}

func validateShellcheckArchiveEntry(
	header *tar.Header,
	version string,
) (name string, err error) {
	switch header.Typeflag {
	case tar.TypeSymlink, tar.TypeLink:
		return "", fmt.Errorf("shellcheck archive contains link entry %q", header.Name)
	}

	rawName := header.Name
	if header.Typeflag == tar.TypeDir {
		rawName = strings.TrimSuffix(rawName, "/")
	}

	name = path.Clean(rawName)
	if name == "." ||
		name != rawName ||
		path.IsAbs(rawName) ||
		strings.HasPrefix(name, "../") ||
		strings.Contains(name, "/../") {
		return "", fmt.Errorf("unsafe shellcheck archive path %q", header.Name)
	}

	root := shellcheckArchiveRoot(version)
	switch name {
	case root,
		path.Join(root, "LICENSE.txt"),
		path.Join(root, "README.txt"),
		path.Join(root, "shellcheck"):
		return name, nil
	default:
		return "", fmt.Errorf("unexpected shellcheck archive entry %q", header.Name)
	}
}

func shellcheckArchiveRoot(version string) (root string) {
	return "shellcheck-v" + version
}
