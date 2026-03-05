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
RUNNER_KIND_SCRIPT="script"
RUNNER_KIND_SCRIPT_SCOPE="script_scope"
RUNNER_KIND_EXECUTOR="runner"
RUNNER_TARGET_GOLANGCI_APP="golangci_app"
RUNNER_TARGET_GOLANGCI_TOOLS="golangci_tools"
RUNNER_TARGET_AST_APP="ast_app"
RUNNER_TARGET_AST_TOOLS="ast_tools"
TIER_ONE="tier1"
TIER_TWO="tier2"
TIER_THREE="tier3"
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
# shellcheck source=tools/scripts/lib/style-registry.sh
source "$SCRIPT_DIR/lib/style-registry.sh"

PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

if ! style_require_command "go" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

GO_BIN_DIR="$(go env GOPATH)/bin"
LOCAL_NPM_BIN="${HOME}/.local/bin"

mkdir -p "$GO_BUILD_CACHE" "$GOLANGCI_CACHE"

prepend_path_if_dir() {
	local candidate_dir="$1"
	if [ -d "$candidate_dir" ]; then
		PATH="$candidate_dir:$PATH"
	fi
}

prepend_path_if_dir "$GO_BIN_DIR"
prepend_path_if_dir "$LOCAL_NPM_BIN"

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

run_check_with_mode() {
	local mode="$1"
	local rule="$2"
	local name="$3"
	local runner_kind="$4"
	local target="$5"

	printf "  %-8s %-45s" "[$rule]" "$name"

	local output
	local exit_code
	if output=$(run_registered_target "$runner_kind" "$target" 2>&1); then
		exit_code=0
	else
		exit_code=$?
	fi

	if [ "$exit_code" -eq 0 ]; then
		echo -e "${GREEN}PASS${COLOUR_RESET}"
		PASSED=$((PASSED + 1))
		return
	fi

	if [ "$mode" = "$LEVEL_RECOMMENDATION" ]; then
		if [ "$STRICT_RECOMMENDATIONS" = true ]; then
			echo -e "${RED}FAIL${COLOUR_RESET}"
			FAILED=$((FAILED + 1))
		else
			echo -e "${YELLOW}WARN${COLOUR_RESET}"
			WARNED=$((WARNED + 1))
		fi
	else
		echo -e "${RED}FAIL${COLOUR_RESET}"
		FAILED=$((FAILED + 1))
	fi

	if [ "$VERBOSE" = true ]; then
		echo ""
		echo "$output" | head -"${VERBOSE_OUTPUT_LINE_LIMIT}" | sed 's/^/\t/'
		echo ""
	fi
}

print_tier_heading() {
	local tier_name="$1"

	echo ""
	echo -e "${CYAN}${tier_name}${COLOUR_RESET}"
	echo ""
}

run_required_check() {
	local rule="$1"
	local name="$2"
	local runner_kind="$3"
	local target="$4"

	run_check_with_mode "$LEVEL_REQUIRED" "$rule" "$name" "$runner_kind" "$target"
}

run_recommendation_check() {
	local rule="$1"
	local name="$2"
	local runner_kind="$3"
	local target="$4"

	run_check_with_mode "$LEVEL_RECOMMENDATION" "$rule" "$name" "$runner_kind" "$target"
}

scope_allows() {
	local required_scope="$1"

	case "$required_scope" in
	"$STYLE_SCOPE_ALL")
		return 0
		;;
	"$STYLE_SCOPE_APP")
		[ "$SCOPE" = "$STYLE_SCOPE_APP" ] || [ "$SCOPE" = "$STYLE_SCOPE_ALL" ]
		return
		;;
	"$STYLE_SCOPE_TOOLS")
		[ "$SCOPE" = "$STYLE_SCOPE_TOOLS" ] || [ "$SCOPE" = "$STYLE_SCOPE_ALL" ]
		return
		;;
	*)
		return 1
		;;
	esac
}

