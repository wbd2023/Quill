#!/bin/bash
#
# tools/style/checks/go/line-length.sh
# Checks Go line length using visual-width evaluation (STYLE.md 1.1).
#
# Rules:
#	- Go lines must be at most 100 characters.
#
# Usage:
#	./tools/style/checks/go/line-length.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
MAX_LINE_LENGTH=100
TAB_WIDTH=4
FOUND=0
RULE_LABEL="1.1"
RULE_GO_VISUAL_PREFIX="[${RULE_LABEL}] Go line exceeds ${MAX_LINE_LENGTH} columns:"

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

mapfile -t GO_FILES < <(style_collect_files "$PROJECT_ROOT" "$SCOPE" "$STYLE_PATTERN_GO")

for file in "${GO_FILES[@]}"; do
	line_number=0
	while IFS= read -r line || [ -n "$line" ]; do
		line_number=$((line_number + 1))
		expanded_line="${line//$'\t'/    }"
		line_length=${#expanded_line}

		if [ "$line_length" -le "$MAX_LINE_LENGTH" ]; then
			continue
		fi

		echo "$RULE_GO_VISUAL_PREFIX"
		echo "  ${file}:${line_number} (${line_length} columns, tab width ${TAB_WIDTH})"
		echo ""
		FOUND=1
	done <"$file"
done

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
