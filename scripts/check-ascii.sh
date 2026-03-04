#!/bin/bash
#
# tools/scripts/check-ascii.sh
# Checks ASCII-only character usage (STYLE.md 1.4).
#
# Rules:
#	- Files should use ASCII unless non-ASCII is strictly necessary.
#	- Vendored and build artefact directories are excluded.
#	- Non-ASCII exceptions must include an inline 'style: allow-non-ascii' marker.
#
# Usage:
#	./tools/scripts/check-ascii.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="1.4"
NON_ASCII_REGEX='[^\x00-\x7F]'
NON_ASCII_ALLOW_MARKER="style: allow-non-ascii"
ASCII_VIOLATION_MESSAGE="[${RULE_LABEL}] Non-ASCII characters detected:"
RG_GLOB_GIT='!.git/**'
RG_GLOB_VENDOR='!vendor/**'
RG_GLOB_BIN='!bin/**'

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_require_command "rg" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	SEARCH_PATHS=("$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD")
	;;
"$STYLE_SCOPE_TOOLS")
	SEARCH_PATHS=("$PROJECT_ROOT/$STYLE_PATH_TOOLS")
	;;
"$STYLE_SCOPE_ALL")
	SEARCH_PATHS=("$PROJECT_ROOT")
	;;
esac

matches=$(rg -nP "$NON_ASCII_REGEX" \
	--hidden \
	--glob "$RG_GLOB_GIT" \
	--glob "$RG_GLOB_VENDOR" \
	--glob "$RG_GLOB_BIN" \
	"${SEARCH_PATHS[@]}" |
	awk '
		{
			if (index($0, marker) > 0) {
				next
			}
			print
		}' marker="$NON_ASCII_ALLOW_MARKER" || true)

if [ -n "$matches" ]; then
	echo "$ASCII_VIOLATION_MESSAGE"
	echo "$matches"
	echo ""
	FOUND=1
fi

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
