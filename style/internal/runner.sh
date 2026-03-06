#!/bin/bash
#
# tools/style/internal/runner.sh
# Shared execution helpers for check-style.sh.

set -euo pipefail

# ---------------------------------------------- Paths ---------------------------------------------

STYLE_RUNNER_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/style/internal/common.sh
source "$STYLE_RUNNER_LIB_DIR/common.sh"

# ----------------------------------------- Runtime Context ----------------------------------------

# Provide defaults for sourced-runner globals so shellcheck can reason about this library. The
# entrypoint assigns the real values before invoking the runner helpers.
: "${BLUE:=}" \
	"${COLOUR_RESET:=}" \
	"${CYAN:=}" \
	"${FAILED:=0}" \
	"${GOLANGCI_CACHE:=}" \
	"${GO_BUILD_CACHE:=}" \
	"${GREEN:=}" \
	"${LEVEL_RECOMMENDATION:=recommendation}" \
	"${LEVEL_REQUIRED:=required}" \
	"${MESSAGE_UNKNOWN_RUNNER_KIND:=}" \
	"${MESSAGE_UNKNOWN_RUNNER_TARGET:=}" \
	"${PASSED:=0}" \
	"${PROJECT_ROOT:=}" \
	"${RED:=}" \
	"${SCRIPT_DIR:=}" \
	"${SCOPE:=}" \
	"${SEPARATOR_LINE:=}" \
	"${STRICT_RECOMMENDATIONS:=false}" \
	"${SUMMARY_ALL_PASS:=}" \
	"${SUMMARY_FAILURE:=}" \
	"${SUMMARY_WITH_WARNINGS:=}" \
	"${USAGE_EXIT_CODE:=2}" \
	"${VERBOSE:=false}" \
	"${VERBOSE_OUTPUT_LINE_LIMIT:=30}" \
	"${WARNED:=0}" \
	"${YELLOW:=}"

# --------------------------------------------- Helpers --------------------------------------------

style_runner_prepend_path_if_dir() {
	local candidate_dir="$1"
	if [ -d "$candidate_dir" ]; then
		PATH="$candidate_dir:$PATH"
	fi
}

style_runner_print_tier_heading() {
	local tier_name="$1"

	echo ""
	echo -e "${CYAN}${tier_name}${COLOUR_RESET}"
	echo ""
}

style_runner_scope_allows() {
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

style_runner_run_registered_target() {
	local runner_kind="$1"
	local target="$2"

	case "$runner_kind" in
	"$STYLE_RUNNER_SCRIPT")
		bash "$SCRIPT_DIR/$target"
		;;
	"$STYLE_RUNNER_SCRIPT_SCOPE")
		bash "$SCRIPT_DIR/$target" --scope "$SCOPE"
		;;
	"$STYLE_RUNNER_EXECUTOR")
		case "$target" in
		"$STYLE_RUNNER_TARGET_GOLANGCI_APP")
			(
				cd "$PROJECT_ROOT"
				GOCACHE="$GO_BUILD_CACHE" \
					GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
					golangci-lint run ./...
			)
			;;
		"$STYLE_RUNNER_TARGET_GOLANGCI_TOOLS")
			(
				cd "$PROJECT_ROOT/tools/style/ast"
				GOCACHE="$GO_BUILD_CACHE" \
					GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
					golangci-lint run ./...
			)
			;;
		"$STYLE_RUNNER_TARGET_AST_APP")
			(
				cd "$PROJECT_ROOT/tools/style/ast"
				GOCACHE="$GO_BUILD_CACHE" \
					go run ./cmd/stylecheck \
					"$PROJECT_ROOT/internal" \
					"$PROJECT_ROOT/cmd" \
					"$PROJECT_ROOT/tests"
			)
			;;
		"$STYLE_RUNNER_TARGET_AST_TOOLS")
			(
				cd "$PROJECT_ROOT/tools/style/ast"
				GOCACHE="$GO_BUILD_CACHE" \
					go run ./cmd/stylecheck \
					"$PROJECT_ROOT/tools/style/ast/cmd/stylecheck" \
					"$PROJECT_ROOT/tools/style/ast/internal/checker" \
					"$PROJECT_ROOT/tools/style/tests"
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

style_runner_run_check_with_mode() {
	local mode="$1"
	local rule="$2"
	local name="$3"
	local runner_kind="$4"
	local target="$5"

	printf "  %-8s %-45s" "[$rule]" "$name"

	local output
	local exit_code
	if output=$(style_runner_run_registered_target "$runner_kind" "$target" 2>&1); then
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

style_runner_run_registered_checks() {
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

		if ! style_runner_scope_allows "${CHECK_SCOPES[$index]}"; then
			continue
		fi

		if [ "$level_filter" = "$LEVEL_RECOMMENDATION" ]; then
			style_runner_run_check_with_mode \
				"$LEVEL_RECOMMENDATION" \
				"${CHECK_RULES[$index]}" \
				"${CHECK_NAMES[$index]}" \
				"${CHECK_RUNNERS[$index]}" \
				"${CHECK_TARGETS[$index]}"
			continue
		fi

		style_runner_run_check_with_mode \
			"$LEVEL_REQUIRED" \
			"${CHECK_RULES[$index]}" \
			"${CHECK_NAMES[$index]}" \
			"${CHECK_RUNNERS[$index]}" \
			"${CHECK_TARGETS[$index]}"
	done
}

style_runner_print_summary() {
	echo ""
	echo -e "${BLUE}${SEPARATOR_LINE}${COLOUR_RESET}"
	echo ""
	echo -e "  Results: ${GREEN}$PASSED passed${COLOUR_RESET},"
	echo -e "           ${YELLOW}$WARNED warned${COLOUR_RESET},"
	echo -e "           ${RED}$FAILED failed${COLOUR_RESET}"
	echo ""

	if [ "$FAILED" -eq 0 ]; then
		if [ "$WARNED" -gt 0 ]; then
			echo -e "  ${YELLOW}${SUMMARY_WITH_WARNINGS}${COLOUR_RESET}"
			echo ""
			return 0
		fi

		echo -e "  ${GREEN}${SUMMARY_ALL_PASS}${COLOUR_RESET}"
		echo ""
		return 0
	fi

	echo -e "  ${RED}${SUMMARY_FAILURE}${COLOUR_RESET}"
	echo ""
	return 1
}