run_registered_target() {
	local runner_kind="$1"
	local target="$2"

	case "$runner_kind" in
	"$RUNNER_KIND_SCRIPT")
		bash "$SCRIPT_DIR/$target"
		;;
	"$RUNNER_KIND_SCRIPT_SCOPE")
		bash "$SCRIPT_DIR/$target" --scope "$SCOPE"
		;;
	"$RUNNER_KIND_EXECUTOR")
		case "$target" in
		"$RUNNER_TARGET_GOLANGCI_APP")
			(
				cd "$PROJECT_ROOT"
				GOCACHE="$GO_BUILD_CACHE" \
					GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
					golangci-lint run ./...
			)
			;;
		"$RUNNER_TARGET_GOLANGCI_TOOLS")
			(
				cd "$PROJECT_ROOT/tools/stylecheck"
				GOCACHE="$GO_BUILD_CACHE" \
					GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
					golangci-lint run ./...
			)
			;;
		"$RUNNER_TARGET_AST_APP")
			(
				cd "$PROJECT_ROOT/tools/stylecheck"
				GOCACHE="$GO_BUILD_CACHE" \
					go run . "$PROJECT_ROOT/internal" "$PROJECT_ROOT/cmd" "$PROJECT_ROOT/tests"
			)
			;;
		"$RUNNER_TARGET_AST_TOOLS")
			(
				cd "$PROJECT_ROOT/tools/stylecheck"
				GOCACHE="$GO_BUILD_CACHE" \
					go run . "$PROJECT_ROOT/tools/stylecheck"
			)
			;;
		*)
			echo "$MESSAGE_UNKNOWN_RUNNER_TARGET $target"
			return "$USAGE_EXIT_CODE"
			;;
		esac
		;;
	*)
		echo "$MESSAGE_UNKNOWN_RUNNER_KIND $runner_kind"
		return "$USAGE_EXIT_CODE"
		;;
	esac
}

run_registered_checks() {
	local tier_filter="$1"
	local level_filter="$2"

	local index
	for index in "${!CHECK_TIERS[@]}"; do
		if [ "${CHECK_TIERS[$index]}" != "$tier_filter" ]; then
			continue
		fi

		if [ "${CHECK_LEVELS[$index]}" != "$level_filter" ]; then
			continue
		fi

		if ! scope_allows "${CHECK_SCOPES[$index]}"; then
			continue
		fi

		if [ "$level_filter" = "$LEVEL_RECOMMENDATION" ]; then
			run_recommendation_check \
				"${CHECK_RULES[$index]}" \
				"${CHECK_NAMES[$index]}" \
				"${CHECK_RUNNERS[$index]}" \
				"${CHECK_TARGETS[$index]}"
			continue
		fi

		run_required_check \
			"${CHECK_RULES[$index]}" \
			"${CHECK_NAMES[$index]}" \
			"${CHECK_RUNNERS[$index]}" \
			"${CHECK_TARGETS[$index]}"
	done
}

style_register_default_checks

echo ""
echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
echo -e "${BLUE}${TITLE_TEXT}${COLOUR_RESET}"
echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
echo ""

# ---------------------------------- Tier 1: Go linters (non-1.1) ----------------------------------

print_tier_heading "$TIER_ONE_NAME"

run_registered_checks "$TIER_ONE" "$LEVEL_REQUIRED"

# --------------------------------- Tier 2: Text and script checks ---------------------------------

print_tier_heading "$TIER_TWO_NAME"

run_registered_checks "$TIER_TWO" "$LEVEL_REQUIRED"

if [ "$PROFILE" = "$PROFILE_ALL" ]; then
	run_registered_checks "$TIER_TWO" "$LEVEL_RECOMMENDATION"
fi

# -------------------------------------- Tier 3: AST analysis --------------------------------------

print_tier_heading "$TIER_THREE_NAME"

run_registered_checks "$TIER_THREE" "$LEVEL_REQUIRED"

# --------------------------------------------- Summary --------------------------------------------

echo ""
echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
echo ""
echo -e "  Results: ${GREEN}$PASSED passed${COLOUR_RESET}, ${YELLOW}$WARNED warned${COLOUR_RESET},"
echo -e "           ${RED}$FAILED failed${COLOUR_RESET}"
echo ""

if [ "$FAILED" -eq 0 ]; then
	if [ "$WARNED" -gt 0 ]; then
		echo -e "  ${YELLOW}${SUMMARY_WITH_WARNINGS}${COLOUR_RESET}"
		echo ""
		exit 0
	fi

	echo -e "  ${GREEN}${SUMMARY_ALL_PASS}${COLOUR_RESET}"
	echo ""
	exit 0
fi

echo -e "  ${RED}${SUMMARY_FAILURE}${COLOUR_RESET}"
echo ""
exit 1
