#!/bin/bash
#
# tools/style/checks/repository/naming.sh
# Checks naming conventions using text patterns (STYLE.md 2.2).
#
# Rules:
#	- Use "Repository" not "Store" in type names.
#	- Use "xRepository" not "xRepo" abbreviations.
#	- Use descriptive constant names in Bash scripts (for example COLOUR_RESET over NC).
#
# Usage:
#	./tools/style/checks/repository/naming.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="2.2"
STORE_TYPE_REGEX='type\s+\w*Store\s+'
REPOSITORY_ABBREV_REGEX='\b\w+Repo\b'
REPOSITORY_NAME_FILTER='Repository'
COMMENT_LINE_FILTER='^[^:]+:[0-9]+:\s*//'
NC_CONSTANT_REGEX='(^|[[:space:]])(local[[:space:]]+)?NC='
MESSAGE_STORE_TYPES="[${RULE_LABEL}] Use 'Repository' not 'Store' in type names:"
MESSAGE_REPO_ABBREV="[${RULE_LABEL}] Use 'xRepository' not 'xRepo':"
MESSAGE_NC_NAME="[${RULE_LABEL}] Use descriptive constant names in Bash scripts"
MESSAGE_NC_NAME+=" (prefer COLOUR_RESET over NC):"

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

if [ "$SCOPE" = "$STYLE_SCOPE_APP" ] || [ "$SCOPE" = "$STYLE_SCOPE_ALL" ]; then
	# "Store" in type names (should be "Repository").
	store_types=$(rg -n --glob "$STYLE_PATTERN_GO" --glob '!*_test.go' "$STORE_TYPE_REGEX" \
		"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" |
		grep -v "$STYLE_PATH_VENDOR" || true)

	if [ -n "$store_types" ]; then
		echo "$MESSAGE_STORE_TYPES"
		echo "$store_types"
		echo ""
		FOUND=1
	fi

	# "xRepo" abbreviation (should be "xRepository").
	repo_abbrev=$(rg -n --glob "$STYLE_PATTERN_GO" "$REPOSITORY_ABBREV_REGEX" \
		"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" |
		grep -v "$STYLE_PATH_VENDOR" |
		grep -v "$REPOSITORY_NAME_FILTER" |
		grep -Ev "$COMMENT_LINE_FILTER" || true)

	if [ -n "$repo_abbrev" ]; then
		echo "$MESSAGE_REPO_ABBREV"
		echo "$repo_abbrev"
		echo ""
		FOUND=1
	fi
fi

if [ "$SCOPE" = "$STYLE_SCOPE_TOOLS" ] || [ "$SCOPE" = "$STYLE_SCOPE_ALL" ]; then
	non_descriptive_nc=$(rg -n --glob "$STYLE_PATTERN_SHELL" "$NC_CONSTANT_REGEX" \
		"$PROJECT_ROOT/$STYLE_PATH_TOOLS" || true)

	if [ -n "$non_descriptive_nc" ]; then
		echo "$MESSAGE_NC_NAME"
		echo "$non_descriptive_nc"
		echo ""
		FOUND=1
	fi
fi

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
