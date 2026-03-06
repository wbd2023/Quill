package scripts_test

import "strings"

const (
	styleRegistryFieldCount = 7

	styleRegistryTableEnvName = "STYLE_REGISTRY_TABLE_FILE"
	styleTestLogEnvName       = "STYLE_TEST_LOG"

	registryTierOne   = "tier1"
	registryTierTwo   = "tier2"
	registryTierThree = "tier3"

	registryLevelRequired       = "required"
	registryLevelRecommendation = "recommendation"
	styleProfileAll             = "all"

	registryScopeApp   = "app"
	registryScopeTools = "tools"
	registryScopeAll   = "all"

	registryRunnerScript      = "script"
	registryRunnerScriptScope = "script_scope"
	registryRunnerExecutor    = "runner"

	registryTargetGolangciApp   = "golangci_app"
	registryTargetGolangciTools = "golangci_tools"
	registryTargetASTApp        = "ast_app"
	registryTargetASTTools      = "ast_tools"

	checkPathBashLineLength       = "checks/bash/line-length.sh"
	checkPathBashMagicValues      = "checks/bash/magic-values.sh"
	checkPathBashShellcheck       = "checks/bash/shellcheck.sh"
	checkPathBashShfmt            = "checks/bash/shfmt.sh"
	checkPathBashStyle            = "checks/bash/style.sh"
	checkPathGeneralASCII         = "checks/general/ascii.sh"
	checkPathGeneralMarkdown      = "checks/general/markdown.sh"
	checkPathGeneralNaming        = "checks/general/naming.sh"
	checkPathGeneralHeaders       = "checks/general/section-headers.sh"
	checkPathGeneralSpelling      = "checks/general/spelling.sh"
	checkPathGoArchitecture       = "checks/go/architecture-imports.sh"
	checkPathGoLineLength         = "checks/go/line-length.sh"
	checkPathGoMagicValues        = "checks/go/magic-values.sh"
	checkPathGoVerticalSpacing    = "checks/go/vertical-spacing.sh"
	missingCheckScriptPath        = "checks/test/missing.sh"
	requiredCheckScriptPath       = "checks/test/required.sh"
	appOnlyCheckScriptPath        = "checks/test/app-only.sh"
	recommendationCheckScriptPath = "checks/test/recommendation.sh"
)

type styleRegistryRow struct {
	Tier   string
	Level  string
	Rule   string
	Name   string
	Scope  string
	Runner string
	Target string
}

/* ----------------------------------------- Row Helpers ---------------------------------------- */

func newStyleRegistryRow(
	tier string,
	level string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (row styleRegistryRow) {
	return styleRegistryRow{
		Tier:   tier,
		Level:  level,
		Rule:   rule,
		Name:   name,
		Scope:  scope,
		Runner: runner,
		Target: target,
	}
}

func registryTableRow(
	tier string,
	level string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	return strings.Join([]string{tier, level, rule, name, scope, runner, target}, "|")
}

func requiredRegistryRow(
	tier string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	return registryTableRow(tier, registryLevelRequired, rule, name, scope, runner, target)
}

func recommendationRegistryRow(
	tier string,
	rule string,
	name string,
	scope string,
	runner string,
	target string,
) (registryRowText string) {
	return registryTableRow(
		tier,
		registryLevelRecommendation,
		rule,
		name,
		scope,
		runner,
		target,
	)
}
