package cli

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/testutil"
)

/* ------------------------------------------ List Text ----------------------------------------- */

func TestListPacksReportsAvailablePacks(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{"list", "--repository-root", repo, "packs"})
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr %q)", exitCode, stderr.String())
	}

	// The Quill repository enables every Shipped Pack, so each is active and shipped.
	output := stdout.String()
	for _, snippet := range []string{"Packs", "project", "shipped", "active"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("list packs output missing %q:\n%s", snippet, output)
		}
	}
}

func TestListPacksDistinguishesActiveFromInactive(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	root := t.TempDir()

	if exitCode := tool.Run(context.Background(), []string{
		"init", "--repository-root", root,
	}); exitCode != 0 {
		t.Fatalf("init failed: %d (stderr %q)", exitCode, stderr.String())
	}

	// The minimal profile enables only the project Pack, so every other shipped Pack is
	// available but inactive: list must show both states.
	stdout.Reset()
	if exitCode := tool.Run(context.Background(), []string{
		"list", "--repository-root", root, "packs",
	}); exitCode != 0 {
		t.Fatalf("list packs failed: %d (stderr %q)", exitCode, stderr.String())
	}

	output := stdout.String()
	for _, snippet := range []string{"active", "inactive", "project"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("minimal-repo list packs missing %q:\n%s", snippet, output)
		}
	}
}

func TestListRulesShowsActiveBindings(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	_ = tool.Run(context.Background(), []string{"list", "--repository-root", repo, "rules"})

	output := stdout.String()
	for _, snippet := range []string{
		"Rules", "Pack", "Provenance", "Enforcement", "Scope", "active", "required", "shipped",
	} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("list rules output missing %q:\n%s", snippet, output)
		}
	}
}

func TestListToolsReportsPinnedTools(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	_ = tool.Run(context.Background(), []string{"list", "--repository-root", repo, "tools"})

	output := stdout.String()
	for _, snippet := range []string{"Tools", "Pin", "Packs"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("list tools output missing %q:\n%s", snippet, output)
		}
	}
	if strings.Contains(output, "Source") {
		t.Fatalf("list tools must not imply unsupported Tool provenance:\n%s", output)
	}
}

func TestListScopesReportsConfiguredScopes(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	_ = tool.Run(context.Background(), []string{"list", "--repository-root", repo, "scopes"})

	output := stdout.String()
	for _, snippet := range []string{"Scopes", "Default", "Roots"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("list scopes output missing %q:\n%s", snippet, output)
		}
	}
}

/* --------------------------------------- List Validation -------------------------------------- */

func TestListRejectsMissingSelector(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{"list", "--repository-root", repo})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), `expected "<selector>"`) {
		t.Fatalf("expected missing-selector error, got %q", stderr.String())
	}
}

func TestListRejectsInvalidSelector(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"list", "--repository-root", repo, "widgets",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), `<selector> must be one of`) {
		t.Fatalf("expected invalid-selector error, got %q", stderr.String())
	}
}

func TestListRejectsMultipleSelectors(t *testing.T) {
	t.Parallel()

	tool, _, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"list", "--repository-root", repo, "packs", "rules",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d for multiple selectors, got %d", usageExitCode, exitCode)
	}
}

/* ------------------------------------- List JSON Contract ------------------------------------- */

func TestListRulesJSONEnvelope(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"list", "rules", "--format", "json", "--repository-root", repo,
	})
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr %q)", exitCode, stderr.String())
	}

	envelope := assertResultEnvelope(t, stdout.Bytes(), "list")

	var payload struct {
		Rules []struct {
			ID          string `json:"id"`
			Pack        string `json:"pack"`
			Provenance  string `json:"provenance"`
			Name        string `json:"name"`
			Active      bool   `json:"active"`
			Enforcement string `json:"enforcement"`
			Scope       string `json:"scope"`
			Fix         bool   `json:"fix"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(envelope.Result, &payload); err != nil {
		t.Fatalf("decode list rules result: %v\n%s", err, envelope.Result)
	}

	if len(payload.Rules) == 0 {
		t.Fatalf("list rules carried no rules: %s", envelope.Result)
	}

	var sawActive bool
	for _, rule := range payload.Rules {
		if rule.ID == "" || rule.Pack == "" || rule.Provenance == "" || rule.Name == "" {
			t.Fatalf("rule missing identity or provenance fields: %+v", rule)
		}
		if rule.Active {
			sawActive = true
			if rule.Enforcement == "" || rule.Scope == "" {
				t.Fatalf("active rule %q missing enforcement or scope: %+v", rule.ID, rule)
			}
		}
	}

	if !sawActive {
		t.Fatalf("expected at least one active rule, got %s", envelope.Result)
	}
}

func TestListJSONStdoutIsSingleEnvelope(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	_ = tool.Run(context.Background(), []string{
		"list", "--format", "json", "--repository-root", repo, "packs",
	})

	if !json.Valid(stdout.Bytes()) {
		t.Fatalf("list JSON stdout must be exactly one JSON document, got:\n%s", stdout.Bytes())
	}
}

func TestListJSONInvalidSelectorIsInvalidArgument(t *testing.T) {
	t.Parallel()

	tool, stdout, _ := newTestCLI()
	repo := testutil.RepositoryRoot(t)

	exitCode := tool.Run(context.Background(), []string{
		"list", "--format", "json", "--repository-root", repo, "widgets",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "list", report.ErrorCodeInvalidArgument)
}
