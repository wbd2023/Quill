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
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

/* --------------------------------------- External Driver -------------------------------------- */

// TestExternalCheckDriver is the end-to-end acceptance test for the external Pack subprocess
// protocol. It compiles a real helper binary that speaks the JSONL protocol, then drives every
// success and failure path through the flat external driver: valid diagnostics, completion failure,
// malformed output, missing completion, invalid range, nonzero exit, timeout, and truncation all
// become structured results or safe execution errors. Pack and rule provenance for the request is
// supplied by the Rule, never by the ExternalCheck Job.
func TestExternalCheckDriver(t *testing.T) {

	helper := compilePackHelper(t)
	packDir := stagePackBinary(t, helper)

	run := execution.RunContext{
		RepoRoot: t.TempDir(),
		Scope:    "all",
		Profile:  profile.Profile{},
	}
	driver := externalCheckDriver()
	ctx := context.Background()

	t.Run("diagnostic flows through as a structured result", func(t *testing.T) {
		t.Parallel()

		result, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "diagnostic", 5*time.Second), nil,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
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
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "fail-completion", 5*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected completion failure error")
		}
	})

	t.Run("malformed output becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "malformed", 5*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected malformed output error")
		}
	})

	t.Run("missing completion becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "no-completion", 5*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected missing completion error")
		}
	})

	t.Run("invalid range becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "bad-range", 5*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected invalid range error")
		}
	})

	t.Run("nonzero exit becomes an execution error", func(t *testing.T) {
		t.Parallel()
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "nonzero", 5*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected nonzero exit error")
		}
	})

	t.Run("timeout becomes an execution error", func(t *testing.T) {
		t.Parallel()
		result, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "timeout", 200*time.Millisecond), nil,
		)
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if !result.TimedOut {
			t.Fatal("expected TimedOut to be set on the result")
		}
	})

	t.Run("output truncation becomes an execution error", func(t *testing.T) {
		t.Parallel()
		result, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "truncate", 10*time.Second), nil,
		)
		if err == nil {
			t.Fatal("expected truncation error")
		}
		if !result.Truncated {
			t.Fatal("expected Truncated to be set on the result")
		}
	})

	t.Run("stderr is captured but does not fail a clean run", func(t *testing.T) {
		t.Parallel()
		_, err := driver(
			ctx, run, externalRule(), externalJob(packDir, "stderr-debug", 5*time.Second), nil,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestExternalCheckRequestShape asserts the request contract: policy is an empty JSON object
// (never null) when no Pack Policy exists, files are repository-relative slash paths (never
// absolute), and Pack/Rule provenance is derived from the Rule, not the Job. The helper echoes all
// of these back so the assertions cover the wire format the subprocess receives.
func TestExternalCheckRequestShape(t *testing.T) {

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
		Profile: profile.Profile{
			Repository: profile.RepositoryConfig{
				ScopeRoots: map[style.Scope][]string{"all": {"."}},
			},
		},
	}

	driver := externalCheckDriver()
	result, err := driver(
		context.Background(),
		run,
		externalRule(),
		externalJob(packDir, "inspect-request", 5*time.Second),
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
	}

	message := result.Diagnostics[0].Message
	if !strings.Contains(message, "pack=testpack") {
		t.Fatalf("expected pack provenance from the rule, got %q", message)
	}
	if !strings.Contains(message, "rule=testpack/rule") {
		t.Fatalf("expected rule provenance from the rule, got %q", message)
	}
	if !strings.Contains(message, "policy={}") {
		t.Fatalf("expected empty policy object {}, got %q", message)
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

func TestStderrContextPreservesMessageBeforeTrailingWhitespace(t *testing.T) {
	t.Parallel()

	stderr := "external Pack failed" + strings.Repeat(" ", stderrExcerpt+1)
	if got, want := stderrContext(stderr), "\nexternal Pack failed"; got != want {
		t.Fatalf("stderrContext() = %q, want %q", got, want)
	}
}

func externalRule() (rule style.Rule) {
	return style.Rule{ID: "testpack/rule", PackID: "testpack"}
}

func externalJob(
	packDirectory string,
	checkID string,
	timeout time.Duration,
) (job style.ExternalCheck) {
	return style.ExternalCheck{
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
