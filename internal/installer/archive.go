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

	"github.com/wbd2023/Quill/internal/toolchain"

	"github.com/ulikunitz/xz"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	maxArchiveEntries = 4096
	maxArchiveSize    = 128 << 20
)

/* ----------------------------------------- Extraction ----------------------------------------- */

func extractBinary(
	archive string,
	dir string,
	install toolchain.GitHubInstall,
	version string,
) (extracted string, err error) {
	return extractBinaryUpTo(archive, dir, install, version, maxArchiveSize)
}

func extractBinaryUpTo(
	archive string,
	dir string,
	install toolchain.GitHubInstall,
	version string,
	maxSize int64,
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

	var entryCount int
	var uncompressedSize int64
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
		entryCount++
		if entryCount > maxArchiveEntries {
			return "", fmt.Errorf("archive exceeds maximum entry count")
		}

		name, err := validateArchiveEntry(header, expected)
		if err != nil {
			return "", err
		}

		if header.Size < 0 || header.Size > maxSize-uncompressedSize {
			return "", fmt.Errorf(
				"archive uncompressed size exceeds maximum of %d bytes",
				maxSize,
			)
		}
		uncompressedSize += header.Size

		switch header.Typeflag {

		case tar.TypeDir:
			continue

		case tar.TypeReg:
			if name != expected {
				continue
			}

			found = true
			if err = writeExecutable(dir, target, tarReader); err != nil {
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
