#!/bin/bash
#
# tools/scripts/check-bash-line-length.sh
# Checks Bash script line length limits (STYLE.md 1.1).
#
# Rules:
#	- Bash script lines must be at most 100 visual columns.
#	- Tabs are counted as 4 columns.
#	- Lines with '# style: allow-long-line' are exempt.
#
# Usage:
#	./tools/scripts/check-bash-line-length.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
MAX_LINE_LENGTH=100
TAB_WIDTH=4
FOUND=0
RULE_LABEL="1.1"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] Bash line exceeds ${MAX_LINE_LENGTH} columns:"
LONG_LINE_MARKER="# style: allow-long-line"

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

mapfile -t SCRIPT_FILES < <(style_collect_files "$PROJECT_ROOT" "$SCOPE" "$STYLE_PATTERN_SHELL")

if [ "${#SCRIPT_FILES[@]}" -eq 0 ]; then
	exit 0
fi

for file in "${SCRIPT_FILES[@]}"; do
	line_number=0

	while IFS= read -r line || [ -n "$line" ]; do
		line_number=$((line_number + 1))
		expanded_line="${line//$'\t'/    }"
		line_length=${#expanded_line}

		if [ "$line_length" -le "$MAX_LINE_LENGTH" ]; then
			continue
		fi

		if [[ "$line" == *"$LONG_LINE_MARKER"* ]]; then
			continue
		fi

		echo "$RULE_MESSAGE_PREFIX"
		echo "  ${file}:${line_number} (${line_length} columns, tab width ${TAB_WIDTH})"
		echo ""
		FOUND=1
	done <"$file"
done

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
