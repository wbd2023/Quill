package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// repositoryScanDriver returns the repository-scan driver for check execution.
func repositoryScanDriver(scanners RepositoryScanners) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		rule style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		scan, ok := job.(style.RepositoryScan)
		if !ok {
			return style.ExecutionResult{},
				errors.New("repository-scan driver received unsupported execution job")
		}

		scanner, found := scanners.Lookup(rule.PackID, scan.Scanner)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown repository scanner %q for pack %q",
				scan.Scanner,
				rule.PackID,
			)
		}

		return scanner(ctx, context, scan)
	}
}
