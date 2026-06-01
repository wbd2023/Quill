package contract_test

import (
	"testing"

	"ciphera/tools/internal/contract"
)

func TestExecutionSpecUsesTargets(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		spec contract.ExecutionSpec
		want bool
	}{
		{
			name: "target command",
			spec: contract.ExecutionSpec{Detail: contract.TargetCommandExecution{}},
			want: true,
		},
		{
			name: "target check",
			spec: contract.ExecutionSpec{Detail: contract.TargetCheckExecution{}},
			want: true,
		},
		{
			name: "file command",
			spec: contract.ExecutionSpec{Detail: contract.FileCommandExecution{}},
			want: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.spec.UsesTargets(); got != test.want {
				t.Fatalf("UsesTargets() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestExecutionSpecRequiresTargetCheckPaths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		spec contract.ExecutionSpec
		want bool
	}{
		{
			name: "target check",
			spec: contract.ExecutionSpec{Detail: contract.TargetCheckExecution{}},
			want: true,
		},
		{
			name: "target command",
			spec: contract.ExecutionSpec{Detail: contract.TargetCommandExecution{}},
			want: false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.spec.RequiresTargetCheckPaths(); got != test.want {
				t.Fatalf("RequiresTargetCheckPaths() = %t, want %t", got, test.want)
			}
		})
	}
}
