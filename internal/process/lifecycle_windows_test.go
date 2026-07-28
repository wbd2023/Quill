//go:build windows

package process

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

// TestCancelProcessReportsDoneWithoutProcess is the Windows guard symmetric to the Unix lifecycle
// test: an absent process yields os.ErrProcessDone and records no kill, so Run does not inject a
// context error after a natural exit.
func TestCancelProcessReportsDoneWithoutProcess(t *testing.T) {
	state := &killState{}
	err := cancelProcess(&exec.Cmd{}, state)()
	if !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("expected os.ErrProcessDone when process is absent, got %v", err)
	}

	if state.wasKilled() {
		t.Fatal("absent process must not record a kill")
	}
}
