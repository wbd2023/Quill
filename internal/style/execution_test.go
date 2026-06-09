package style_test

import (
	"testing"

	"ciphera/tools/internal/style"
)

func TestExecutionSpecUsesTargets(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		spec style.ExecutionSpec
		want bool
	}{
		{
			name: "target command",
			spec: style.ExecutionSpec{Detail: style.TargetCommandExecution{}},
			want: true,
		},
		{
			name: "target check",
			spec: style.ExecutionSpec{Detail: style.TargetCheckExecution{}},
			want: true,
		},
		{
			name: "file command",
			spec: style.ExecutionSpec{Detail: style.FileCommandExecution{}},
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
		spec style.ExecutionSpec
		want bool
	}{
		{
			name: "target check",
			spec: style.ExecutionSpec{Detail: style.TargetCheckExecution{}},
			want: true,
		},
		{
			name: "target command",
			spec: style.ExecutionSpec{Detail: style.TargetCommandExecution{}},
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
