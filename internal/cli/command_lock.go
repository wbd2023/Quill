package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/engine"
	"github.com/wbd2023/Quill/internal/lockfile"
)

/* ------------------------------------------- Command ------------------------------------------ */

func runLock(tool Tool, options lockOptions) (exitCode int) {
	engineInstance, err := engine.New(
		options.repoRoot,
		engine.WithProgressWriter(tool.stdout),
	)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	result, err := engineInstance.Lock(context.Background())
	if err != nil {
		tool.writeError(err)
		return 1
	}

	archiveByID := make(map[string]lockfile.Archive, len(result.Archives))
	for _, archive := range result.Archives {
		archiveByID[archive.Tool] = archive
	}

	contents, err := lockfile.Encode(lockfile.Lockfile{Archives: archiveByID})
	if err != nil {
		tool.writeError(err)
		return 1
	}

	path := filepath.Join(options.repoRoot, lockfile.DefaultFilename)
	if err = writeLockfile(path, contents); err != nil {
		tool.writeError(err)
		return 1
	}

	if _, err = fmt.Fprintf(
		tool.stdout, "Wrote %s (%d tools)\n", path, len(result.Archives),
	); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}

const (
	standardDirectoryPermissions os.FileMode = 0o755
	standardLockfilePermissions  os.FileMode = 0o644
)

// writeLockfile writes contents to path atomically via a temp-file rename in the same directory.
func writeLockfile(path string, contents string) (err error) {
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, standardDirectoryPermissions); err != nil {
		return err
	}

	temp, err := os.CreateTemp(dir, ".lock-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()

	defer func() {
		if err != nil {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err = temp.WriteString(contents); err != nil {
		_ = temp.Close() // preserve the original write error
		return err
	}

	if err = temp.Chmod(standardLockfilePermissions); err != nil {
		_ = temp.Close() // preserve the original permission error
		return err
	}

	if err = temp.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

/* ------------------------------------------- Parsing ------------------------------------------ */

func parseLockOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options lockOptions, err error) {
	const summary = "resolve and write archive-tool hashes to quill.lock"
	flagSet := newLockFlagSet(&options)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newLockFlagSet(options *lockOptions) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("lock")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	return flagSet
}

func lockUsageText() (usage string) {
	const summary = "resolve and write archive-tool hashes to quill.lock"
	var options lockOptions
	return commandUsage("lock", summary, newLockFlagSet(&options))
}
