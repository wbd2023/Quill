package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/wbd2023/quill/internal/cli"
)

type commandRunner interface {
	Run(context.Context, []string) (code int)
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	status := run(
		os.Args[1:],
		signals,
		func() { signal.Stop(signals) },
		cli.New(os.Stdout, os.Stderr, version()),
	)
	os.Exit(status)
}

func run(
	arguments []string,
	signals <-chan os.Signal,
	stopSignals func(),
	runner commandRunner,
) (code int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stopped sync.Once
	stop := func() { stopped.Do(stopSignals) }
	defer stop()

	cause := make(chan os.Signal, 1)
	done := make(chan struct{})

	go func() {
		select {
		case received := <-signals:
			cause <- received
			cancel()
			stop()

		case <-done:
		}
	}()
	defer close(done)

	code = runner.Run(ctx, arguments)

	select {
	case received := <-cause:
		return getExitCode(received)
	default:
		return code
	}
}

func getExitCode(signal os.Signal) (code int) {
	// Shells conventionally report signal termination as 128 plus the signal number.
	const offset = 128

	switch signal {
	case os.Interrupt:
		return offset + int(syscall.SIGINT)
	case syscall.SIGTERM:
		return offset + int(syscall.SIGTERM)
	default:
		return 1
	}
}
