package drivers

import (
	"context"
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/pack/external"
	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// externalOutputLimit bounds each captured stream (standard output and standard error) from an
// external Pack subprocess. Each stream is bounded independently; overflow is reported as an
// execution error rather than silently truncated findings.
const externalOutputLimit int64 = 1 << 20

// stderrExcerpt is the largest tail of standard error folded into an external execution error so a
// failing Pack leaves a debugging breadcrumb without flooding the report.
const stderrExcerpt = 512

/* --------------------------------------- External Driver -------------------------------------- */

// externalCheckDriver returns the flat driver for external Pack checks. Every external check is
// self-describing - the bound job carries its executable command, Pack directory, and runtime
// limits - so one driver serves all external Packs without a Pack-qualified binding registry. The
// driver resolves the executable beneath the Pack directory, writes one JSON request to standard
// input, and reads JSON Lines diagnostics followed by exactly one terminal completion from
// standard output. Standard error is captured separately and surfaced in error context only.
func externalCheckDriver() (driver execution.Driver) {
	return func(
		ctx context.Context,
		run execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		externalJob, found := job.(style.ExternalCheckJob)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"external check driver received an empty job",
			)
		}

		files, err := collectExternalFiles(run, externalJob.FileSet)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		executable, err := external.ResolveExecutable(
			externalJob.PackDirectory, externalJob.Command,
		)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		configuration, found := run.Profile.PackConfigs.Lookup(externalJob.PackID)
		if !found || configuration == nil {
			configuration = map[string]any{}
		}

		payload, err := external.EncodeRequest(external.Request{
			Protocol:       external.ProtocolVersion,
			Operation:      "check",
			RepositoryRoot: run.RepoRoot,
			PackID:         externalJob.PackID,
			RuleID:         externalJob.RuleID,
			CheckID:        externalJob.CheckID,
			Scope:          string(run.Scope),
			Files:          files,
			Configuration:  configuration,
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}

		timeout := externalJob.Timeout
		if timeout <= 0 {
			timeout = external.DefaultTimeout
		}

		commandResult, runErr := process.RunCommand(ctx, process.CommandRequest{
			Name:             executable,
			Environment:      process.EnvironmentIsolated,
			Directory:        run.RepoRoot,
			Stdin:            payload,
			Timeout:          timeout,
			OutputLimitBytes: externalOutputLimit,
		})

		result = style.ExecutionResult{
			PackID:    externalJob.PackID,
			ExitCode:  commandResult.ExitCode,
			TimedOut:  commandResult.TimedOut,
			Truncated: commandResult.Truncated,
		}

		if runErr != nil {
			return result, externalRunError(externalJob, commandResult, runErr)
		}

		if commandResult.Truncated {
			return result, fmt.Errorf(
				"external pack %q rule %q exceeded the output limit%s",
				externalJob.PackID, externalJob.RuleID, stderrContext(commandResult.Stderr),
			)
		}

		outcome, parseErr := external.ParseResponse(commandResult.Stdout)
		if parseErr != nil {
			return result, fmt.Errorf(
				"external pack %q rule %q: %w%s",
				externalJob.PackID, externalJob.RuleID, parseErr,
				stderrContext(commandResult.Stderr),
			)
		}

		if !outcome.Success {
			return result, fmt.Errorf(
				"external pack %q rule %q reported failure: %s%s",
				externalJob.PackID, externalJob.RuleID, outcome.Error,
				stderrContext(commandResult.Stderr),
			)
		}

		result.Diagnostics = outcome.Diagnostics
		return result, nil
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func collectExternalFiles(run execution.RunContext, fileSet string) (files []string, err error) {
	var candidates []string
	if fileSet == "" {
		candidates, err = execution.CollectScopeFiles(run)
	} else {
		candidates, err = execution.CollectFileSetFiles(run, fileSet)
	}
	if err != nil {
		return nil, err
	}

	files = make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		relative, relErr := filewalk.RelativePath(run.RepoRoot, candidate)
		if relErr != nil {
			return nil, relErr
		}
		files = append(files, relative)
	}
	return files, nil
}

func externalRunError(
	job style.ExternalCheckJob,
	result process.CommandResult,
	runErr error,
) (err error) {
	switch {
	case result.TimedOut:
		return fmt.Errorf(
			"external pack %q rule %q timed out%s",
			job.PackID, job.RuleID, stderrContext(result.Stderr),
		)

	case result.Canceled:
		return fmt.Errorf(
			"external pack %q rule %q was canceled",
			job.PackID, job.RuleID,
		)

	default:
		return fmt.Errorf(
			"external pack %q rule %q exited with status %d: %v%s",
			job.PackID, job.RuleID, result.ExitCode, runErr, stderrContext(result.Stderr),
		)
	}
}

// stderrContext folds a bounded tail of standard error into an error message when present, so a
// failing Pack leaves a debugging breadcrumb. Standard error never carries diagnostics.
func stderrContext(stderr string) (suffix string) {
	if strings.TrimSpace(stderr) == "" {
		return ""
	}

	excerpt := strings.TrimSpace(stderr)
	if len(excerpt) > stderrExcerpt {
		excerpt = excerpt[len(excerpt)-stderrExcerpt:]
	}
	return "\n" + excerpt
}
