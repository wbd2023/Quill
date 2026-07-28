package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/toolchain"
)

// recordingRunner satisfies toolchain.CommandRunner and counts every process resolution or
// execution it is asked to perform. Coverage must never reach tool inspection, so a successful
// Coverage operation leaves the count at zero.
type recordingRunner struct {
	calls int
}

var _ toolchain.CommandRunner = (*recordingRunner)(nil)

func (runner *recordingRunner) ResolvePath(
	_ context.Context,
	_ map[string]string,
	_ string,
) (path string, err error) {
	runner.calls++
	return "", nil
}

func (runner *recordingRunner) Run(
	_ context.Context,
	_ map[string]string,
	_ string,
	_ []string,
) (output string, err error) {
	runner.calls++
	return "", nil
}

// TestCoverageInspectsNoTools asserts the observable side effect of Coverage being metadata-only:
// it shares the prepared document and never spawns a process or inspects a tool. The check is on
// actual command execution, not on an internal construction seam.
func TestCoverageInspectsNoTools(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{}
	engine, err := New(testutil.RepositoryRoot(t), WithCommandRunner(runner))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if _, err = engine.Coverage(context.Background()); err != nil {
		t.Fatalf("Coverage: %v", err)
	}

	if runner.calls != 0 {
		t.Fatalf("Coverage invoked the command runner %d time(s), want 0", runner.calls)
	}
}

func TestCoverageRejectsCancelledContext(t *testing.T) {
	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	operationContext, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err = engine.Coverage(operationContext); !errors.Is(err, context.Canceled) {
		t.Fatalf("Coverage error = %v, want context.Canceled", err)
	}
}
