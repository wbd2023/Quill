package text

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/policy"
	textrules "ciphera/tools/internal/rules/text"
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
	ruleGroupExternal contract.RuleGroup = "external_tools"
	ruleGroupText     contract.RuleGroup = "text_scanners"
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

func rules() (rules []contract.RuleDefinition) {
	return []contract.RuleDefinition{
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
) (rule contract.RuleDefinition) {
	return contract.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupExternal,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutionFileCommand,
			Detail: contract.FileCommandExecution{
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
) (rule contract.RuleDefinition) {
	return scanRule(id, name, ruleGroupText, scanner)
}

func scanRule(
	id string,
	name string,
	group contract.RuleGroup,
	scanner string,
) (rule contract.RuleDefinition) {
	return contract.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: group,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutionRepositoryScan,
			Detail: contract.RepositoryScanExecution{
				Scanner: scanner,
			},
		},
	}
}

func lineLengthRule() (rule contract.RuleDefinition) {
	rule = scannerRule(
		"text/line-length",
		"Line length",
		ScannerLineLength,
	)
	execution := rule.Check.Detail.(contract.RepositoryScanExecution)
	execution.FileSet = "line_length"
	rule.Check.Detail = execution
	return rule
}
