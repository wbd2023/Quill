#!/bin/bash
#
# tools/style/checks/bash/magic-values.sh
# Recommends avoiding magic numeric values in Bash scripts (STYLE.md 2.10).
#
# Rules:
#	- Prefer named constants for non-trivial numeric values in logic.
#	- Trivial values (0, 1, -1) are allowed.
#
# Usage:
#	./tools/style/checks/bash/magic-values.sh [--scope app|tools|all]
#
# Exit code: 0 if no recommendation findings, 1 if findings exist.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
TRIVIAL_ZERO="0"
TRIVIAL_ONE="1"
TRIVIAL_NEGATIVE_ONE="-1"
EXIT_LITERAL_REGEX='^[[:space:]]*exit[[:space:]]+-?[0-9]+[[:space:]]*$'
COMPARISON_LITERAL_REGEX='^[[:space:]]*[^#].*\-(eq|ne|gt|lt|ge|le)[[:space:]]+-?[0-9]+'
COMPARISON_CAPTURE_REGEX='-(eq|ne|gt|lt|ge|le)[[:space:]]+(-?[0-9]+)'
HEAD_LIMIT_REGEX='^[[:space:]]*[^#].*head[[:space:]]+-[0-9]+'
HEAD_CAPTURE_REGEX='head[[:space:]]+-([0-9]+)'
ARG_COUNT_TOKEN='$#'
RULE_LABEL="R2.10"
MESSAGE_EXIT_CODES="[${RULE_LABEL}] Use named constants for non-trivial exit codes in Bash scripts:"
MESSAGE_COMPARISON="[${RULE_LABEL}] Use named constants for non-trivial numeric comparisons"
MESSAGE_COMPARISON+=" in Bash scripts:"
MESSAGE_HEAD_LIMITS="[${RULE_LABEL}] Use named constants for non-trivial output limits"
MESSAGE_HEAD_LIMITS+=" in Bash scripts:"

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

# --------------------------------------------- Helpers --------------------------------------------

exit_literal_matches=$(
	rg -n "$EXIT_LITERAL_REGEX" "${SCRIPT_FILES[@]}" |
		awk -F: '
			{
				value = $NF
				sub(/^[[:space:]]*exit[[:space:]]+/, "", value)
				sub(/[[:space:]]*$/, "", value)
				if (value != zero && value != one && value != negative_one) {
					print $0
				}
			}' zero="$TRIVIAL_ZERO" one="$TRIVIAL_ONE" negative_one="$TRIVIAL_NEGATIVE_ONE" || true
)

if [ -n "$exit_literal_matches" ]; then
	echo "$MESSAGE_EXIT_CODES"
	echo "$exit_literal_matches"
	echo ""
	FOUND=1
fi

test_literal_matches=$(
	rg -n "$COMPARISON_LITERAL_REGEX" "${SCRIPT_FILES[@]}" |
		awk '
				{
					if (index($0, arg_count_token) > 0) {
						next
					}

				if (match($0, comparison_capture_regex, arr)) {
					value = arr[2]
					if (value != zero && value != one && value != negative_one) {
						print $0
					}
				}
			}' \
			arg_count_token="$ARG_COUNT_TOKEN" \
			comparison_capture_regex="$COMPARISON_CAPTURE_REGEX" \
			zero="$TRIVIAL_ZERO" one="$TRIVIAL_ONE" negative_one="$TRIVIAL_NEGATIVE_ONE" || true
)

if [ -n "$test_literal_matches" ]; then
	echo "$MESSAGE_COMPARISON"
	echo "$test_literal_matches"
	echo ""
	FOUND=1
fi

head_limit_matches=$(
	rg -n "$HEAD_LIMIT_REGEX" "${SCRIPT_FILES[@]}" |
		awk '
			{
				if (match($0, head_capture_regex, arr)) {
					value = arr[1]
					if (value != zero && value != one) {
						print $0
					}
				}
			}' \
			head_capture_regex="$HEAD_CAPTURE_REGEX" \
			zero="$TRIVIAL_ZERO" \
			one="$TRIVIAL_ONE" || true
)

if [ -n "$head_limit_matches" ]; then
	echo "$MESSAGE_HEAD_LIMITS"
	echo "$head_limit_matches"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
