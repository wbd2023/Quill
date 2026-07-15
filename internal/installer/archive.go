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

	"ciphera/tools/internal/toolchain"

	"github.com/ulikunitz/xz"
)

/* ----------------------------------------- Extraction ----------------------------------------- */

func extractBinary(
	archive string,
	dir string,
	install toolchain.GitHubInstall,
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

	tag := fmt.Sprintf(install.Tag, version)
	expected := fmt.Sprintf(install.Path, tag)
	target := filepath.Join(dir, filepath.FromSlash(expected))
	found := false

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if !found {
				return "", fmt.Errorf("archive missing %s", expected)
			}

			return target, nil
		}

		if err != nil {
			return "", err
		}

		name, err := validateArchiveEntry(header, expected)
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
			return "", fmt.Errorf("unsupported archive entry %q", header.Name)
		}
	}
}

/* ----------------------------------------- Validation ----------------------------------------- */

// validateArchiveEntry checks that a tar header is safe (no symlinks, no path traversal) and
// either is or sits under the expected binary's directory within the archive.
func validateArchiveEntry(header *tar.Header, expected string) (name string, err error) {
	if header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink {
		return "", fmt.Errorf("archive contains link entry %q", header.Name)
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
		return "", fmt.Errorf("unsafe archive path %q", header.Name)
	}

	root := path.Dir(expected)
	if name != root && !strings.HasPrefix(name, root+"/") {
		return "", fmt.Errorf("unexpected archive entry %q", header.Name)
	}

	return name, nil
}
