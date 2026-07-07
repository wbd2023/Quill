package installer

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
	archive string,
	dir string,
	version string,
) (extracted string, err error) {
	file, err := os.Open(archive)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", archive, closeErr)
		}
	}()

	reader, err := xz.NewReader(file)
	if err != nil {
		return "", err
	}

	expected := path.Join(shellcheckArchiveRoot(version), "shellcheck")
	target := filepath.Join(dir, filepath.FromSlash(expected))
	found := false

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if !found {
				return "", fmt.Errorf("shellcheck archive missing %s", expected)
			}

			return target, nil
		}

		if err != nil {
			return "", err
		}

		name, err := validateShellcheckEntry(header, version)
		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue

		case tar.TypeReg:
			if name != expected {
				continue
			}

			found = true
			if err = writeExecutable(target, tarReader); err != nil {
				return "", err
			}

		default:
			return "", fmt.Errorf("unsupported shellcheck archive entry %q", header.Name)
		}
	}
}

/* ------------------------------------- Archive Validation ------------------------------------- */

// validateShellcheckEntry checks that a tar header is safe (no symlinks, no path traversal) and
// matches the expected shellcheck archive layout for the given version.
func validateShellcheckEntry(header *tar.Header, version string) (name string, err error) {
	if header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink {
		return "", fmt.Errorf("shellcheck archive contains link entry %q", header.Name)
	}

	raw := header.Name
	if header.Typeflag == tar.TypeDir {
		raw = strings.TrimSuffix(raw, "/")
	}

	name = path.Clean(raw)
	if name == "." ||
		name != raw ||
		path.IsAbs(raw) ||
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
