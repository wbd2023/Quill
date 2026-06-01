package builtin

import (
	"ciphera/tools/internal/policy"
	textrules "ciphera/tools/internal/rules/text"
)

func textPack() (pack Pack) {
	return Pack{
		ID:   PackText,
		Name: "Text",
		Tools: selectTools(
			ToolMisspell,
		),
		FileSets: textFileSets(),
		Config: PackConfig{
			Required: true,
			Validate: textrules.ValidatePackConfig,
		},
		Rules: []RuleDefinition{
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
				RuleGroupText,
				ScannerSectionHeaders,
			),
			scanRule(
				"text/section-header-density",
				"Section header density",
				RuleGroupText,
				ScannerSectionHeaderDensity,
			),
			scanRule(
				"text/section-header-names",
				"Section header naming",
				RuleGroupText,
				ScannerSectionHeaderNames,
			),
		},
	}
}

func textFileSets() (fileSets policy.FileSets) {
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
