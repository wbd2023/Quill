package installer

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"testing"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func testArchiveCapability(platforms ...string) (capability toolchain.Capability) {
	platformMap := make(map[string]string, len(platforms))
	for _, p := range platforms {
		platformMap[p] = p + "-asset"
	}

	return toolchain.Capability{
		ID:          "test-tool",
		Name:        "Test",
		InstallKind: toolchain.InstallKindArchive,
		Archive: &toolchain.ArchiveSpec{
			Format:    toolchain.ArchiveFormatXz,
			Platforms: platformMap,
			URL: func(version string, platform string) string {
				return "https://example.com/" + platform
			},
			BinaryPath: func(version string) string { return "test-v" + version + "/test" },
		},
	}
}

// stubResolver returns a fixed hash per platform, exercising the assembly logic in
// resolveArchive without network I/O.
func stubResolver(
	hashesByPlatform map[string]string,
	failPlatform string,
	failErr error,
) (resolver platformResolver) {
	return func(
		_ io.Writer,
		_ toolchain.ArchiveSpec,
		_ style.Tool,
		platformKey string,
	) (hash string, err error) {
		if platformKey == failPlatform {
			return "", failErr
		}

		return hashesByPlatform[platformKey], nil
	}
}

/* -------------------------------------------- Tests ------------------------------------------- */

func TestResolveArchiveCollectsAllPlatformHashes(t *testing.T) {
	t.Parallel()

	tool := style.Tool{ID: "test-tool", Name: "Test", PinnedVersion: "1.0.0"}
	capability := testArchiveCapability("linux/amd64", "darwin/arm64")
	hashes := map[string]string{
		"linux/amd64":  "aaa",
		"darwin/arm64": "bbb",
	}

	archive, err := resolveArchive(
		io.Discard,
		tool,
		capability,
		stubResolver(hashes, "", nil),
	)
	if err != nil {
		t.Fatalf("resolveArchive: %v", err)
	}

	if archive.Tool != "test-tool" {
		t.Fatalf("tool = %q, want test-tool", archive.Tool)
	}

	if archive.Version != "1.0.0" {
		t.Fatalf("version = %q, want 1.0.0", archive.Version)
	}

	gotPlatforms := make([]string, 0, len(archive.Hashes))
	for p := range archive.Hashes {
		gotPlatforms = append(gotPlatforms, p)
	}

	sort.Strings(gotPlatforms)
	wantPlatforms := []string{"darwin/arm64", "linux/amd64"}
	if len(gotPlatforms) != 2 ||
		gotPlatforms[0] != wantPlatforms[0] ||
		gotPlatforms[1] != wantPlatforms[1] {
		t.Fatalf("platforms = %v, want %v", gotPlatforms, wantPlatforms)
	}

	if archive.Hashes["linux/amd64"] != "aaa" || archive.Hashes["darwin/arm64"] != "bbb" {
		t.Fatalf("hashes = %v, want mapped values", archive.Hashes)
	}
}

func TestResolveArchivePropagatesPlatformError(t *testing.T) {
	t.Parallel()

	tool := style.Tool{ID: "test-tool", Name: "Test", PinnedVersion: "1.0.0"}
	capability := testArchiveCapability("linux/amd64")
	platformErr := errors.New("network down")

	_, err := resolveArchive(
		io.Discard,
		tool,
		capability,
		stubResolver(nil, "linux/amd64", platformErr),
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, platformErr) {
		t.Fatalf("expected wrapped platformErr, got %v", err)
	}
}

// stubArchiveResolver returns a fixed entry for any archive tool, exercising the iteration
// and filtering in resolveWith without network I/O.
func stubArchiveResolver(
	hashesByPlatform map[string]string,
) (resolver archiveResolver) {
	return func(
		_ io.Writer,
		tool style.Tool,
		_ toolchain.Capability,
		_ platformResolver,
	) (archive lockfile.Archive, err error) {
		return lockfile.Archive{
			Tool:    tool.ID,
			Version: tool.PinnedVersion,
			Hashes:  hashesByPlatform,
		}, nil
	}
}

func TestResolveFiltersNonArchiveTools(t *testing.T) {
	t.Parallel()

	tools := []style.Tool{
		{ID: "go-binary", Name: "Go Tool", PinnedVersion: "1.0.0"},
		{ID: "archive-tool", Name: "Archive Tool", PinnedVersion: "2.0.0"},
	}
	capabilities := map[string]toolchain.Capability{
		"go-binary": {
			ID:          "go-binary",
			InstallKind: toolchain.InstallKindGoBinary,
		},
		"archive-tool": testArchiveCapability("linux/amd64"),
	}

	entries, err := resolveWith(
		io.Discard,
		tools,
		capabilities,
		stubArchiveResolver(map[string]string{"linux/amd64": "abc"}),
	)
	if err != nil {
		t.Fatalf("resolveWith: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (archive only), got %d", len(entries))
	}

	if entries[0].Tool != "archive-tool" {
		t.Fatalf("entry tool = %q, want archive-tool", entries[0].Tool)
	}
}

func TestResolveReportsMissingCapability(t *testing.T) {
	t.Parallel()

	tools := []style.Tool{
		{ID: "missing", Name: "Missing", PinnedVersion: "1.0.0"},
	}

	var buf bytes.Buffer
	_, err := resolveWith(&buf, tools, nil, stubArchiveResolver(nil))
	if err == nil {
		t.Fatal("expected missing-capability error, got nil")
	}
}
