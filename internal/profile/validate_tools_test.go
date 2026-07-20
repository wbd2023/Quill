package profile_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/profile/internal/profiletest"
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

			config := profiletest.Config()

			test.adjust(&config.Tools[0])
			err := profile.Validate(config)
			requireErrorContains(t, err, test.error)
		})
	}
}
