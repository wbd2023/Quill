package rulepack

/* ------------------------------------ Control-Plane Checks ------------------------------------ */

const (
	ControlPlaneCheckEnforcementLevels = "enforcement_levels"
	ControlPlaneCheckGlobalExclusions  = "global_exclusions"
	ControlPlaneCheckQualityTargets    = "quality_targets"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

const (
	RepositoryScannerArchitecture       = "architecture"
	RepositoryScannerASCII              = "ascii"
	RepositoryScannerBashMagicValues    = "bash_magic_values"
	RepositoryScannerBashSafety         = "bash_safety"
	RepositoryScannerBashStructure      = "bash_structure"
	RepositoryScannerBashTestHygiene    = "bash_test_hygiene"
	RepositoryScannerExceptionMarkers   = "exception_markers"
	RepositoryScannerGuardClauseSpacing = "guard_clause_spacing"
	RepositoryScannerLineLength         = "line_length"
	RepositoryScannerMaintenanceMarkers = "maintenance_markers"
	RepositoryScannerNaming             = "naming"
	RepositoryScannerSecrets            = "secrets"
	RepositoryScannerSectionHeaderNames = "section_header_names"
	RepositoryScannerSectionHeaders     = "section_headers"
	RepositoryScannerSwitchCaseSpacing  = "switch_case_spacing"
)
