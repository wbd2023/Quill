package lockfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	standardDirectoryPermissions os.FileMode = 0o755
	standardLockfilePermissions  os.FileMode = 0o644
)

// Write encodes lockfile and atomically persists it to <root>/quill.lock via a
// temp-file rename in the same directory. It owns directory creation, encoding,
// temporary-file cleanup on failure, the shared 0644 lockfile permission, and
// the final replacement. Cancellation before the final replacement leaves any
// existing lockfile unchanged. It returns the absolute path of the written
// lockfile.
func Write(ctx context.Context, root string, lockfile Lockfile) (path string, err error) {
	if err = ctx.Err(); err != nil {
		return "", err
	}

	contents, err := Encode(lockfile)
	if err != nil {
		return "", fmt.Errorf("encode lockfile: %w", err)
	}

	path = filepath.Join(root, DefaultFilename)
	if err = writeAtomically(ctx, path, contents); err != nil {
		return "", fmt.Errorf("write lockfile %q: %w", path, err)
	}

	return path, nil
}

// writeAtomically writes contents to path via a temp file in the same directory,
// then renames it into place. The temp file is removed if any step fails.
func writeAtomically(ctx context.Context, path string, contents string) (err error) {
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

	if err = ctx.Err(); err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}
