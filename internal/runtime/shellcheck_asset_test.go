package runtime

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ulikunitz/xz"
)

/* -------------------------------------- ShellCheck Assets ------------------------------------- */

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

/* ------------------------------------- ShellCheck Archives ------------------------------------ */

func TestExtractShellcheckBinaryExtractsExpectedAsset(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/LICENSE.txt",
			Body: "licence",
		},
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/README.txt",
			Body: "readme",
		},
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/shellcheck",
			Body: "#!/bin/sh\n",
		},
	)

	destination := t.TempDir()
	binaryPath, err := extractShellcheckBinary(archivePath, destination, "0.10.0")
	if err != nil {
		t.Fatalf("extractShellcheckBinary: %v", err)
	}

	contents, err := os.ReadFile(binaryPath)
	if err != nil {
		t.Fatalf("read extracted binary: %v", err)
	}

	if string(contents) != "#!/bin/sh\n" {
		t.Fatalf("unexpected extracted contents: %q", contents)
	}

	readmePath := filepath.Join(destination, "shellcheck-v0.10.0", "README.txt")
	if _, err := os.Stat(readmePath); err == nil {
		t.Fatal("expected README to be ignored, not extracted")
	}
}

func TestExtractShellcheckBinaryRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/../shellcheck",
			Body: "bad",
		},
	)

	if _, err := extractShellcheckBinary(archivePath, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected path traversal entry to fail")
	}
}

func TestExtractShellcheckBinaryRejectsLinks(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name:     "shellcheck-v0.10.0/shellcheck",
			Typeflag: tar.TypeSymlink,
			Linkname: "/bin/sh",
		},
	)

	if _, err := extractShellcheckBinary(archivePath, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected link entry to fail")
	}
}

/* --------------------------------------- Archive Helpers -------------------------------------- */

type shellcheckArchiveEntry struct {
	Name     string
	Body     string
	Typeflag byte
	Linkname string
}

func writeShellcheckArchive(
	t *testing.T,
	entries ...shellcheckArchiveEntry,
) (path string) {
	t.Helper()

	path = filepath.Join(t.TempDir(), "shellcheck.tar.xz")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}

	xzWriter, err := xz.NewWriter(file)
	if err != nil {
		t.Fatalf("create xz writer: %v", err)
	}

	tarWriter := tar.NewWriter(xzWriter)
	for _, entry := range entries {
		typeflag := entry.Typeflag
		if typeflag == 0 {
			typeflag = tar.TypeReg
		}

		header := &tar.Header{
			Name:     entry.Name,
			Mode:     0o755,
			Size:     int64(len(entry.Body)),
			Typeflag: typeflag,
			Linkname: entry.Linkname,
		}
		if typeflag != tar.TypeReg {
			header.Size = 0
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("write tar header: %v", err)
		}

		if header.Size > 0 {
			if _, err := tarWriter.Write([]byte(entry.Body)); err != nil {
				t.Fatalf("write tar body: %v", err)
			}
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}

	if err := xzWriter.Close(); err != nil {
		t.Fatalf("close xz writer: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}

	return path
}
