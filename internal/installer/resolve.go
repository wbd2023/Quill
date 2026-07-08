package installer

import (
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
// independent tools or platforms are collected and returned as a joined error.
func Resolve(
	writer io.Writer,
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
) (entries []lockfile.Archive, err error) {
	return resolveWith(writer, tools, capabilities, resolveArchive)
}

// archiveResolver assembles one tool's lockfile entry by resolving each platform. Extracted as a
// parameter so the iteration and filtering in resolveWith is testable without network I/O.
type archiveResolver func(
	writer io.Writer,
	tool style.Tool,
	install toolchain.ArchiveInstall,
	resolveOne platformResolver,
) (archive lockfile.Archive, err error)

func resolveWith(
	writer io.Writer,
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
	resolveArchiveStep archiveResolver,
) (entries []lockfile.Archive, err error) {
	var errs []error
	for _, tool := range tools {
		capability, found := capabilities[tool.ID]
		if !found {
			errs = append(errs, fmt.Errorf("missing tool capability %q", tool.ID))
			continue
		}

		install, ok := capability.Install.(toolchain.ArchiveInstall)
		if !ok {
			continue
		}

		entry, resolveErr := resolveArchiveStep(writer, tool, install, resolvePlatform)
		if resolveErr != nil {
			errs = append(errs, resolveErr)
			continue
		}

		entries = append(entries, entry)
	}

	err = errors.Join(errs...)
	return entries, err
}

// platformResolver hashes one platform's archive for a tool. Extracted as a parameter so the
// assembly logic in resolveArchive is testable without network I/O.
type platformResolver func(
	writer io.Writer,
	spec toolchain.ArchiveSpec,
	tool style.Tool,
	platformKey string,
) (hash string, err error)

func resolveArchive(
	writer io.Writer,
	tool style.Tool,
	install toolchain.ArchiveInstall,
	resolveOne platformResolver,
) (archive lockfile.Archive, err error) {
	spec := install.Spec
	hashes := make(map[string]string, len(spec.Platforms))

	for platformKey := range spec.Platforms {
		hash, hashErr := resolveOne(writer, spec, tool, platformKey)
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
	url := fmt.Sprintf(spec.URLFormat, tool.PinnedVersion, platform)

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

	return hashFile(archive)
}
