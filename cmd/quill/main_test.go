package main

import (
	"context"
	"os"
	"syscall"
	"testing"
)

type commandRunnerFunc func(context.Context, []string) int

func (run commandRunnerFunc) Run(ctx context.Context, arguments []string) (code int) {
	return run(ctx, arguments)
}

func TestRunCancelsCommandAndPreservesSignalExitStatus(t *testing.T) {
	cases := []struct {
		name   string
		signal os.Signal
		want   int
	}{
		{name: "interrupt", signal: os.Interrupt, want: 130},
		{name: "termination", signal: syscall.SIGTERM, want: 143},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			signals := make(chan os.Signal, 1)
			started := make(chan struct{})
			cancelled := make(chan struct{})
			finished := make(chan int, 1)
			stopped := make(chan struct{})

			go func() {
				finished <- run(
					[]string{"check", "--format", "json"},
					signals,
					func() { close(stopped) },
					commandRunnerFunc(func(ctx context.Context, _ []string) (code int) {
						close(started)
						<-ctx.Done()
						close(cancelled)
						return 1
					}),
				)
			}()

			<-started
			signals <- test.signal

			if got := <-finished; got != test.want {
				t.Fatalf("run signal exit = %d, want %d", got, test.want)
			}
			select {
			case <-cancelled:
			default:
				t.Fatal("command context was not cancelled")
			}
			select {
			case <-stopped:
			default:
				t.Fatal("signal notification was not stopped after cancellation")
			}
		})
	}
}

func TestGetExitCode(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		signal os.Signal
		want   int
	}{
		{name: "interrupt", signal: os.Interrupt, want: 130},
		{name: "termination", signal: syscall.SIGTERM, want: 143},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := getExitCode(test.signal); got != test.want {
				t.Fatalf("getExitCode(%v) = %d, want %d", test.signal, got, test.want)
			}
		})
	}
}
