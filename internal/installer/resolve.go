package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/lockfile"
	"github.com/wbd2023/Quill/internal/toolchain"
)

/* ----------------------------------------- Resolution ----------------------------------------- */

// Resolve downloads every platform archive for each archive-installed tool declared in the
// profile, hashes each one, and returns the entries that make up quill.lock. Failures from
// independent tools or platforms are collected and returned as a joined error.
func Resolve(
	ctx context.Context,
	writer io.Writer,
	tools []toolchain.Tool,
) (entries []lockfile.Archive, err error) {
	return resolveWith(ctx, writer, tools, resolveArchive)
}

// archiveResolver assembles one tool's lockfile entry by resolving each platform. Extracted as a
// parameter so the iteration and filtering in resolveWith is testable without network I/O.
type archiveResolver func(
	ctx context.Context,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.GitHubInstall,
	resolveOne platformResolver,
) (archive lockfile.Archive, err error)

func resolveWith(
	ctx context.Context,
	writer io.Writer,
	tools []toolchain.Tool,
	resolveArchiveStep archiveResolver,
) (entries []lockfile.Archive, err error) {
	var errs []error
	for _, tool := range tools {
		install, ok := tool.Install.(toolchain.GitHubInstall)
		if !ok {
			continue
		}

		entry, resolveErr := resolveArchiveStep(ctx, writer, tool, install, resolvePlatform)
		if resolveErr != nil {
			errs = append(errs, resolveErr)
			continue
		}

		entries = append(entries, entry)
	}

	return entries, errors.Join(errs...)
}

// platformResolver hashes one platform's archive for a tool. Extracted as a parameter so the
// assembly logic in resolveArchive is testable without network I/O.
type platformResolver func(
	ctx context.Context,
	writer io.Writer,
	install toolchain.GitHubInstall,
	tool toolchain.Tool,
	platformKey string,
) (hash string, err error)

func resolveArchive(
	ctx context.Context,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.GitHubInstall,
	resolveOne platformResolver,
) (archive lockfile.Archive, err error) {
	hashes := make(map[string]string, len(install.Platforms))

	for platformKey := range install.Platforms {
		hash, hashErr := resolveOne(ctx, writer, install, tool, platformKey)
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
	ctx context.Context,
	writer io.Writer,
	install toolchain.GitHubInstall,
	tool toolchain.Tool,
	platformKey string,
) (hash string, err error) {
	platform := install.Platforms[platformKey]
	tag := fmt.Sprintf(install.Tag, tool.PinnedVersion)
	asset := fmt.Sprintf(install.Asset, tag, platform)
	url := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		install.Owner,
		install.Repository,
		tag,
		asset,
	)

	dir, err := os.MkdirTemp("", "quill-resolve-*")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	archive := filepath.Join(dir, asset)
	if _, err = fmt.Fprintf(writer, "Resolving %s %s...\n", tool.Name, platformKey); err != nil {
		return "", err
	}

	if err = downloadFile(ctx, url, archive); err != nil {
		return "", err
	}

	return hashFile(archive)
}
