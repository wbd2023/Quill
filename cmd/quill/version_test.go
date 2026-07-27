package main

import (
	"os"
	"runtime/debug"
	"syscall"
	"testing"
)

func TestVersionFromBuildInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		buildInfo *debug.BuildInfo
		ok        bool
		want      string
	}{
		{
			name:      "tagged build",
			buildInfo: &debug.BuildInfo{Main: debug.Module{Version: "v0.1.0"}},
			ok:        true,
			want:      "v0.1.0",
		},
		{
			name:      "development build",
			buildInfo: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			ok:        true,
			want:      developmentVersion,
		},
		{
			name:      "missing build information",
			buildInfo: &debug.BuildInfo{},
			ok:        false,
			want:      developmentVersion,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := versionFromBuildInfo(testCase.buildInfo, testCase.ok)
			if got != testCase.want {
				t.Fatalf("versionFromBuildInfo() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestExitCodeForSignal(t *testing.T) {
	testCases := []struct {
		name   string
		signal os.Signal
		want   int
	}{
		{name: "interrupt", signal: os.Interrupt, want: 130},
		{name: "termination", signal: syscall.SIGTERM, want: 143},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := exitCodeForSignal(testCase.signal); got != testCase.want {
				t.Fatalf("exitCodeForSignal(%v) = %d, want %d", testCase.signal, got, testCase.want)
			}
		})
	}
}
