package drivers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/style"
)

/* --------------------------------------- External Driver -------------------------------------- */

// TestExternalCheckDriver is the end-to-end acceptance test for the external Pack subprocess
// protocol. It compiles a real helper binary that speaks the JSONL protocol, then drives every
// success and failure path through the flat external driver: valid diagnostics, completion failure,
// malformed output, missing completion, invalid range, nonzero exit, timeout, and truncation all
// become structured results or safe execution errors.
func TestExternalCheckDriver(t *testing.T) {
	t.Parallel()

	helper := compilePackHelper(t)
	packDir := stagePackBinary(t, helper)

	run := execution.RunContext{
		RepoRoot: t.TempDir(),
		Scope:    "all",
		Profile:  policy.Profile{},
	}
	driver := externalCheckDriver()
	ctx := context.Background()

	t.Run("diagnostic flows through as a structured result", func(t *testing.T) {
		t.Parallel()

		result, err := driver(ctx, run, externalJob(packDir, "diagnostic", 5*time.Second), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.PackID != "testpack" {
			t.Fatalf("expected pack id provenance, got %q", result.PackID)
		}
		if len(result.Diagnostics) != 1 {
			t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
		}
		diag := result.Diagnostics[0]
		if diag.File != "internal/service/users.go" || diag.Range.Start.Line != 42 {
			t.Fatalf("unexpected diagnostic: %+v", diag)
		}
	})

	t.Run("completion failure becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "fail-completion", 5*time.Second), nil)
		if err == nil {
			t.Fatal("expected completion failure error")
		}
	})

	t.Run("malformed output becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "malformed", 5*time.Second), nil)
		if err == nil {
			t.Fatal("expected malformed output error")
		}
	})

	t.Run("missing completion becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "no-completion", 5*time.Second), nil)
		if err == nil {
			t.Fatal("expected missing completion error")
		}
	})

	t.Run("invalid range becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "bad-range", 5*time.Second), nil)
		if err == nil {
			t.Fatal("expected invalid range error")
		}
	})

	t.Run("nonzero exit becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "nonzero", 5*time.Second), nil)
		if err == nil {
			t.Fatal("expected nonzero exit error")
		}
	})

	t.Run("timeout becomes an execution error", func(t *testing.T) {
		t.Parallel()
		result, err := driver(ctx, run, externalJob(packDir, "timeout", 200*time.Millisecond), nil)
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if !result.TimedOut {
			t.Fatal("expected TimedOut to be set on the result")
		}
	})

	t.Run("output truncation becomes an execution error", func(t *testing.T) {
		t.Parallel()
		result, err := driver(ctx, run, externalJob(packDir, "truncate", 10*time.Second), nil)
		if err == nil {
			t.Fatal("expected truncation error")
		}
		if !result.Truncated {
			t.Fatal("expected Truncated to be set on the result")
		}
	})

	t.Run("stderr is captured but does not fail a clean run", func(t *testing.T) {
		t.Parallel()
		_, err := driver(ctx, run, externalJob(packDir, "stderr-debug", 5*time.Second), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestExternalCheckRequestShape asserts the request contract: configuration is an empty JSON object
// (never null) when no pack config exists, and files are repository-relative slash paths (never
// absolute). The helper echoes both back so the assertions cover the wire format the subprocess
// receives.
func TestExternalCheckRequestShape(t *testing.T) {
	t.Parallel()

	helper := compilePackHelper(t)
	packDir := stagePackBinary(t, helper)

	repoRoot := t.TempDir()
	srcDir := filepath.Join(repoRoot, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	srcFile := filepath.Join(srcDir, "file.go")
	if err := os.WriteFile(srcFile, []byte("package src\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	run := execution.RunContext{
		RepoRoot: repoRoot,
		Scope:    "all",
		Profile: policy.Profile{
			Repository: policy.RepositoryConfig{
				ScopeRoots: map[style.Scope][]string{"all": {"."}},
			},
		},
	}

	result, err := externalCheckDriver()(context.Background(), run,
		externalJob(packDir, "inspect-request", 5*time.Second), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
	}

	message := result.Diagnostics[0].Message
	if !strings.Contains(message, "config={}") {
		t.Fatalf("expected empty configuration object {}, got %q", message)
	}
	if !strings.Contains(message, "files=src/file.go") {
		t.Fatalf("expected repository-relative file path, got %q", message)
	}
	if strings.Contains(message, repoRoot) {
		t.Fatalf("request leaked an absolute repository path: %q", message)
	}

	if !strings.Contains(message, "filenil=false") {
		t.Fatalf("expected files to be an empty array [] not null, got %q", message)
	}
}

func externalJob(
	packDirectory string,
	checkID string,
	timeout time.Duration,
) (job style.ExternalCheckJob) {
	return style.ExternalCheckJob{
		PackID:        "testpack",
		RuleID:        "testpack/rule",
		CheckID:       checkID,
		PackDirectory: packDirectory,
		Command:       "bin/packhelper",
		Timeout:       timeout,
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func stagePackBinary(t *testing.T, source string) (packDirectory string) {
	t.Helper()

	packDirectory = t.TempDir()
	binDirectory := filepath.Join(packDirectory, "bin")
	if err := os.MkdirAll(binDirectory, 0o755); err != nil {
		t.Fatal(err)
	}

	destination := filepath.Join(binDirectory, "packhelper")
	if runtime.GOOS == "windows" {
		destination += ".exe"
	}

	content, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, content, 0o755); err != nil {
		t.Fatal(err)
	}

	return packDirectory
}

func compilePackHelper(t *testing.T) (binary string) {
	t.Helper()

	source := filepath.Join("testdata", "packhelper", "main.go")
	binary = filepath.Join(t.TempDir(), "packhelper")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	build := exec.Command("go", "build", "-o", binary, source)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("compile pack helper: %v\n%s", err, output)
	}

	return binary
}
