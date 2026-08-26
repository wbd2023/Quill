package lockfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestEncodeDecodeRoundtrip(t *testing.T) {
	t.Parallel()

	original := Lockfile{
		Loaded: true,
		Archives: map[string]Archive{
			"shellcheck": {
				Tool:    "shellcheck",
				Version: "0.10.0",
				Hashes: map[string]string{
					"darwin/amd64": "ef27684f",
					"linux/amd64":  "6c881ab0",
				},
			},
		},
	}

	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	archive, ok := decoded.Archives["shellcheck"]
	if !ok {
		t.Fatal("expected shellcheck archive entry")
	}

	if archive.Version != "0.10.0" {
		t.Fatalf("version = %q, want 0.10.0", archive.Version)
	}

	if len(archive.Hashes) != 2 {
		t.Fatalf("expected 2 hashes, got %d", len(archive.Hashes))
	}
}

func TestHashForDistinctErrors(t *testing.T) {
	t.Parallel()

	loaded := Lockfile{
		Loaded: true,
		Archives: map[string]Archive{
			"shellcheck": {
				Tool:    "shellcheck",
				Version: "0.10.0",
				Hashes:  map[string]string{"linux/amd64": "abc"},
			},
		},
	}

	tests := []struct {
		name       string
		lockfile   Lockfile
		tool       string
		version    string
		goos       string
		goarch     string
		wantSubstr string
	}{
		{
			name:       "not loaded",
			lockfile:   Lockfile{Loaded: false},
			tool:       "shellcheck",
			version:    "0.10.0",
			goos:       "linux",
			goarch:     "amd64",
			wantSubstr: "quill.lock not found",
		},
		{
			name:       "tool missing",
			lockfile:   loaded,
			tool:       "nonexistent",
			version:    "0.10.0",
			goos:       "linux",
			goarch:     "amd64",
			wantSubstr: "no lockfile entry for nonexistent",
		},
		{
			name:       "version mismatch",
			lockfile:   loaded,
			tool:       "shellcheck",
			version:    "0.11.0",
			goos:       "linux",
			goarch:     "amd64",
			wantSubstr: "lockfile has shellcheck 0.10.0 but profile pins 0.11.0",
		},
		{
			name:       "platform missing",
			lockfile:   loaded,
			tool:       "shellcheck",
			version:    "0.10.0",
			goos:       "darwin",
			goarch:     "amd64",
			wantSubstr: "no lockfile hash for shellcheck on darwin/amd64",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := test.lockfile.HashFor(test.tool, test.version, test.goos, test.goarch)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.wantSubstr) {
				t.Fatalf("error %q does not contain %q", err.Error(), test.wantSubstr)
			}
		})
	}
}

func TestHashForReturnsHash(t *testing.T) {
	t.Parallel()

	lockfile := Lockfile{
		Loaded: true,
		Archives: map[string]Archive{
			"shellcheck": {
				Tool:    "shellcheck",
				Version: "0.10.0",
				Hashes:  map[string]string{"linux/amd64": "abc123"},
			},
		},
	}

	hash, err := lockfile.HashFor("shellcheck", "0.10.0", "linux", "amd64")
	if err != nil {
		t.Fatalf("HashFor: %v", err)
	}

	if hash != "abc123" {
		t.Fatalf("hash = %q, want abc123", hash)
	}
}

func TestDecodeRejectsUnknownSchemaVersion(t *testing.T) {
	t.Parallel()

	source := `schema_version = 2`
	if _, err := Decode(source); err == nil {
		t.Fatal("expected schema version error")
	}
}

func TestDecodeRejectsDuplicateTool(t *testing.T) {
	t.Parallel()

	source := `
schema_version = 1

[[archive]]
tool = "shellcheck"
version = "0.10.0"
hashes = {}

[[archive]]
tool = "shellcheck"
version = "0.11.0"
hashes = {}
`
	if _, err := Decode(source); err == nil {
		t.Fatal("expected duplicate tool error")
	}
}

func TestLoadRejectsSymlinkLockfile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "elsewhere")
	if err := os.WriteFile(target, []byte("schema_version = 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, DefaultFilename)); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	lockfile, err := Load(root)
	if err == nil {
		t.Fatal("expected error for symlink lockfile, got nil")
	}
	if lockfile.Loaded {
		t.Fatalf("symlink lockfile must not load: %+v", lockfile)
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %q does not mention non-regular file", err.Error())
	}
}

func TestLoadAcceptsLockfileAtByteLimit(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(root, DefaultFilename),
		[]byte(validLockfileOfSize(int(maxLockfileBytes))),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	lockfile, err := Load(root)
	if err != nil {
		t.Fatalf("Load at byte limit: %v", err)
	}
	if !lockfile.Loaded {
		t.Fatal("expected lockfile at byte limit to load")
	}
}

func TestLoadRejectsLockfileOverByteLimit(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(root, DefaultFilename),
		[]byte(validLockfileOfSize(int(maxLockfileBytes)+1)),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	lockfile, err := Load(root)
	if err == nil {
		t.Fatal("expected error for over-limit lockfile, got nil")
	}
	if lockfile.Loaded {
		t.Fatalf("over-limit lockfile must not load: %+v", lockfile)
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("error %q does not mention the size limit", err.Error())
	}
}

// validLockfileOfSize returns valid lockfile TOML of exactly size bytes by padding a minimal
// schema with short comment lines, letting the byte-limit tests hit the exact boundary.
func validLockfileOfSize(size int) (toml string) {
	const header = "schema_version = 1\n"
	if size <= len(header) {
		return header
	}

	var builder strings.Builder
	builder.Grow(size)
	builder.WriteString(header)

	// Pad with short comment lines so no single line approaches a parser
	// line-length assumption; the result stays valid TOML of exactly size bytes.
	remaining := size - len(header)
	const line = "#x\n"
	for remaining >= len(line) {
		builder.WriteString(line)
		remaining -= len(line)
	}
	if remaining > 0 {
		builder.WriteByte('#')
		for index := 1; index < remaining; index++ {
			builder.WriteByte('x')
		}
	}

	return builder.String()
}
