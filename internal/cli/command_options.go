package cli

import (
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/style"
)

type checkOptions struct {
	repoRoot              string
	scope                 style.Scope
	mode                  style.CheckMode
	format                report.OutputFormat
	strictRecommendations bool
	verbose               bool
}

type fixOptions struct {
	repoRoot string
	scope    style.Scope
}

type doctorOptions struct {
	repoRoot string
	format   report.OutputFormat
}

type coverageOptions struct {
	repoRoot string
	format   report.OutputFormat
	verbose  bool
}

type installOptions struct {
	repoRoot string
}

type flagHelpError struct {
	message string
}

func (err flagHelpError) Error() (message string) {
	return err.message
}
