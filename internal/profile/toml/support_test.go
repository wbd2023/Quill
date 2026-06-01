package toml_test

import (
	"testing"

	"ciphera/tools/internal/policy"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func requireEqual[T any](t *testing.T, name string, expected T, actual T) {
	t.Helper()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("%s mismatch (-expected +actual):\n%s", name, diff)
	}
}

func requireConfigEqual(t *testing.T, expected policy.Config, actual policy.Config) {
	t.Helper()

	if diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty()); diff != "" {
		t.Fatalf("config mismatch (-expected +actual):\n%s", diff)
	}
}
