package cli

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/testutil"
)

/* ---------------------------------------- Explain Text ---------------------------------------- */

func TestExplainActiveRule(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "rule:profile/enforcement-levels",
	})
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr %q)", exitCode, stderr.String())
	}

	output := stdout.String()
	for _, snippet := range []string{
		"Rule: profile/enforcement-levels",
		"pack:",
		"project",
		"enforcement:",
		"scope:",
		"requirements:",
		"check:",
		"Pack config",
		"commands",
	} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("explain output missing %q:\n%s", snippet, output)
		}
	}
}

/* ------------------------------------- Explain Validation ------------------------------------- */

func TestExplainRejectsMissingSubject(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{"explain", "--repo-root", repo})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), "expected one subject") {
		t.Fatalf("expected missing-subject error, got %q", stderr.String())
	}
}

func TestExplainRejectsSubjectWithoutKind(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "profile/enforcement-levels",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), "invalid subject") {
		t.Fatalf("expected invalid-subject error, got %q", stderr.String())
	}
}

func TestExplainRejectsUnsupportedKind(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "pack:project",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), "unsupported subject") {
		t.Fatalf("expected unsupported-subject error, got %q", stderr.String())
	}
}

func TestExplainRejectsEmptyRuleID(t *testing.T) {
	t.Parallel()

	tool, _, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{"explain", "--repo-root", repo, "rule:"})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d for empty rule id, got %d", usageExitCode, exitCode)
	}
}

func TestExplainUnknownRuleFails(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "rule:does/not-exist",
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for unknown rule, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "unknown rule") {
		t.Fatalf("expected unknown-rule error, got %q", stderr.String())
	}
}

func TestExplainDeclaredInactiveRuleFails(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	// go/file-shape is declared by the Go Pack but is not bound in this profile.
	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "rule:go/file-shape",
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for inactive rule, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "not active") {
		t.Fatalf("expected not-active error, got %q", stderr.String())
	}
}

func TestExplainJSONEnvelope(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--format", "json", "--repo-root", repo, "rule:profile/enforcement-levels",
	})
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr %q)", exitCode, stderr.String())
	}

	envelope := assertResultEnvelope(t, stdout.Bytes(), "explain")

	var payload struct {
		Rule struct {
			ID           string   `json:"id"`
			Pack         string   `json:"pack"`
			Name         string   `json:"name"`
			External     bool     `json:"external"`
			Enforcement  string   `json:"enforcement"`
			Scope        string   `json:"scope"`
			Requirements []string `json:"requirements"`
			Check        struct {
				Category string `json:"category"`
			} `json:"check"`
		} `json:"rule"`
	}
	if err := json.Unmarshal(envelope.Result, &payload); err != nil {
		t.Fatalf("decode explain result: %v\n%s", err, envelope.Result)
	}

	rule := payload.Rule
	if rule.ID != "profile/enforcement-levels" {
		t.Fatalf("id = %q, want profile/enforcement-levels", rule.ID)
	}
	if rule.Pack != "project" {
		t.Fatalf("pack = %q, want project", rule.Pack)
	}
	if rule.Enforcement == "" || rule.Scope == "" {
		t.Fatalf("active rule missing enforcement/scope: %+v", rule)
	}
	if len(rule.Requirements) == 0 {
		t.Fatalf("active rule missing requirements: %+v", rule)
	}
	if rule.Check.Category == "" {
		t.Fatalf("check execution category missing: %+v", rule)
	}
}

func TestExplainJSONUnknownRuleIsOperationFailed(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--format", "json", "--repo-root", repo, "rule:does/not-exist",
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "explain", report.ErrorCodeOperationFailed)
}

func TestExplainJSONInvalidSubjectIsInvalidArgument(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--format", "json", "--repo-root", repo, "pack:project",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "explain", report.ErrorCodeInvalidArgument)
}

func TestExplainLaunchesNoToolProcess(t *testing.T) {
	t.Parallel()

	// explain is metadata-only: it must succeed for a rule whose check would otherwise require an
	// installed tool, proving no tool inspection or subprocess is launched.
	tool, stdout, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"explain", "--repo-root", repo, "rule:go/lint",
	})
	if exitCode != 0 {
		t.Fatalf("expected exit 0 for metadata-only explain, got %d (stderr %q)",
			exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Rule: go/lint") {
		t.Fatalf("expected go/lint explanation, got %q", stdout.String())
	}
}
