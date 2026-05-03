package policy_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func requireEqual[T any](t *testing.T, want T, got T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected value (-expected +actual):\n%s", diff)
	}
}
