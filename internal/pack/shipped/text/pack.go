package text

import (
	textrules "ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

const (
	PackID = "text"

	ToolMisspell = "misspell"
)

const (
	ScannerASCII                = "ascii"
	ScannerExceptionMarkers     = "exception_markers"
	ScannerLineLength           = "line_length"
	ScannerMaintenanceMarkers   = "maintenance_markers"
	ScannerSectionHeaderDensity = "section_header_density"
	ScannerSectionHeaderNames   = "section_header_names"
	ScannerSectionHeaders       = "section_headers"
)

const (
	ruleGroupExternal style.RuleGroup = "external_tools"
	ruleGroupText     style.RuleGroup = "text_scanners"
)

// Pack returns the Text Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Text",
		Tools:    append([]toolchain.Capability{}, tools...),
		FileSets: fileSets(),
		Config: pack.Config{
			Required: true,
			Validate: textrules.ValidatePackConfig,
		},
		Rules: rules(),
	}
}

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		fileCommandRule(
			"text/spelling",
			"Spelling (non-Go)",
			ToolMisspell,
			"spelling",
			[]string{"-error", "-locale", "UK"},
		),
		lineLengthRule(),
		scannerRule(
			"text/ascii",
			"ASCII-only characters",
			ScannerASCII,
		),
		scannerRule(
			"text/exception-markers",
			"Exception marker syntax",
			ScannerExceptionMarkers,
		),
		scannerRule(
			"text/maintenance-markers",
			"TODO and FIXME marker format",
			ScannerMaintenanceMarkers,
		),
		scanRule(
			"text/section-headers",
			"Section header format",
			ruleGroupText,
			ScannerSectionHeaders,
		),
		scanRule(
			"text/section-header-density",
			"Section header density",
			ruleGroupText,
			ScannerSectionHeaderDensity,
		),
		scanRule(
			"text/section-header-names",
			"Section header naming",
			ruleGroupText,
			ScannerSectionHeaderNames,
		),
	}
}

func fileSets() (fileSets policy.FileSets) {
	fileSets = append(fileSets, policy.FileSetConfig{
		Name: "line_length",
		Exclude: policy.FileSetExclude{
			Files: []string{"go.sum", "package-lock.json"},
		},
	})
	fileSets = append(fileSets, policy.FileSetConfig{
		Name: "spelling",
		Exclude: policy.FileSetExclude{
			Extensions: []string{".go"},
			Files: []string{
				"COPYING",
				"COPYING.*",
				"LICENSE",
				"LICENSE.*",
				"NOTICE",
				"NOTICE.*",
				"package-lock.json",
			},
		},
	})
	return fileSets
}

/* ---------------------------------------- Rule Builders --------------------------------------- */

func fileCommandRule(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupExternal,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionFileCommand,
			Detail: style.FileCommandExecution{
				ToolID:    toolID,
				FileSet:   fileSet,
				Arguments: append([]string{}, arguments...),
			},
		},
	}
}

func scannerRule(
	id string,
	name string,
	scanner string,
) (rule style.RuleDefinition) {
	return scanRule(id, name, ruleGroupText, scanner)
}

func scanRule(
	id string,
	name string,
	group style.RuleGroup,
	scanner string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: group,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionRepositoryScan,
			Detail: style.RepositoryScanExecution{
				Scanner: scanner,
			},
		},
	}
}

func lineLengthRule() (rule style.RuleDefinition) {
	rule = scannerRule(
		"text/line-length",
		"Line length",
		ScannerLineLength,
	)
	execution := rule.Check.Detail.(style.RepositoryScanExecution)
	execution.FileSet = "line_length"
	rule.Check.Detail = execution
	return rule
}
