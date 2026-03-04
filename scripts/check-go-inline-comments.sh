#!/bin/bash
#
# tools/scripts/check-go-inline-comments.sh
# Checks comment conventions not fully covered by linters (STYLE.md 2.3).
#
# Rules:
#   - Inline trailing comments should start with a lower-case letter.
#   - Inline trailing comments should not end with punctuation.
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

FOUND=0
RULE_LABEL="2.3"
INLINE_COMMENT_PATTERN='^[[:space:]]*[^/[:space:]].*//.*'
COMMENT_DIRECTIVE_REGEX='^(nolint|TODO:|FIXME:|go:|Code\ generated)'
MESSAGE_CASE="[${RULE_LABEL}] Inline comment should start lower-case:"
MESSAGE_PUNCTUATION="[${RULE_LABEL}] Inline comment should not end with punctuation:"
COMMENT_PUNCTUATION_REGEX='[.!?]'
COMMENT_CASE_REGEX='[A-Z]'

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# --------------------------------------------- Checks ---------------------------------------------

while IFS= read -r match; do
	[ -z "$match" ] && continue

	file=$(echo "$match" | cut -d: -f1)
	linenum=$(echo "$match" | cut -d: -f2)
	line=$(echo "$match" | cut -d: -f3-)

	# Extract inline comment payload.
	comment="${line#*//}"
	comment="${comment#"${comment%%[![:space:]]*}"}" # ltrim
	comment="${comment%"${comment##*[![:space:]]}"}" # rtrim

	# Skip directives/markers.
	if [[ "$comment" =~ $COMMENT_DIRECTIVE_REGEX ]]; then
		continue
	fi
	if [ -z "$comment" ]; then
		continue
	fi

	first_char="${comment:0:1}"
	last_char="${comment: -1}"

	if [[ "$first_char" =~ $COMMENT_CASE_REGEX ]]; then
		echo "$MESSAGE_CASE"
		echo "  $file:$linenum"
		echo ""
		FOUND=1
	fi

	if [[ "$last_char" =~ $COMMENT_PUNCTUATION_REGEX ]]; then
		echo "$MESSAGE_PUNCTUATION"
		echo "  $file:$linenum"
		echo ""
		FOUND=1
	fi
done < <(
	rg -n --glob "$STYLE_PATTERN_GO" --glob '!*_test.go' \
		"$INLINE_COMMENT_PATTERN" \
		"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD" || true
)

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
