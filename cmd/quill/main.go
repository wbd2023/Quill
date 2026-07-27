package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wbd2023/quill/internal/cli"
)

// signalExitBase is the conventional base status for a process terminated by a signal.
const signalExitBase = 128

func main() {
	operationContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	receivedSignal := make(chan os.Signal, 1)
	done := make(chan struct{})
	go func() {
		select {
		case signal := <-signals:
			receivedSignal <- signal
			cancel()
		case <-done:
		}
	}()
	defer close(done)

	tool := cli.New(os.Stdout, os.Stderr, currentVersion())
	exitCode := tool.Run(operationContext, os.Args[1:])

	select {
	case received := <-receivedSignal:
		exitCode = exitCodeForSignal(received)
	default:
	}

	os.Exit(exitCode)
}

func exitCodeForSignal(signal os.Signal) (exitCode int) {
	switch signal {
	case os.Interrupt:
		return signalExitBase + int(syscall.SIGINT)
	case syscall.SIGTERM:
		return signalExitBase + int(syscall.SIGTERM)
	default:
		return 1
	}
}
