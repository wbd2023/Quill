package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellcheckAssetForReturnsCorrectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{goos: "linux", goarch: "amd64", expected: "linux.x86_64"},
		{goos: "linux", goarch: "arm64", expected: "linux.aarch64"},
		{goos: "darwin", goarch: "amd64", expected: "darwin.x86_64"},
		{goos: "darwin", goarch: "arm64", expected: "darwin.aarch64"},
	}

	for _, test := range tests {
		t.Run(test.goos+"-"+test.goarch, func(t *testing.T) {
			t.Parallel()

			asset, err := shellcheckAssetFor(test.goos, test.goarch)
			if err != nil {
				t.Fatalf("shellcheckAssetFor: %v", err)
			}

			if asset.Name != test.expected {
				t.Fatalf("asset name = %q, want %q", asset.Name, test.expected)
			}
		})
	}
}

func TestShellcheckAssetForRejectsUnsupportedPlatform(t *testing.T) {
	t.Parallel()

	if _, err := shellcheckAssetFor("freebsd", "amd64"); err == nil {
		t.Fatal("expected unsupported platform error")
	}
}

func TestVerifyFileChecksumRejectsMismatch(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "archive.tar.xz")
	if err := os.WriteFile(path, []byte("archive"), 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}

	err := verifyChecksum(path, strings.Repeat("0", 64))
	if err == nil {
		t.Fatal("expected checksum mismatch")
	}
}
