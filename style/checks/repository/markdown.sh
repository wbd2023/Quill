#!/bin/bash
#
# tools/style/checks/repository/markdown.sh
# Checks Markdown formatting and style with markdownlint (STYLE.md 3.2).
#
# Rules:
#	- Markdown files should follow markdownlint policy.
#
# Usage:
#	./tools/style/checks/repository/markdown.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found, 2 if tool unavailable.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
RULE_LABEL="3.2"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] Markdownlint findings:"
MARKDOWNLINT_CONFIG_PATH=".markdownlint.jsonc"
APP_MARKDOWN_FILTER='^README\.md$|^STYLE\.md$|^CONTRIBUTING\.md$|^SECURITY\.md$|^TODO\.md$|'
APP_MARKDOWN_FILTER+='^cmd/|^internal/|^tests/'
TOOLS_MARKDOWN_FILTER='^tools/'
ALL_MARKDOWN_EXCLUDE='^(vendor/|bin/|\.git/)'
VENDOR_EXCLUDE='^vendor/'

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

if ! style_require_command "rg" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_require_command "markdownlint" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files "$STYLE_PATTERN_MARKDOWN" |
			rg -v "$VENDOR_EXCLUDE" |
			rg "$APP_MARKDOWN_FILTER"
	)
	;;
"$STYLE_SCOPE_TOOLS")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files "$STYLE_PATTERN_MARKDOWN" |
			rg -v "$VENDOR_EXCLUDE" |
			rg "$TOOLS_MARKDOWN_FILTER"
	)
	;;
"$STYLE_SCOPE_ALL")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files "$STYLE_PATTERN_MARKDOWN" |
			rg -v "$ALL_MARKDOWN_EXCLUDE"
	)
	;;
esac

if [ "${#FILES[@]}" -eq 0 ]; then
	exit 0
fi

if ! output=$(
	markdownlint -c "$PROJECT_ROOT/$MARKDOWNLINT_CONFIG_PATH" "${FILES[@]/#/$PROJECT_ROOT/}" 2>&1
); then
	echo "$RULE_MESSAGE_PREFIX"
	echo "$output"
	echo ""
	exit 1
fi

exit 0
