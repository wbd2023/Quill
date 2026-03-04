#!/bin/bash
#
# tools/scripts/check-vertical-spacing.sh
# Recommends vertical spacing for readability (STYLE.md R2.12).
#
# Rules:
#	- Recommend a blank line between consecutive guard clauses.
#
# Usage:
#	./tools/scripts/check-vertical-spacing.sh [--scope app|tools|all]
#
# Exit code: 0 if no recommendation findings, 1 if findings exist.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="R2.12"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] Consider a blank line between consecutive guard clauses:"
GUARD_CLAUSE_PATTERN='(?m)^([ \t]*)if[^\n]*\{\n(?:\1\t[^\n]*\n)*\1\treturn[^\n]*\n'
GUARD_CLAUSE_PATTERN+='\1\}\n\1if[^\n]*\{'

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	mapfile -t FILES < <(
		style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_APP" "$STYLE_PATTERN_GO"
	)
	;;
"$STYLE_SCOPE_TOOLS")
	mapfile -t FILES < <(
		style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_GO"
	)
	;;
"$STYLE_SCOPE_ALL")
	mapfile -t FILES < <(
		{
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_APP" "$STYLE_PATTERN_GO"
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_GO"
		} | sort -u
	)
	;;
esac

if [ "${#FILES[@]}" -eq 0 ]; then
	exit 0
fi

# --------------------------------------------- Checks ---------------------------------------------

matches=$(rg --pcre2 -n -U "$GUARD_CLAUSE_PATTERN" "${FILES[@]}" || true)

if [ -n "$matches" ]; then
	echo "$RULE_MESSAGE_PREFIX"
	echo "$matches"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
