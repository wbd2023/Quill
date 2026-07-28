package process

import (
	"sync/atomic"
	"time"
)

// childWaitDelay bounds the wait for the child to exit after the kill signal so a stubborn process
// cannot hang Run indefinitely. It is shared by the platform-specific process-tree configurators.
const childWaitDelay = 2 * time.Second

// killState records whether the context cancellation actually terminated a running process.
// Termination cause (TimedOut/Canceled) is attributed from this flag rather than from a process
// exit code: a Unix signal death reports a negative exit code, but a Windows kill exits with code
// 1, so the exit code alone is not a reliable cross-platform signal.
type killState struct {
	killed atomic.Bool
}

func (state *killState) markKilled() {
	state.killed.Store(true)
}

func (state *killState) wasKilled() (killed bool) {
	return state.killed.Load()
}
