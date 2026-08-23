package drivers

import (
	"context"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
)

/* ---------------------------- Runtime Binding Completeness Contract --------------------------- */

// noOpBindings return closures matching each typed binding signature without doing any work; they
// exist only so a binding can be registered for the "present" cases.
func noOpScanner(
	context.Context,
	execution.RunContext,
	style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return style.ExecutionResult{}, nil
}

func noOpProfileCheck(
	context.Context,
	execution.RunContext,
	style.ProfileCheck,
) (result style.ExecutionResult, err error) {
	return style.ExecutionResult{}, nil
}

func noOpTargetCommand(
	context.Context,
	execution.RunContext,
	style.TargetCommandJob,
) (result style.ExecutionResult, err error) {
	return style.ExecutionResult{}, nil
}

func noOpTargetCheck(
	context.Context,
	execution.RunContext,
	style.TargetCheckJob,
) (result style.ExecutionResult, err error) {
	return style.ExecutionResult{}, nil
}

func validateOne(rule style.Rule, bindings Bindings) (err error) {
	return bindings.Validate(style.Plan{Rules: []style.Rule{rule}})
}

// TestValidateAcceptsNoOpFamiliesWithoutBindings proves toolchain and external checks need no
// Pack-owned runtime binding: they are intrinsic or self-describing.
func TestValidateAcceptsNoOpFamiliesWithoutBindings(t *testing.T) {
	t.Parallel()

	bindings := NewBindings()
	for _, job := range []style.Job{
		style.ToolchainCheck{ToolIDs: []string{"go"}},
		style.ExternalCheck{CheckID: "ext"},
	} {
		rule := style.Rule{ID: "r", PackID: "p", Check: job}
		if err := validateOne(rule, bindings); err != nil {
			t.Fatalf("%T check rejected without a binding: %v", job, err)
		}
	}
}

// TestValidateRejectsCheckJobsMissingBindings proves every check Job that owns a runtime identity
// resolves to exactly one registered binding, and that registering it clears the error.
func TestValidateRejectsCheckJobsMissingBindings(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		job      style.Job
		register func(bindings *Bindings)
	}{
		{
			name: "profile check",
			job:  style.ProfileCheck{Check: "commands"},
			register: func(bindings *Bindings) {
				bindings.AddProfileCheck("p", "commands", noOpProfileCheck)
			},
		},
		{
			name: "repository scanner",
			job:  style.RepositoryScan{Scanner: "secrets"},
			register: func(bindings *Bindings) {
				bindings.AddRepositoryScanner("p", "secrets", noOpScanner)
			},
		},
		{
			name: "target command",
			job:  style.TargetCommandJob{Language: "go", Action: "golangci"},
			register: func(bindings *Bindings) {
				bindings.AddTargetCommand("p", "go", "golangci", noOpTargetCommand)
			},
		},
		{
			name: "target check",
			job:  style.TargetCheckJob{Language: "go", Check: "comments"},
			register: func(bindings *Bindings) {
				bindings.AddTargetCheck("p", "go", "comments", noOpTargetCheck)
			},
		},
		{
			name: "file command interpreter",
			job:  style.FileCommand{ToolID: "markdownlint", FileSet: "markdown"},
			register: func(bindings *Bindings) {
				bindings.AddFileInterpreter("markdownlint", InterpretPlainText(ExitFindings, "x"))
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rule := style.Rule{ID: "r", PackID: "p", Check: test.job}

			if err := validateOne(rule, NewBindings()); err == nil {
				t.Fatalf("expected missing-binding error, got nil")
			}

			bindings := NewBindings()
			test.register(&bindings)
			if err := validateOne(rule, bindings); err != nil {
				t.Fatalf("expected registered binding to validate, got %v", err)
			}
		})
	}
}

// TestValidateReportsPackProvenanceInMissingBindingMessage proves the error names the Pack the
// binding identity is scoped to, since binding keys are Pack-qualified.
func TestValidateReportsPackProvenanceInMissingBindingMessage(t *testing.T) {
	t.Parallel()

	rule := style.Rule{
		ID: "go/comments", PackID: "go",
		Check: style.TargetCheckJob{Language: "go", Check: "comments"},
	}
	err := validateOne(rule, NewBindings())
	if err == nil {
		t.Fatal("expected missing-binding error")
	}
	if !strings.Contains(err.Error(), "pack \"go\"") {
		t.Fatalf("error %q does not name the pack", err)
	}
	if !strings.Contains(err.Error(), "go/comments") {
		t.Fatalf("error %q does not name the local identity", err)
	}
}

// TestValidateFixAsymmetry proves fixes support file-command and target-command only: a
// file-command fix needs no interpreter, a target-command fix still resolves its binding, and every
// other family is rejected as an unsupported fix at preparation.
func TestValidateFixAsymmetry(t *testing.T) {
	t.Parallel()

	t.Run("file command fix needs no interpreter", func(t *testing.T) {
		t.Parallel()
		rule := style.Rule{
			ID: "r", PackID: "p",
			Fix: style.FileCommand{ToolID: "markdownlint", FileSet: "markdown"},
		}
		if err := validateOne(rule, NewBindings()); err != nil {
			t.Fatalf("file-command fix should not require an interpreter: %v", err)
		}
	})

	t.Run("target command fix resolves its binding", func(t *testing.T) {
		t.Parallel()
		rule := style.Rule{
			ID: "r", PackID: "p",
			Fix: style.TargetCommandJob{Language: "go", Action: "go_format"},
		}
		if err := validateOne(rule, NewBindings()); err == nil {
			t.Fatal("expected target-command fix to require a binding")
		}
		bindings := NewBindings()
		bindings.AddTargetCommand("p", "go", "go_format", noOpTargetCommand)
		if err := validateOne(rule, bindings); err != nil {
			t.Fatalf("registered target-command fix should validate: %v", err)
		}
	})

	unsupported := []struct {
		name string
		job  style.Job
	}{
		{"toolchain", style.ToolchainCheck{ToolIDs: []string{"go"}}},
		{"profile", style.ProfileCheck{Check: "commands"}},
		{"repository scan", style.RepositoryScan{Scanner: "secrets"}},
		{"target check", style.TargetCheckJob{Language: "go", Check: "comments"}},
		{"external", style.ExternalCheck{CheckID: "ext"}},
	}
	for _, test := range unsupported {
		t.Run(test.name+" fix is unsupported", func(t *testing.T) {
			t.Parallel()
			rule := style.Rule{ID: "r", PackID: "p", Fix: test.job}
			err := validateOne(rule, NewBindings())
			if err == nil {
				t.Fatalf("expected %s fix to be rejected as unsupported", test.name)
			}
			if !strings.Contains(
				err.Error(), "fixes support file-command and target-command only",
			) {
				t.Fatalf("error %q does not name the fix constraint", err)
			}
		})
	}
}
