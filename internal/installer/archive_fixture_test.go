package installer

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"

	"github.com/ulikunitz/xz"
)

/* --------------------------------------- Archive Entries -------------------------------------- */

type archiveEntry struct {
	Name     string
	Body     string
	Typeflag byte
	Linkname string
}

/* --------------------------------------- Fixture Writers -------------------------------------- */

func writeTestArchive(
	t *testing.T,
	entries ...archiveEntry,
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

func writeTestArchiveHeader(t *testing.T, name string, size int64) (path string) {
	t.Helper()

	path = filepath.Join(t.TempDir(), "oversized.tar.xz")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}

	xzWriter, err := xz.NewWriter(file)
	if err != nil {
		t.Fatalf("create xz writer: %v", err)
	}
	tarWriter := tar.NewWriter(xzWriter)
	if err = tarWriter.WriteHeader(&tar.Header{
		Name:     name,
		Mode:     0o755,
		Size:     size,
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("write tar header: %v", err)
	}

	// leave the declared body absent so the extractor must reject from the header
	if err = xzWriter.Close(); err != nil {
		t.Fatalf("close xz writer: %v", err)
	}
	if err = file.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}
	return path
}
