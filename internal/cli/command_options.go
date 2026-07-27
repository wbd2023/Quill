package cli

import (
	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/style"
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
	format   report.OutputFormat
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
	format   report.OutputFormat
}

type lockOptions struct {
	repoRoot string
	format   report.OutputFormat
}

type flagHelpError struct {
	message string
}

func (err flagHelpError) Error() (message string) {
	return err.message
}
