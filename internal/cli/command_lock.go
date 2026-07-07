package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/installer"
	"ciphera/tools/internal/lockfile"
)

/* ------------------------------------------- Command ------------------------------------------ */

func runLock(tool Tool, options lockOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, "")
	if err != nil {
		tool.writeError(err)
		return 1
	}

	entries, err := installer.Resolve(
		tool.stdout,
		context.Effective.Tools,
		context.ToolCapabilities,
	)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	archiveByID := make(map[string]lockfile.Archive, len(entries))
	for _, entry := range entries {
		archiveByID[entry.Tool] = entry
	}

	contents, err := lockfile.Encode(lockfile.Lockfile{Archives: archiveByID})
	if err != nil {
		tool.writeError(err)
		return 1
	}

	path := filepath.Join(options.repoRoot, lockfile.DefaultFilename)
	if err = writeFile(path, contents); err != nil {
		tool.writeError(err)
		return 1
	}

	if _, err = fmt.Fprintf(tool.stdout, "Wrote %s (%d tools)\n", path, len(entries)); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}

const standardLockPermissions os.FileMode = 0o755

// writeFile writes contents to path atomically via a temp-file rename in the same directory.
func writeFile(path string, contents string) (err error) {
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, standardLockPermissions); err != nil {
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
		_ = temp.Close()
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
