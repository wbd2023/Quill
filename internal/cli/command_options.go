package cli

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
)

type checkOptions struct {
	repoRoot              string
	scope                 contract.Scope
	mode                  contract.CheckMode
	format                report.OutputFormat
	strictRecommendations bool
	verbose               bool
}

type fixOptions struct {
	repoRoot string
	scope    contract.Scope
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
