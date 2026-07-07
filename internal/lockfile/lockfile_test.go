package lockfile

import (
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
