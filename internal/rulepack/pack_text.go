package rulepack

import (
	"ciphera/tools/internal/contract"
)

/* ------------------------------------------ Text Pack ----------------------------------------- */

func textPack() (pack Pack) {
	return Pack{
		ID:   PackText,
		Name: "Text",
		Tools: selectTools(
			contract.ToolMisspell,
		),
		Rules: []RuleDefinition{
			fileCommandRule(
				"docs/spelling",
				"Spelling (non-Go)",
				contract.ToolMisspell,
				"spelling",
				[]string{"-error", "-locale", "UK"},
			),
			lineLengthRule(),
			repoScanRule(
				"repo/ascii",
				"ASCII-only characters",
				RepositoryScannerASCII,
			),
			repoScanRule(
				"repo/exception-markers",
				"Exception marker syntax",
				RepositoryScannerExceptionMarkers,
			),
			repoScanRule(
				"repo/secrets",
				"Committed secrets",
				RepositoryScannerSecrets,
			),
			repoScanRule(
				"repo/maintenance-markers",
				"TODO and FIXME marker format",
				RepositoryScannerMaintenanceMarkers,
			),
			repoScanRule(
				"repo/section-headers",
				"Section header format",
				RepositoryScannerSectionHeaders,
			),
			repoScanRule(
				"repo/section-header-names",
				"Section header naming",
				RepositoryScannerSectionHeaderNames,
			),
		},
	}
}
