package toolchain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
)

// InspectTools reports the status of each tool in tools, sorted by tool ID.
// Parent-context cancellation aborts inspection and is returned as an operation
// error; ordinary tool probe failures remain individual invalid statuses.
func InspectTools(
	ctx context.Context,
	runner CommandRunner,
	tools map[string]Tool,
	environment map[string]string,
) (statuses []Status, err error) {
	ids := make([]string, 0, len(tools))
	for id := range tools {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	statuses = make([]Status, 0, len(ids))

	for _, id := range ids {
		if err = ctx.Err(); err != nil {
			return statuses, err
		}

		status, inspectErr := inspectTool(ctx, runner, tools[id], environment)
		if inspectErr != nil {
			return statuses, inspectErr
		}
		statuses = append(statuses, status)
	}

	if err = ctx.Err(); err != nil {
		return statuses, err
	}

	return statuses, nil
}

func inspectTool(
	ctx context.Context,
	runner CommandRunner,
	tool Tool,
	environment map[string]string,
) (status Status, inspectionError error) {
	status = Status{Tool: tool}

	path, err := runner.ResolvePath(ctx, environment, tool.Command)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return status, ctxErr
		}
		status.Issue = "missing from PATH"
		return status, nil
	}

	status.Path = path
	version, err := tool.Version(ctx, runner, environment, path)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return status, ctxErr
		}
		status.Issue = err.Error()
		return status, nil
	}

	status.Version = version
	if normaliseVersion(version) != normaliseVersion(tool.PinnedVersion) {
		status.Issue = fmt.Sprintf("requires pinned version %s", tool.PinnedVersion)
		return status, nil
	}

	status.Valid = true
	return status, nil
}

// IsInstalled reports whether a tool matching the pinned version is already installed at the given
// path.
func IsInstalled(
	ctx context.Context,
	runner CommandRunner,
	tool Tool,
	path string,
) (installed bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	probe := tool
	probe.Command = path
	statuses, err := InspectTools(
		ctx,
		runner,
		map[string]Tool{tool.ID: probe},
		nil,
	)
	if err != nil {
		return false, err
	}
	if len(statuses) != 1 {
		return false, fmt.Errorf("inspect local tool %s: missing status", tool.ID)
	}

	return statuses[0].Valid, nil
}
