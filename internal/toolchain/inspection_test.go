package toolchain

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

/* ----------------------------------- Inspection Test Doubles ---------------------------------- */

type blockingInspectionRunner struct {
	started chan struct{}
	calls   []string
}

var _ CommandRunner = (*blockingInspectionRunner)(nil)

func (runner *blockingInspectionRunner) ResolvePath(
	ctx context.Context,
	_ map[string]string,
	command string,
) (path string, err error) {
	runner.calls = append(runner.calls, command)
	close(runner.started)
	<-ctx.Done()
	return "", ctx.Err()
}

func (*blockingInspectionRunner) Run(
	context.Context,
	map[string]string,
	string,
	[]string,
) (output string, err error) {
	return "", nil
}

/* -------------------------------- Cancellation Behaviour Tests -------------------------------- */

func TestInspectToolsReturnsParentCancellationAndStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := &blockingInspectionRunner{started: make(chan struct{})}
	result := make(chan error, 1)
	go func() {
		_, err := InspectTools(ctx, runner, map[string]Tool{
			"first": {ID: "first", Command: "first"},
			"later": {ID: "later", Command: "later"},
		}, nil)
		result <- err
	}()

	<-runner.started
	cancel()

	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("InspectTools error = %v, want context.Canceled", err)
	}
	if len(runner.calls) != 1 || runner.calls[0] != "first" {
		t.Fatalf("inspection calls = %v, want [first]", runner.calls)
	}
}

func TestInspectToolsRejectsCancelledContextWithoutTools(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	statuses, err := InspectTools(ctx, &blockingInspectionRunner{}, nil, nil)
	if len(statuses) != 0 {
		t.Fatalf("InspectTools statuses = %v, want none", statuses)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("InspectTools error = %v, want context.Canceled", err)
	}
}

func TestInspectToolsObservesCancellationFromFinalProbe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	statuses, err := InspectTools(ctx, immediateInspectionRunner{}, map[string]Tool{
		"tool": {
			ID:            "tool",
			Command:       "tool",
			PinnedVersion: "1.0.0",
			Version: func(
				context.Context,
				CommandRunner,
				map[string]string,
				string,
			) (string, error) {
				cancel()
				return "1.0.0", nil
			},
		},
	}, nil)
	if len(statuses) != 1 {
		t.Fatalf("InspectTools statuses = %v, want one completed probe", statuses)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("InspectTools error = %v, want context.Canceled", err)
	}
}

/* ------------------------------------ Immediate Test Double ----------------------------------- */

type immediateInspectionRunner struct{}

func (immediateInspectionRunner) ResolvePath(
	context.Context,
	map[string]string,
	string,
) (path string, err error) {
	return "/tool", nil
}

func (immediateInspectionRunner) Run(
	context.Context,
	map[string]string,
	string,
	[]string,
) (output string, err error) {
	return "", nil
}

func TestIsInstalledReturnsParentCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tool")
	if err := os.WriteFile(path, nil, 0o755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	installed, err := IsInstalled(ctx, &blockingInspectionRunner{}, Tool{
		ID:      "tool",
		Command: path,
	}, path)
	if installed {
		t.Fatal("IsInstalled reported an installed tool after cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("IsInstalled error = %v, want context.Canceled", err)
	}
}
