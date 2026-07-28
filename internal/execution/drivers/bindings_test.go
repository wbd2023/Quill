package drivers

import (
	"context"
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------- Pack-local Identity Collisions ------------------------------- */

// Two Packs may each register the same local binding identity without colliding: the Pack-qualified
// key keeps them distinct. Only a true duplicate (same Pack, same local id) is rejected.
func TestBindingsDistinguishSameLocalIDAcrossPacks(t *testing.T) {
	scanner := func(
		_ context.Context,
		_ execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	bindings := NewBindings()
	bindings.AddRepositoryScanner("alpha", "shared", scanner)
	bindings.AddRepositoryScanner("beta", "shared", scanner)

	if _, found := bindings.LookupRepositoryScanner("alpha", "shared"); !found {
		t.Fatal("expected alpha/shared scanner registered")
	}
	if _, found := bindings.LookupRepositoryScanner("beta", "shared"); !found {
		t.Fatal("expected beta/shared scanner registered")
	}
}

func TestProfileChecksRejectSamePackLocalCollision(t *testing.T) {
	check := func(
		_ context.Context,
		_ execution.RunContext,
		_ style.ProfileExecution,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	bindings := NewBindings()
	bindings.AddProfileCheck("project", "commands", check)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate Pack-local profile check to panic")
		}
	}()

	bindings.AddProfileCheck("project", "commands", check)
}

func TestTargetCommandsRejectSamePackLocalCollision(t *testing.T) {
	command := func(
		_ context.Context,
		_ execution.RunContext,
		_ style.TargetCommandJob,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	bindings := NewBindings()
	bindings.AddTargetCommand("go", "go", "golangci", command)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate Pack-local target command to panic")
		}
	}()

	bindings.AddTargetCommand("go", "go", "golangci", command)
}

func TestTargetChecksRejectSamePackLocalCollision(t *testing.T) {
	check := func(
		_ context.Context,
		_ execution.RunContext,
		_ style.TargetCheckJob,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	bindings := NewBindings()
	bindings.AddTargetCheck("go", "go", "comments", check)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate Pack-local target check to panic")
		}
	}()

	bindings.AddTargetCheck("go", "go", "comments", check)
}

// File interpreters stay keyed by global Tool ID, so a duplicate Tool ID always collides.
func TestFileInterpretersRejectDuplicateToolID(t *testing.T) {
	bindings := NewBindings()
	bindings.AddFileInterpreter("misspell", InterpretPlainText(0, "x"))
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate file interpreter to panic")
		}
	}()

	bindings.AddFileInterpreter("misspell", InterpretPlainText(0, "x"))
}

// Satisfies the RuntimeBindings contract so a drivers.Bindings value is usable wherever the
// completeness validation needs a resolver.
func TestBindingsSatisfiesRuntimeBindingsContract(t *testing.T) {
	var _ style.RuntimeBindings = NewBindings()
}
