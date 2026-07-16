package profile

import (
	"strings"
	"testing"
)

func requireErrorContainsInternal(tb testing.TB, err error, text string) {
	tb.Helper()

	if err == nil {
		tb.Fatalf("expected error containing %q, got nil", text)
	}

	if !strings.Contains(err.Error(), text) {
		tb.Fatalf("expected error containing %q, got %v", text, err)
	}
}
