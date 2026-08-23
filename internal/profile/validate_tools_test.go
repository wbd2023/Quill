package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
)

func TestCheckRejectsNegativeToolExecutionLimits(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		adjust func(*profile.PinnedTool)
		error  string
	}{
		{
			name: "negative timeout",
			adjust: func(tool *profile.PinnedTool) {
				tool.TimeoutSeconds = -1
			},
			error: "timeout_seconds",
		},
		{
			name: "negative output limit",
			adjust: func(tool *profile.PinnedTool) {
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
