package validation_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/fixture"
	"ciphera/tools/internal/profile/internal/validation"
)

func TestCheckRejectsNegativeToolExecutionLimits(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		adjust func(*policy.PinnedTool)
		error  string
	}{
		{
			name: "negative timeout",
			adjust: func(tool *policy.PinnedTool) {
				tool.TimeoutSeconds = -1
			},
			error: "timeout_seconds",
		},
		{
			name: "negative output limit",
			adjust: func(tool *policy.PinnedTool) {
				tool.OutputLimitBytes = -1
			},
			error: "output_limit_bytes",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := fixture.Config()

			test.adjust(&config.Tools[0])
			err := validation.Check(config)
			requireErrorContains(t, err, test.error)
		})
	}
}
