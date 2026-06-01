package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellcheckAssetName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{goos: "linux", goarch: "amd64", expected: "linux.x86_64"},
		{goos: "linux", goarch: "arm64", expected: "linux.aarch64"},
		{goos: "darwin", goarch: "amd64", expected: "darwin.x86_64"},
		{goos: "darwin", goarch: "arm64", expected: "darwin.aarch64"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.goos+"-"+testCase.goarch, func(t *testing.T) {
			t.Parallel()

			actual, err := shellcheckAssetName(testCase.goos, testCase.goarch)
			if err != nil {
				t.Fatalf("shellcheckAssetName: %v", err)
			}

			if actual != testCase.expected {
				t.Fatalf("asset name = %q, want %q", actual, testCase.expected)
			}
		})
	}
}

func TestShellcheckAssetNameRejectsUnsupportedPlatform(t *testing.T) {
	t.Parallel()

	if _, err := shellcheckAssetName("freebsd", "amd64"); err == nil {
		t.Fatal("expected unsupported platform error")
	}
}

func TestVerifyFileChecksumRejectsMismatch(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "archive.tar.xz")
	if err := os.WriteFile(path, []byte("archive"), 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}

	err := verifyFileChecksum(path, "archive.tar.xz", strings.Repeat("0", 64))
	if err == nil {
		t.Fatal("expected checksum mismatch")
	}
}
