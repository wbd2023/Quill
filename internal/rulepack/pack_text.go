package rulepack

func textPack() (pack Pack) {
	return Pack{
		ID:   PackText,
		Name: "Text",
		Tools: selectTools(
			ToolMisspell,
		),
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
			scannerRule(
				"text/section-headers",
				"Section header format",
				ScannerSectionHeaders,
			),
			scannerRule(
				"text/section-header-density",
				"Section header density",
				ScannerSectionHeaderDensity,
			),
			scannerRule(
				"text/section-header-names",
				"Section header naming",
				ScannerSectionHeaderNames,
			),
		},
	}
}
