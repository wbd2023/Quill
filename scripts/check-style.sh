#!/bin/bash
#
# tools/scripts/check-style.sh
# Master STYLE.md compliance runner.
#
# Usage:
#	./tools/scripts/check-style.sh [--verbose] [--scope app|tools|all] \
#		[--profile required|all] [--strict-recommendations]
#
# Exit code: 0 if all checks pass, 1 if any fail.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
VERBOSE_OUTPUT_LINE_LIMIT=30
SEPARATOR_LINE="-------------------------------------------------------------------------------"
TITLE_TEXT="                           STYLE.md Compliance Check"
PROFILE_REQUIRED="required"
PROFILE_ALL="all"
LEVEL_REQUIRED="$PROFILE_REQUIRED"
LEVEL_RECOMMENDATION="recommendation"
TIER_ONE_NAME="Tier 1: Go linters"
TIER_TWO_NAME="Tier 2: Text and script checks"
TIER_THREE_NAME="Tier 3: AST analysis"
USAGE_LINE_ONE_SUFFIX=""
USAGE_LINE_TWO='       [--profile required|all] [--strict-recommendations]'
PROFILE_ERROR_PREFIX="invalid profile:"
PROFILE_ERROR_SUFFIX="(expected required|all)"
SUMMARY_WITH_WARNINGS="Required checks passed with recommendations."
SUMMARY_ALL_PASS="All STYLE.md checks passed."
SUMMARY_FAILURE="Some checks failed. Run with --verbose for details."
MESSAGE_UNKNOWN_RUNNER_TARGET="unknown runner target:"
MESSAGE_UNKNOWN_RUNNER_KIND="unknown runner kind:"

CACHE_ROOT="${TMPDIR:-/tmp}/ciphera-stylecheck-${USER:-user}"
GO_BUILD_CACHE="$CACHE_ROOT/go-build"
GOLANGCI_CACHE="$CACHE_ROOT/golangci-lint"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[0;33m'
COLOUR_RESET='\033[0m'

VERBOSE=false
SCOPE=""
PROFILE="$PROFILE_REQUIRED"
STRICT_RECOMMENDATIONS=false
FAILED=0
PASSED=0
WARNED=0

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"
# shellcheck source=tools/scripts/lib/style-runner.sh
source "$SCRIPT_DIR/lib/style-runner.sh"
# shellcheck source=tools/scripts/lib/style-registry.sh
source "$SCRIPT_DIR/lib/style-registry.sh"

PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

if ! style_require_command "go" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

GO_BIN_DIR="$(go env GOPATH)/bin"
LOCAL_NPM_BIN="${HOME}/.local/bin"

mkdir -p "$GO_BUILD_CACHE" "$GOLANGCI_CACHE"

style_runner_prepend_path_if_dir "$GO_BIN_DIR"
style_runner_prepend_path_if_dir "$LOCAL_NPM_BIN"

SCOPE="$STYLE_SCOPE_ALL"
USAGE_LINE_ONE_SUFFIX="[--verbose] [--scope $STYLE_SCOPE_USAGE]"

# ---------------------------------------------- Args ----------------------------------------------

print_usage() {
	printf 'usage: %s %s\n' "$0" "$USAGE_LINE_ONE_SUFFIX" >&2
	echo "$USAGE_LINE_TWO" >&2
}

while [ $# -gt 0 ]; do
	case "$1" in
	--verbose)
		VERBOSE=true
		shift
		;;
	--scope)
		if [ "$#" -lt 2 ]; then
			print_usage
			exit "$USAGE_EXIT_CODE"
		fi
		SCOPE="$2"
		shift 2
		;;
	--profile)
		if [ "$#" -lt 2 ]; then
			print_usage
			exit "$USAGE_EXIT_CODE"
		fi
		PROFILE="$2"
		shift 2
		;;
	--strict-recommendations)
		STRICT_RECOMMENDATIONS=true
		shift
		;;
	*)
		print_usage
		exit "$USAGE_EXIT_CODE"
		;;
	esac
done

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if [[ ! "$PROFILE" =~ ^($PROFILE_REQUIRED|$PROFILE_ALL)$ ]]; then
	echo "$PROFILE_ERROR_PREFIX $PROFILE $PROFILE_ERROR_SUFFIX"
	exit "$USAGE_EXIT_CODE"
fi

style_register_default_checks

echo ""
echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
echo -e "${BLUE}${TITLE_TEXT}${COLOUR_RESET}"
echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
echo ""

# ---------------------------------- Tier 1: Go linters (non-1.1) ----------------------------------

style_runner_print_tier_heading "$TIER_ONE_NAME"

style_runner_run_registered_checks "$STYLE_TIER_ONE" "$LEVEL_REQUIRED"

# --------------------------------- Tier 2: Text and script checks ---------------------------------

style_runner_print_tier_heading "$TIER_TWO_NAME"

style_runner_run_registered_checks "$STYLE_TIER_TWO" "$LEVEL_REQUIRED"

if [ "$PROFILE" = "$PROFILE_ALL" ]; then
	style_runner_run_registered_checks "$STYLE_TIER_TWO" "$LEVEL_RECOMMENDATION"
fi

# -------------------------------------- Tier 3: AST analysis --------------------------------------

style_runner_print_tier_heading "$TIER_THREE_NAME"

style_runner_run_registered_checks "$STYLE_TIER_THREE" "$LEVEL_REQUIRED"

# --------------------------------------------- Summary --------------------------------------------

if style_runner_print_summary; then
	exit 0
fi

exit 1
