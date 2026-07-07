package installer

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ----------------------------------------- Resolution ----------------------------------------- */

// Resolve downloads every platform archive for each archive-installed tool declared in the
// profile, hashes each one, and returns the entries that make up quill.lock. Failures from
// independent tools or platforms are collected and returned as a joined error; the partial
// results are still returned so the caller can write what resolved successfully if it chooses.
func Resolve(
	writer io.Writer,
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
) (entries []lockfile.Archive, err error) {
	var errs []error
	for _, tool := range tools {
		capability, found := capabilities[tool.ID]
		if !found {
			errs = append(errs, fmt.Errorf("missing tool capability %q", tool.ID))
			continue
		}

		if capability.InstallKind != toolchain.InstallKindArchive || capability.Archive == nil {
			continue
		}

		entry, resolveErr := resolveArchive(writer, tool, capability)
		if resolveErr != nil {
			errs = append(errs, resolveErr)
			continue
		}

		entries = append(entries, entry)
	}

	err = errors.Join(errs...)
	return entries, err
}

func resolveArchive(
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (archive lockfile.Archive, err error) {
	spec := *capability.Archive
	hashes := make(map[string]string, len(spec.Platforms))

	for platformKey := range spec.Platforms {
		hash, hashErr := resolvePlatform(writer, spec, tool, platformKey)
		if hashErr != nil {
			return lockfile.Archive{}, fmt.Errorf(
				"resolve %s %s: %w",
				tool.ID,
				platformKey,
				hashErr,
			)
		}

		hashes[platformKey] = hash
	}

	return lockfile.Archive{
		Tool:    tool.ID,
		Version: tool.PinnedVersion,
		Hashes:  hashes,
	}, nil
}

func resolvePlatform(
	writer io.Writer,
	spec toolchain.ArchiveSpec,
	tool style.Tool,
	platformKey string,
) (hash string, err error) {
	platform := spec.Platforms[platformKey]
	url := spec.URL(tool.PinnedVersion, platform)

	dir, err := os.MkdirTemp("", "quill-resolve-*")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	archive := filepath.Join(dir, "archive")
	if _, err = fmt.Fprintf(writer, "Resolving %s %s...\n", tool.Name, platformKey); err != nil {
		return "", err
	}

	if err = downloadFile(url, archive); err != nil {
		return "", err
	}

	return sha256File(archive)
}

func sha256File(path string) (hash string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %q: %w", path, closeErr)
		}
	}()

	digest := sha256.New()
	if _, err = io.Copy(digest, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}
