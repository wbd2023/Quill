#!/bin/bash
#
# tools/style/checks/bash/shellcheck.sh
# Checks Bash scripts with shellcheck (STYLE.md 2.11).
#
# Rules:
#	- All Bash scripts must pass shellcheck static analysis.
#
# Usage:
#	./tools/style/checks/bash/shellcheck.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found, 2 if tool unavailable.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
RULE_LABEL="2.11"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] shellcheck findings:"
SHELLCHECK_MODE="-x"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STYLE_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
# shellcheck source=tools/style/internal/common.sh
source "$STYLE_DIR/internal/common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$STYLE_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

mapfile -t SCRIPT_FILES < <(style_collect_files "$PROJECT_ROOT" "$SCOPE" "$STYLE_PATTERN_SHELL")

if [ "${#SCRIPT_FILES[@]}" -eq 0 ]; then
	exit 0
fi

if ! style_require_command "shellcheck" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! output=$(shellcheck "$SHELLCHECK_MODE" "${SCRIPT_FILES[@]}" 2>&1); then
	echo "$RULE_MESSAGE_PREFIX"
	echo "$output"
	echo ""
	exit 1
fi

exit 0
