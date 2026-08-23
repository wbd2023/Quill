package cli

import (
	"context"
	"fmt"
)

// versionCmd prints the build version. It accepts no flags or positional arguments (Kong rejects
// any).
type versionCmd struct{}

func (*versionCmd) run(_ context.Context, runner Runner) (exitCode int) {
	_, _ = fmt.Fprintln(runner.stdout, runner.version)
	return 0
}
