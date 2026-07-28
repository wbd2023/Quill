package profile_test

import (
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------- Runtime Binding Completeness -------------------------------- */

// allMissingBindings reports no binding as registered.
type allMissingBindings struct{}

func (allMissingBindings) HasProfileCheck(string, string) (found bool) {
	return false
}
func (allMissingBindings) HasRepositoryScanner(string, string) (found bool) {
	return false
}
func (allMissingBindings) HasTargetCommand(string, string, string) (found bool) {
	return false
}
func (allMissingBindings) HasTargetCheck(string, string, string) (found bool) {
	return false
}
func (allMissingBindings) HasFileInterpreter(string) (found bool) {
	return false
}

// allPresentBindings reports every binding as registered.
type allPresentBindings struct{}

func (allPresentBindings) HasProfileCheck(string, string) (found bool) {
	return true
}
func (allPresentBindings) HasRepositoryScanner(string, string) (found bool) {
	return true
}
func (allPresentBindings) HasTargetCommand(string, string, string) (found bool) {
	return true
}
func (allPresentBindings) HasTargetCheck(string, string, string) (found bool) {
	return true
}
func (allPresentBindings) HasFileInterpreter(string) (found bool) {
	return true
}

func planWithCheck(ruleID string, check style.Job) (plan profile.EffectiveProfile) {
	return profile.EffectiveProfile{
		Effective: style.Plan{
			Rules: []style.Rule{{ID: ruleID, Check: check}},
		},
	}
}

func planWithFix(ruleID string, fix style.Job) (plan profile.EffectiveProfile) {
	return profile.EffectiveProfile{
		Effective: style.Plan{
			Rules: []style.Rule{{
				ID:    ruleID,
				Check: style.ToolchainExecution{ToolIDs: []string{"go"}},
				Fix:   fix,
			}},
		},
	}
}

/* ------------------------------------- Check-side bindings ------------------------------------ */

func TestValidateRuntimeBindingsRejectsMissingProfileCheck(t *testing.T) {
	effective := planWithCheck("project/commands", style.ProfileExecution{
		PackID: "project",
		Check:  "commands",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing profile check to be rejected")
	}

	if !strings.Contains(err.Error(), "profile check") ||
		!strings.Contains(err.Error(), "commands") {
		t.Fatalf("error = %q, want profile check commands", err)
	}
}

func TestValidateRuntimeBindingsRejectsMissingRepositoryScanner(t *testing.T) {
	effective := planWithCheck("text/ascii", style.RepositoryScanExecution{
		PackID:  "text",
		Scanner: "ascii",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing repository scanner to be rejected")
	}

	if !strings.Contains(err.Error(), "repository scanner") {
		t.Fatalf("error = %q, want repository scanner", err)
	}
}

func TestValidateRuntimeBindingsRejectsMissingTargetCommand(t *testing.T) {
	effective := planWithCheck("go/golangci", style.TargetCommandJob{
		PackID:   "go",
		Language: "go",
		Action:   "golangci",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing target command to be rejected")
	}

	if !strings.Contains(err.Error(), "target command") {
		t.Fatalf("error = %q, want target command", err)
	}
}

func TestValidateRuntimeBindingsRejectsMissingTargetCheck(t *testing.T) {
	effective := planWithCheck("go/comments", style.TargetCheckJob{
		PackID:   "go",
		Language: "go",
		Check:    "comments",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing target check to be rejected")
	}

	if !strings.Contains(err.Error(), "target check") {
		t.Fatalf("error = %q, want target check", err)
	}
}

func TestValidateRuntimeBindingsRejectsMissingFileInterpreter(t *testing.T) {
	effective := planWithCheck("text/spelling", style.FileCommandExecution{
		ToolID: "misspell",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing file interpreter to be rejected")
	}

	if !strings.Contains(err.Error(), "file interpreter") {
		t.Fatalf("error = %q, want file interpreter", err)
	}
}

func TestValidateRuntimeBindingsSkipsToolchainJobs(t *testing.T) {
	effective := planWithCheck("toolchain/versions", style.ToolchainExecution{
		PackID:  "project",
		ToolIDs: []string{"go"},
	})

	// Toolchain execution has no Pack-owned runtime binding, so an empty registry is complete.
	if err := profile.ValidateRuntimeBindings(effective, allMissingBindings{}); err != nil {
		t.Fatalf("expected toolchain job to need no binding, got: %v", err)
	}
}

/* -------------------------------------- Fix-side bindings ------------------------------------- */

// File-command fixes run the tool without interpreting output, so they need no file interpreter.
// An absent interpreter must not fail a file-command fix even when nothing is registered.
func TestValidateRuntimeBindingsAcceptsFileCommandFixWithoutInterpreter(t *testing.T) {
	effective := planWithFix("text/spelling", style.FileCommandExecution{
		ToolID: "misspell",
	})

	if err := profile.ValidateRuntimeBindings(effective, allMissingBindings{}); err != nil {
		t.Fatalf("expected file-command fix to need no interpreter, got: %v", err)
	}
}

// Target-command fixes resolve their Pack-qualified command binding like checks.
func TestValidateRuntimeBindingsAcceptsTargetCommandFix(t *testing.T) {
	effective := planWithFix("go/format", style.TargetCommandJob{
		PackID:   "go",
		Language: "go",
		Action:   "go_format",
	})

	if err := profile.ValidateRuntimeBindings(effective, allPresentBindings{}); err != nil {
		t.Fatalf("expected registered target-command fix to validate, got: %v", err)
	}
}

func TestValidateRuntimeBindingsRejectsMissingTargetCommandFix(t *testing.T) {
	effective := planWithFix("go/format", style.TargetCommandJob{
		PackID:   "go",
		Language: "go",
		Action:   "go_format",
	})

	err := profile.ValidateRuntimeBindings(effective, allMissingBindings{})
	if err == nil {
		t.Fatal("expected missing target command fix to be rejected")
	}

	if !strings.Contains(err.Error(), "fix") || !strings.Contains(err.Error(), "target command") {
		t.Fatalf("error = %q, want fix target command", err)
	}
}

// The fix driver set supports file-command and target-command only. A scanner fix is unsupported
// and must be rejected at preparation.
func TestValidateRuntimeBindingsRejectsScannerFix(t *testing.T) {
	effective := planWithFix("text/ascii", style.RepositoryScanExecution{
		PackID:  "text",
		Scanner: "ascii",
	})

	err := profile.ValidateRuntimeBindings(effective, allPresentBindings{})
	if err == nil {
		t.Fatal("expected scanner fix to be rejected as unsupported")
	}

	if !strings.Contains(err.Error(), "fix") || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("error = %q, want unsupported fix", err)
	}
}

func TestValidateRuntimeBindingsRejectsProfileCheckFix(t *testing.T) {
	effective := planWithFix("project/commands", style.ProfileExecution{
		PackID: "project",
		Check:  "commands",
	})

	if err := profile.ValidateRuntimeBindings(effective, allPresentBindings{}); err == nil {
		t.Fatal("expected profile-check fix to be rejected as unsupported")
	}
}

func TestValidateRuntimeBindingsAcceptsCompletePlan(t *testing.T) {
	effective := profile.EffectiveProfile{
		Effective: style.Plan{
			Rules: []style.Rule{
				{
					ID:    "project/commands",
					Check: style.ProfileExecution{PackID: "project", Check: "commands"},
				},
				{
					ID: "go/format",
					Check: style.TargetCommandJob{
						PackID: "go", Language: "go", Action: "go_format",
					},
					Fix: style.FileCommandExecution{ToolID: "goimports"},
				},
			},
		},
	}

	if err := profile.ValidateRuntimeBindings(effective, allPresentBindings{}); err != nil {
		t.Fatalf("expected complete plan to validate, got: %v", err)
	}
}

func TestValidateRuntimeBindingsSkipsNilFix(t *testing.T) {
	effective := profile.EffectiveProfile{
		Effective: style.Plan{
			Rules: []style.Rule{
				{
					ID:    "project/commands",
					Check: style.ProfileExecution{PackID: "project", Check: "commands"},
					Fix:   nil,
				},
			},
		},
	}

	if err := profile.ValidateRuntimeBindings(effective, allPresentBindings{}); err != nil {
		t.Fatalf("expected nil fix to be skipped, got: %v", err)
	}
}
